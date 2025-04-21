# 🎰 Bet Settlement Engine

**`bet-settlement-engine`** is a simple and efficient REST API service written in **Golang** that allows users to place and settle bets in real-time using **in-memory storage**. Designed for simplicity, performance, and concurrency safety, this project demonstrates the core functionality of a betting backend system without relying on any external databases.

---

## 🎯 Key Features

- 🧾 **Place a Bet** – Users can place bets by specifying the amount, odds, and event ID.
- ✅ **Settle a Bet** – Bets can be settled based on the result of an event (win/lose).
- 💰 **Check Balance** – Users can retrieve their current balance.
- 🔒 **Concurrency-safe** – All critical operations are protected using `sync.Mutex` to avoid race conditions in a concurrent environment.
- 🧠 **In-Memory Logic** – Fast, ephemeral storage ideal for lightweight simulations and testing.

---

## 🧪 Use Case

This project is ideal for:

- Demonstrating understanding of Go’s concurrency and memory management.
- Simulating core functionality of betting platforms.
- Building the foundation for more robust systems with persistent databases later.

---

## 🛠️ Built With

- **Go (Golang)** – Fast, simple, and memory-safe language.
- **`net/http`** – Native HTTP server for API endpoints.
- **`sync` package** – For thread-safe in-memory data structures.

---

## 📂 Repository Structure



