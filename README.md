# Admin Statistics API

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![MongoDB](https://img.shields.io/badge/MongoDB-7.0-green.svg)](https://mongodb.com)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)](https://docker.com)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A high-performance GoLang REST API for aggregating casino transaction statistics using MongoDB and Redis caching.

> **Note**: This is a demonstration project created for technical interview purposes, showcasing advanced GoLang and MongoDB development skills.

## Features

- **High Performance**: Optimized MongoDB aggregation queries with Redis caching
- **Scalable Architecture**: Clean separation of concerns with services, handlers, and middleware
- **Authentication**: Secure API endpoints with token-based authentication
- **Input Validation**: Comprehensive request validation using go-playground/validator
- **Data Generation**: Utility script to generate 2M+ sample transactions
- **Efficient Queries**: Advanced MongoDB aggregation pipelines for complex statistics

## Technology Stack

- **Language**: Go 1.21+
- **Database**: MongoDB
- **Cache**: Redis (optional)
- **Web Framework**: Gin
- **Validation**: go-playground/validator
- **MongoDB Driver**: Official MongoDB Go Driver

## API Endpoints

### Authentication
All API endpoints (except `/health`) require an `Authorization` header:
```
Authorization: your-secret-token-here
```

> **Security Note**: The default token is for development only. In production, use a strong, unique token and consider implementing JWT or OAuth2.

### Endpoints

1. **Health Check**
   ```
   GET /health
   ```

2. **Gross Gaming Revenue**
   ```
   GET /gross_gaming_rev?from=2024-01-01&to=2024-12-31
   ```
   Calculates GGR (Wagers - Payouts) by currency and USD.

3. **Daily Wager Volume**
   ```
   GET /daily_wager_volume?from=2024-01-01&to=2024-12-31
   ```
   Returns daily wager volumes by currency and USD.

4. **User Wager Percentile**
   ```
   GET /user/{user_id}/wager_percentile?from=2024-01-01&to=2024-12-31
   ```
   Calculates user's wager percentile ranking.

## Quick Start

### Clone the Repository
```bash
git clone https://github.com/jayant99acharya/admin_statistics_api.git
```

### Option 1: Docker (Recommended)

1. **Prerequisites**
   - Docker and Docker Compose installed

2. **Quick Start**
   ```bash
   docker-compose up -d
   ```
   This will start MongoDB, Redis, and the API automatically.

3. **Generate Sample Data**
   ```bash
   # Connect to the running container
   docker-compose exec app go run scripts/generate_data.go
   ```

4. **Access the API**
   - API: http://localhost:8090
   - MongoDB Express (DB Admin): http://localhost:8081 (admin/admin)

### Option 2: Local Development

1. **Prerequisites**
   - Go 1.21+
   - MongoDB running on localhost:27017
   - Redis running on localhost:6379 (optional)

2. **Setup**
   ```bash
   cp .env.example .env
   go mod tidy
   go run scripts/generate_data.go
   go run main.go
   ```

The API will be available at `http://localhost:8080`

## Usage Examples

### Using curl

1. **Health Check**
   ```bash
   # Docker setup (port 8090)
   curl http://localhost:8090/health
   
   # Local development (port 8080)
   curl http://localhost:8080/health
   ```

2. **Get Gross Gaming Revenue**
   ```bash
   # Docker setup
   curl -H "Authorization: your-secret-token-here" \
        "http://localhost:8090/gross_gaming_rev?from=2024-01-01&to=2024-12-31"
   
   # Local development
   curl -H "Authorization: your-secret-token-here" \
        "http://localhost:8080/gross_gaming_rev?from=2024-01-01&to=2024-12-31"
   ```

3. **Get Daily Wager Volume**
   ```bash
   # Docker setup
   curl -H "Authorization: your-secret-token-here" \
        "http://localhost:8090/daily_wager_volume?from=2024-01-01&to=2024-12-31"
   ```

4. **Get User Wager Percentile**
   ```bash
   # First, get a user ID from the database
   # Docker setup
   curl -H "Authorization: your-secret-token-here" \
        "http://localhost:8090/user/507f1f77bcf86cd799439011/wager_percentile?from=2024-01-01&to=2024-12-31"
   ```

### Sample Responses

**Gross Gaming Revenue:**
```json
{
  "success": true,
  "data": {
    "from": "2024-01-01",
    "to": "2024-12-31",
    "gross_gaming_revenue": [
      {
        "currency": "BTC",
        "amount": 125.45,
        "usdValue": 5647250.00
      },
      {
        "currency": "ETH",
        "amount": 2340.67,
        "usdValue": 7022010.00
      },
      {
        "currency": "USDT",
        "amount": 1250000.00,
        "usdValue": 1250000.00
      }
    ]
  }
}
```

## Database Schema

### Transaction Collection
```go
type Transaction struct {
    ID        primitive.ObjectID   `bson:"_id"`
    CreatedAt time.Time            `bson:"createdAt"`
    UserID    primitive.ObjectID   `bson:"userId"`
    RoundID   string               `bson:"roundId"`
    Type      string               `bson:"type"`      // "Wager" or "Payout"
    Amount    primitive.Decimal128 `bson:"amount"`
    Currency  string               `bson:"currency"`  // "ETH", "BTC", or "USDT"
    USDAmount primitive.Decimal128 `bson:"usdAmount"`
}
```

## Performance Optimizations

1. **Database Indexes**: Automatically created on startup
   - `createdAt` (for time-based queries)
   - `userId` (for user-specific queries)
   - `userId + createdAt` (compound index)
   - `roundId` (for round-based queries)
   - `type` (for filtering wagers/payouts)

2. **Redis Caching**: Results cached for 5 minutes
   - Gross Gaming Revenue
   - Daily Wager Volume
   - User Wager Percentiles

3. **Efficient Aggregation**: MongoDB aggregation pipelines optimized for large datasets

## Development

### Project Structure
```
admin_stats_api/
├── config/          # Database and cache configuration
├── handlers/        # HTTP request handlers
├── middleware/      # Authentication middleware
├── models/          # Data models and structs
├── scripts/         # Utility scripts (data generation)
├── services/        # Business logic and database operations
├── utils/           # Helper functions
├── main.go          # Application entry point
├── go.mod           # Go module definition
├── Dockerfile       # Container build configuration
├── docker-compose.yml # Multi-service setup
└── README.md        # This file
```

### Running Tests
```bash
# Local testing
go test ./...

# Docker testing
docker-compose exec app go test ./...
```

### Building for Production
```bash
# Docker build
docker build -t admin-stats-api .

# Or use docker-compose
docker-compose up --build
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `MONGODB_URI` | MongoDB connection string | `mongodb://localhost:27017` |
| `DB_NAME` | Database name | `admin_statistics` |
| `REDIS_ADDR` | Redis server address | `localhost:6379` |
| `REDIS_PASSWORD` | Redis password | `` |
| `REDIS_DB` | Redis database number | `0` |
| `AUTH_TOKEN` | API authentication token | `your-secret-token-here` |
| `PORT` | Server port | `8080` |
| `GIN_MODE` | Gin framework mode | `debug` |

> **Security Warning**: Always use strong, unique tokens in production environments.

## Troubleshooting

1. **MongoDB Connection Issues**
   - Ensure MongoDB is running: `brew services list | grep mongodb`
   - Check connection string in `.env` file

2. **Redis Connection Issues**
   - Redis is optional; the API will work without it
   - Check Redis status: `brew services list | grep redis`

3. **Data Generation Takes Too Long**
   - The script generates 2M+ transactions, which may take 10-30 minutes
   - Monitor progress in the console output
   - Reduce `MIN_ROUNDS` in `scripts/generate_data.go` for faster testing

4. **Memory Issues During Data Generation**
   - Increase batch size in `generate_data.go`
   - Ensure sufficient RAM (recommended: 8GB+)

## Contributing

This project is primarily for demonstration purposes. If you'd like to contribute:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built for technical interview demonstration
- Showcases advanced GoLang and MongoDB development patterns
- Demonstrates production-ready API architecture and best practices