package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"admin_statistics_api/config"
	"admin_statistics_api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type StatisticsService struct {
	collection *mongo.Collection
}

func NewStatisticsService() *StatisticsService {
	return &StatisticsService{
		collection: config.DB.Collection("transactions"),
	}
}

// GetGrossGamingRevenue calculates GGR (Wagers - Payouts) by currency
func (s *StatisticsService) GetGrossGamingRevenue(from, to time.Time) ([]models.GrossGamingRevenue, error) {
	ctx := context.Background()
	
	// Try to get from cache first
	cacheKey := fmt.Sprintf("ggr:%d:%d", from.Unix(), to.Unix())
	if cached, err := config.GetCache(cacheKey); err == nil {
		var result []models.GrossGamingRevenue
		if json.Unmarshal([]byte(cached), &result) == nil {
			return result, nil
		}
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"createdAt": bson.M{
					"$gte": from,
					"$lte": to,
				},
			},
		},
		{
			"$group": bson.M{
				"_id": bson.M{
					"currency": "$currency",
					"type":     "$type",
				},
				"totalAmount":    bson.M{"$sum": bson.M{"$toDouble": "$amount"}},
				"totalUSDAmount": bson.M{"$sum": bson.M{"$toDouble": "$usdAmount"}},
			},
		},
		{
			"$group": bson.M{
				"_id": "$_id.currency",
				"wagers": bson.M{
					"$sum": bson.M{
						"$cond": bson.A{
							bson.M{"$eq": bson.A{"$_id.type", "Wager"}},
							"$totalAmount",
							0,
						},
					},
				},
				"payouts": bson.M{
					"$sum": bson.M{
						"$cond": bson.A{
							bson.M{"$eq": bson.A{"$_id.type", "Payout"}},
							"$totalAmount",
							0,
						},
					},
				},
				"wagersUSD": bson.M{
					"$sum": bson.M{
						"$cond": bson.A{
							bson.M{"$eq": bson.A{"$_id.type", "Wager"}},
							"$totalUSDAmount",
							0,
						},
					},
				},
				"payoutsUSD": bson.M{
					"$sum": bson.M{
						"$cond": bson.A{
							bson.M{"$eq": bson.A{"$_id.type", "Payout"}},
							"$totalUSDAmount",
							0,
						},
					},
				},
			},
		},
		{
			"$project": bson.M{
				"currency": "$_id",
				"ggr":      bson.M{"$subtract": bson.A{"$wagers", "$payouts"}},
				"ggrUSD":   bson.M{"$subtract": bson.A{"$wagersUSD", "$payoutsUSD"}},
			},
		},
	}

	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.GrossGamingRevenue
	for cursor.Next(ctx) {
		var doc struct {
			Currency string  `bson:"currency"`
			GGR      float64 `bson:"ggr"`
			GGRUSD   float64 `bson:"ggrUSD"`
		}
		if err := cursor.Decode(&doc); err != nil {
			continue
		}
		results = append(results, models.GrossGamingRevenue{
			Currency: doc.Currency,
			Amount:   doc.GGR,
			USDValue: doc.GGRUSD,
		})
	}

	// Cache the result for 5 minutes
	if resultJSON, err := json.Marshal(results); err == nil {
		config.SetCache(cacheKey, string(resultJSON), 5*time.Minute)
	}

	return results, nil
}

// GetDailyWagerVolume calculates daily wager volume by currency
func (s *StatisticsService) GetDailyWagerVolume(from, to time.Time) ([]models.DailyWagerVolume, error) {
	ctx := context.Background()
	
	// Try to get from cache first
	cacheKey := fmt.Sprintf("daily_wager:%d:%d", from.Unix(), to.Unix())
	if cached, err := config.GetCache(cacheKey); err == nil {
		var result []models.DailyWagerVolume
		if json.Unmarshal([]byte(cached), &result) == nil {
			return result, nil
		}
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"createdAt": bson.M{
					"$gte": from,
					"$lte": to,
				},
				"type": "Wager",
			},
		},
		{
			"$group": bson.M{
				"_id": bson.M{
					"date": bson.M{
						"$dateToString": bson.M{
							"format": "%Y-%m-%d",
							"date":   "$createdAt",
						},
					},
					"currency": "$currency",
				},
				"totalAmount":    bson.M{"$sum": bson.M{"$toDouble": "$amount"}},
				"totalUSDAmount": bson.M{"$sum": bson.M{"$toDouble": "$usdAmount"}},
			},
		},
		{
			"$project": bson.M{
				"date":     "$_id.date",
				"currency": "$_id.currency",
				"amount":   "$totalAmount",
				"usdValue": "$totalUSDAmount",
			},
		},
		{
			"$sort": bson.M{
				"date":     1,
				"currency": 1,
			},
		},
	}

	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.DailyWagerVolume
	for cursor.Next(ctx) {
		var doc struct {
			Date     string  `bson:"date"`
			Currency string  `bson:"currency"`
			Amount   float64 `bson:"amount"`
			USDValue float64 `bson:"usdValue"`
		}
		if err := cursor.Decode(&doc); err != nil {
			continue
		}
		results = append(results, models.DailyWagerVolume{
			Date:     doc.Date,
			Currency: doc.Currency,
			Amount:   doc.Amount,
			USDValue: doc.USDValue,
		})
	}

	// Cache the result for 5 minutes
	if resultJSON, err := json.Marshal(results); err == nil {
		config.SetCache(cacheKey, string(resultJSON), 5*time.Minute)
	}

	return results, nil
}

// GetUserWagerPercentile calculates user's wager percentile
func (s *StatisticsService) GetUserWagerPercentile(userID primitive.ObjectID, from, to time.Time) (*models.UserWagerPercentile, error) {
	ctx := context.Background()
	
	// Try to get from cache first
	cacheKey := fmt.Sprintf("user_percentile:%s:%d:%d", userID.Hex(), from.Unix(), to.Unix())
	if cached, err := config.GetCache(cacheKey); err == nil {
		var result models.UserWagerPercentile
		if json.Unmarshal([]byte(cached), &result) == nil {
			return &result, nil
		}
	}

	// First, get all users' total wager amounts
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"createdAt": bson.M{
					"$gte": from,
					"$lte": to,
				},
				"type": "Wager",
			},
		},
		{
			"$group": bson.M{
				"_id":              "$userId",
				"totalWageredUSD": bson.M{"$sum": bson.M{"$toDouble": "$usdAmount"}},
			},
		},
		{
			"$sort": bson.M{
				"totalWageredUSD": -1,
			},
		},
	}

	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var userWagers []struct {
		UserID           primitive.ObjectID `bson:"_id"`
		TotalWageredUSD float64            `bson:"totalWageredUSD"`
	}

	for cursor.Next(ctx) {
		var doc struct {
			UserID           primitive.ObjectID `bson:"_id"`
			TotalWageredUSD float64            `bson:"totalWageredUSD"`
		}
		if err := cursor.Decode(&doc); err != nil {
			continue
		}
		userWagers = append(userWagers, doc)
	}

	if len(userWagers) == 0 {
		return nil, fmt.Errorf("no wager data found for the specified time period")
	}

	// Find the target user's position and calculate percentile
	var targetUserWager float64
	var userRank int
	found := false

	for i, wager := range userWagers {
		if wager.UserID == userID {
			targetUserWager = wager.TotalWageredUSD
			userRank = i + 1
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("user not found in wager data for the specified time period")
	}

	totalUsers := len(userWagers)
	percentile := float64(totalUsers-userRank+1) / float64(totalUsers) * 100

	result := &models.UserWagerPercentile{
		UserID:       userID.Hex(),
		TotalWagered: targetUserWager,
		Percentile:   percentile,
		Rank:         userRank,
		TotalUsers:   totalUsers,
	}

	// Cache the result for 5 minutes
	if resultJSON, err := json.Marshal(result); err == nil {
		config.SetCache(cacheKey, string(resultJSON), 5*time.Minute)
	}

	return result, nil
}