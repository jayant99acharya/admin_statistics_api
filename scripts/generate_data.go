package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"admin_statistics_api/config"
	"admin_statistics_api/models"
	"admin_statistics_api/utils"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	MIN_ROUNDS = 2000000
	MIN_USERS  = 500
	BATCH_SIZE = 1000
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using default values")
	}

	// Connect to database

	config.ConnectDatabase()
	defer config.DisconnectDatabase()

	fmt.Println("Starting data generation...")
	fmt.Printf("Target: %d rounds with %d unique users\n", MIN_ROUNDS, MIN_USERS)

	// Generate user IDs
	userIDs := generateUserIDs(MIN_USERS)
	fmt.Printf("Generated %d unique user IDs\n", len(userIDs))

	// Generate transactions
	err := generateTransactions(userIDs, MIN_ROUNDS)
	if err != nil {
		log.Fatal("Failed to generate transactions:", err)
	}

	fmt.Println("Data generation completed successfully!")
}

func generateUserIDs(count int) []primitive.ObjectID {
	userIDs := make([]primitive.ObjectID, count)
	for i := 0; i < count; i++ {
		userIDs[i] = primitive.NewObjectID()
	}
	return userIDs
}

func generateTransactions(userIDs []primitive.ObjectID, rounds int) error {
	collection := config.DB.Collection("transactions")
	ctx := context.Background()

	// Clear existing data
	fmt.Println("Clearing existing transactions...")
	_, err := collection.DeleteMany(ctx, map[string]interface{}{})
	if err != nil {
		return fmt.Errorf("failed to clear existing data: %v", err)
	}

	var transactions []interface{}
	rand.Seed(time.Now().UnixNano())

	fmt.Println("Generating transactions...")
	startTime := time.Now()

	for round := 0; round < rounds; round++ {
		if round%100000 == 0 {
			fmt.Printf("Generated %d rounds (%.1f%%)\n", round, float64(round)/float64(rounds)*100)
		}

		// Random user for this round
		userID := userIDs[rand.Intn(len(userIDs))]
		roundID := fmt.Sprintf("round_%d_%d", time.Now().UnixNano(), round)
		currency := utils.GetRandomCurrency()

		// Random time within the past year
		now := time.Now()
		pastYear := now.AddDate(-1, 0, 0)
		randomTime := pastYear.Add(time.Duration(rand.Int63n(int64(now.Sub(pastYear)))))

		// Generate wager amount
		wagerAmount := utils.GetRandomAmount(currency)
		wagerUSDAmount := utils.ConvertToUSD(wagerAmount, currency)

		// Create wager transaction
		wagerAmountDecimal, _ := primitive.ParseDecimal128(fmt.Sprintf("%.8f", wagerAmount))
		wagerUSDDecimal, _ := primitive.ParseDecimal128(fmt.Sprintf("%.2f", wagerUSDAmount))
		
		wagerTransaction := models.Transaction{
			ID:        primitive.NewObjectID(),
			CreatedAt: randomTime,
			UserID:    userID,
			RoundID:   roundID,
			Type:      "Wager",
			Amount:    wagerAmountDecimal,
			Currency:  currency,
			USDAmount: wagerUSDDecimal,
		}

		transactions = append(transactions, wagerTransaction)

		// Generate payout (always after wager)
		payoutTime := randomTime.Add(time.Duration(rand.Intn(300)) * time.Second) // 0-5 minutes later
		payoutMultiplier := utils.GetRandomPayoutMultiplier()
		payoutAmount := wagerAmount * payoutMultiplier
		payoutUSDAmount := utils.ConvertToUSD(payoutAmount, currency)

		// Create payout transaction
		payoutAmountDecimal, _ := primitive.ParseDecimal128(fmt.Sprintf("%.8f", payoutAmount))
		payoutUSDDecimal, _ := primitive.ParseDecimal128(fmt.Sprintf("%.2f", payoutUSDAmount))
		
		payoutTransaction := models.Transaction{
			ID:        primitive.NewObjectID(),
			CreatedAt: payoutTime,
			UserID:    userID,
			RoundID:   roundID,
			Type:      "Payout",
			Amount:    payoutAmountDecimal,
			Currency:  currency,
			USDAmount: payoutUSDDecimal,
		}

		transactions = append(transactions, payoutTransaction)

		// Insert in batches
		if len(transactions) >= BATCH_SIZE {
			if err := insertBatch(collection, transactions); err != nil {
				return err
			}
			transactions = []interface{}{}
		}
	}

	// Insert remaining transactions
	if len(transactions) > 0 {
		if err := insertBatch(collection, transactions); err != nil {
			return err
		}
	}

	elapsed := time.Since(startTime)
	fmt.Printf("Generated %d rounds (%d transactions) in %v\n", rounds, rounds*2, elapsed)

	return nil
}

func insertBatch(collection *mongo.Collection, transactions []interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := collection.InsertMany(ctx, transactions)
	if err != nil {
		return fmt.Errorf("failed to insert batch: %v", err)
	}

	return nil
}