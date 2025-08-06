package utils

import (
	"math/rand"
	"time"
)

// Static conversion rates for the example (in a real app, these would come from an API)
var CurrencyRates = map[string]float64{
	"BTC":  45000.0, // 1 BTC = $45,000
	"ETH":  3000.0,  // 1 ETH = $3,000
	"USDT": 1.0,     // 1 USDT = $1
}

func ConvertToUSD(amount float64, currency string) float64 {
	rate, exists := CurrencyRates[currency]
	if !exists {
		return 0
	}
	return amount * rate
}

func GetRandomCurrency() string {
	currencies := []string{"BTC", "ETH", "USDT"}
	return currencies[rand.Intn(len(currencies))]
}

func GetRandomAmount(currency string) float64 {
	rand.Seed(time.Now().UnixNano())
	
	switch currency {
	case "BTC":
		// Random amount between 0.001 and 0.1 BTC
		return 0.001 + rand.Float64()*(0.1-0.001)
	case "ETH":
		// Random amount between 0.01 and 1 ETH
		return 0.01 + rand.Float64()*(1.0-0.01)
	case "USDT":
		// Random amount between 10 and 1000 USDT
		return 10 + rand.Float64()*(1000-10)
	default:
		return 0
	}
}

func GetRandomPayoutMultiplier() float64 {
	// Random payout multiplier between 0 (total loss) and 2.5 (150% profit)
	return rand.Float64() * 2.5
}