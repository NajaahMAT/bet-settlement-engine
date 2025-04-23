
# 🎰 Bet Settlement Engine

**`bet-settlement-engine`** is a simple and efficient REST API service written in **Golang** that allows users to place and settle bets in real-time using **Redis** as a backend data store. Designed for simplicity, performance, and concurrency safety, this project demonstrates the core functionality of a betting backend system with clean code, safe operations, and extensible design.

---

## 🎯 Key Features

- 🧾 **Place a Bet** – Users can place bets by specifying the amount, odds, and event ID.
- ✅ **Settle a Bet** – Bets can be settled based on the result of an event (`win` or `lose`).
- 💰 **Check Balance** – Users can retrieve their current balance.
- 💵 **Initial Balance** – New users are initialized with a default balance of `1000` units.
- 🔒 **Concurrency-Safe Operations** – Uses `sync.Mutex` to lock user balance updates and avoid race conditions.
- ⚡ **Redis Storage** – Leverages Redis for fast, in-memory storage of user balances and bet details.

---

## 🧪 Use Case

This project is ideal for:

- Demonstrating understanding of Go’s concurrency and memory management.
- Simulating core functionality of betting platforms.
- Learning to build stateless APIs with Redis integration.
- Serving as a backend skeleton for further enhancements with persistent databases or complex logic.

---

## 🛠️ Tech Stack

- **Go (Golang)** – High-performance, compiled backend language.
- **Redis** – In-memory key-value store for managing balances and bet info.
- **`net/http`** – Native HTTP server in Go for routing and handling API requests.
- **`sync` package** – Ensures concurrency-safe data operations using `Mutex`.

---

## 🧠 Concurrency Safety with `sync.Mutex`

When a user places a bet or when a bet is settled, the application updates the user's balance. These operations are critical and can become problematic in concurrent environments.

To ensure **thread safety**, the application wraps balance update operations with a `sync.Mutex` lock. This prevents race conditions such as:

- Deducting more than the available balance.
- Overwriting updates in a multi-threaded environment.

Even though Redis supports atomic operations, this additional layer ensures safe composite operations (like get-update-set sequences).

---

## 🧾 Redis Integration

The service uses Redis for all backend storage operations.

### Redis Configuration

Environment variables used for Redis:

```bash
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

### Connection Code (Simplified)

```go
r.rdb = redis.NewClient(&redis.Options{
    Addr:     fmt.Sprintf("%s:%s", r.conf.Host, r.conf.Port),
    Password: r.conf.Password,
    DB:       r.conf.DB,
})
```

On startup, the application will attempt to `PING` the Redis instance to ensure it's reachable.

---

## 📦 API Endpoints

### ➕ Place a Bet

**POST** `/api/v1/bet/place`

**Request Body**:

```json
{
  "user_id": "Najaah",
  "event_id": "event001",
  "amount": 100,
  "odds": 2.5
}
```

**Response**:

```json
{
  "msg": "Bet saved successfully",
  "bet_id": "f57ceb92-4df1-4c18-aae6-683915d7a039",
  "user_id": "Najaah",
  "event_id": "event001",
  "amount": 100,
  "odds": 2.5,
  "result": "placed"
}
```

---

### ✅ Settle a Bet

**POST** `/api/v1/bet/settle`

**Request Body**:

```json
{
  "bet_id": "f57ceb92-4df1-4c18-aae6-683915d7a039",
  "result": "win"
}
```

**Response**:

```json
{
  "msg": "Bet settled successfully",
  "bet_id": "f57ceb92-4df1-4c18-aae6-683915d7a039",
  "user_id": "Najaah",
  "amount_won": 250
}
```

---

### 💰 Check Balance

**GET** `/api/v1/balance?user_id=Najaah`

**Response**:

```json
{
  "user_id": "Najaah",
  "balance": 1150
}
```

---

## 🚀 Running the Project

```bash
# Set environment variables or add to .env
export REDIS_HOST=localhost
export REDIS_PORT=6379
export REDIS_PASSWORD=
export REDIS_DB=0

# Run the project
go run cmd/main.go
```

---

## ✅ Sample Test Flow

1. **Check balance** – should be 1000 if new user.
2. **Place a bet** – e.g., 100 at odds 2.5.
3. **Check balance again** – should be deducted.
4. **Settle the bet as win** – balance should increase with winnings.
5. **Settle the bet as lose** – no balance change.

---

## ✍️ Author

Made by [Najaah](https://github.com/NajaahMAT)

---
