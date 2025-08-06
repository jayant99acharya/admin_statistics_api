package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	ID        primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time              `bson:"createdAt" json:"createdAt"`
	UserID    primitive.ObjectID     `bson:"userId" json:"userId"`
	RoundID   string                 `bson:"roundId" json:"roundId"`
	Type      string                 `bson:"type" json:"type" validate:"required,oneof=Wager Payout"`
	Amount    primitive.Decimal128   `bson:"amount" json:"amount"`
	Currency  string                 `bson:"currency" json:"currency" validate:"required,oneof=ETH BTC USDT"`
	USDAmount primitive.Decimal128   `bson:"usdAmount" json:"usdAmount"`
}

type GrossGamingRevenue struct {
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
	USDValue float64 `json:"usdValue"`
}

type DailyWagerVolume struct {
	Date     string  `json:"date"`
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
	USDValue float64 `json:"usdValue"`
}

type UserWagerPercentile struct {
	UserID         string  `json:"userId"`
	TotalWagered   float64 `json:"totalWagered"`
	Percentile     float64 `json:"percentile"`
	Rank           int     `json:"rank"`
	TotalUsers     int     `json:"totalUsers"`
}