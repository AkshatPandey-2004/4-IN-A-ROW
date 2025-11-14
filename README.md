# 4-IN-A-ROW 🎮

A real-time multiplayer Connect Four game built with Go backend and React frontend, featuring WebSocket communication, Kafka event streaming, and MongoDB persistence.

## 🌐 Live Demo

**🎮 Play Now:** [http://13.235.128.137:8082/](http://13.235.128.137:8082/)

---

## ✨ Features

- 🎯 **Real-time Multiplayer** - Play against other players via WebSocket
- 🤖 **AI Bot Opponent** - Practice against an intelligent minimax algorithm bot
- 📊 **Live Leaderboard** - Track top players with real-time statistics
- 🔄 **Event Streaming** - Kafka-powered game analytics
- 💾 **Data Persistence** - MongoDB for game history and player stats
- 🎨 **Modern UI** - Responsive React frontend with animations

---

## 🏗️ Tech Stack

**Frontend:** React.js, WebSocket, CSS3  
**Backend:** Go 1.21, Gorilla WebSocket, Kafka Client  
**Infrastructure:** Docker, Nginx, Kafka, Zookeeper, MongoDB Atlas  

---

## 🚀 Running Locally

### Prerequisites

Before you begin, ensure you have:

- **Go** 1.21+ ([Download](https://golang.org/dl/))
- **Node.js** 18+ and npm ([Download](https://nodejs.org/))
- **Docker Desktop** ([Download](https://www.docker.com/products/docker-desktop/))
- **Git** ([Download](https://git-scm.com/))

> **Note:** MongoDB connection credentials are already configured in the `.env` file (included in repo for placement evaluation purposes).

---

### Step 1: Clone the Repository

```bash
git clone https://github.com/AkshatPandey-2004/4-IN-A-ROW.git
cd 4-IN-A-ROW
```

---

### Step 2: Start Kafka & Zookeeper

Kafka is required for event streaming. Run it using Docker:

```bash
# Navigate to deployments directory
cd deployments

# Start Kafka and Zookeeper
docker-compose up -d zookeeper kafka

# Wait for services to start (about 30 seconds)
sleep 30

# Verify they're running
docker ps
```

You should see:
```
CONTAINER ID   IMAGE                         STATUS         PORTS
xxxxx          confluentinc/cp-kafka         Up 30 seconds  0.0.0.0:9092->9092/tcp
xxxxx          confluentinc/cp-zookeeper     Up 30 seconds  0.0.0.0:2181->2181/tcp
```

---

### Step 3: Configure Environment (Already Done!)

The `.env` file is already included in the repository:

```bash
# .env (at project root)
PORT=8081
MONGODB_URI="mongodb+srv://Admin:admin@cluster0.wrnmkpn.mongodb.net/fourinarow?retryWrites=true&w=majority&appName=Cluster0"
MONGODB_DATABASE="fourinarow"
KAFKA_BROKERS=kafka:9092
```

> **No additional setup needed!** MongoDB is already configured and accessible.

---

### Step 4: Run the Backend (Go Server)

Open a new terminal:

```bash
# Navigate to project root
cd 4-IN-A-ROW

# Install Go dependencies
go mod download

# Load environment variables and run
source .env && go run cmd/server/main.go
```

**Or on Windows (PowerShell):**
```powershell
$env:MONGODB_URI="mongodb+srv://Admin:admin@cluster0.wrnmkpn.mongodb.net/fourinarow?retryWrites=true&w=majority&appName=Cluster0"
$env:KAFKA_BROKERS="localhost:9092"
$env:PORT="8081"
go run cmd/server/main.go
```

**Expected Output:**
```
2025/11/14 00:44:55 MongoDB connected successfully
2025/11/14 00:44:55 MongoDB indexes created successfully
2025/11/14 00:44:55 Kafka producer created successfully
2025/11/14 00:44:55 Kafka consumer created successfully
2025/11/14 00:44:55 Server starting on port 8081
```

**Test Backend:**
```bash
# In another terminal
curl http://localhost:8081/api/health
# Should return: {"status":"ok"}

curl http://localhost:8081/api/leaderboard
# Should return player data
```

---

### Step 5: Run the Frontend (React App)

Open another new terminal:

```bash
# Navigate to frontend directory
cd 4-IN-A-ROW/frontend

# Install dependencies (first time only)
npm install

# Start development server
npm start
```

**Expected Output:**
```
Compiled successfully!

You can now view frontend in the browser.

  Local:            http://localhost:3000
  On Your Network:  http://192.168.x.x:3000

webpack compiled successfully
```

The browser should automatically open `http://localhost:3000`

---

### Step 6: Play the Game! 🎮

1. **Open:** `http://localhost:3000` in your browser
2. **Enter your username**
3. **Choose:**
   - 🤖 **Play vs Bot** - Instant game against AI
   - 👥 **Find Opponent** - Wait for another player to join

---

## 🎮 Testing Multiplayer Locally

To test multiplayer functionality:

### Option 1: Two Browser Windows

1. Open `http://localhost:3000` in **Chrome**
2. Open `http://localhost:3000` in **Firefox** (or Chrome Incognito)
3. Enter **different usernames** in each
4. Click **"Find Opponent"** in both windows
5. The system will match you together!

### Option 2: Multiple Tabs

1. Open two tabs at `http://localhost:3000`
2. Use different usernames
3. Click "Find Opponent" in both
4. Play against yourself!

---

## 📁 Project Structure

```
4-IN-A-ROW/
├── .env                        # Environment variables (included)
├── cmd/
│   └── server/
│       └── main.go            # Backend entry point
├── internal/
│   ├── game/
│   │   ├── board.go          # Game logic (6x7 grid, win detection)
│   │   ├── game.go           # Game instance management
│   │   ├── manager.go        # Thread-safe game state manager
│   │   └── bot.go            # Minimax AI implementation
│   ├── handlers/
│   │   ├── websocket.go      # WebSocket connection handlers
│   │   └── api.go            # REST API handlers
│   └── storage/
│       └── mongodb.go        # MongoDB operations
├── pkg/
│   ├── models/
│   │   ├── game.go           # Game data models
│   │   └── player.go         # Player data models
│   └── kafka/
│       ├── producer.go       # Kafka event producer
│       └── consumer.go       # Kafka event consumer
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   ├── Board.jsx     # Game board UI
│   │   │   ├── Leaderboard.jsx
│   │   │   └── GameInfo.jsx
│   │   ├── hooks/
│   │   │   └── useWebSocket.js
│   │   ├── App.js            # Main React component
│   │   └── index.js
│   ├── package.json
│   └── public/
├── deployments/
│   ├── docker-compose.yml    # Docker services config
│   ├── Dockerfile.backend
│   ├── Dockerfile.frontend
│   └── nginx.conf            # Nginx reverse proxy config
├── go.mod
├── go.sum
└── README.md
```

---

## 🐳 Alternative: Run Everything with Docker

If you prefer to run everything in Docker (no local Go/Node setup needed):

```bash
cd deployments

# Build and start all services
docker-compose up -d --build

# Wait for build to complete (2-3 minutes)
# Access at http://localhost:8082
```

This starts:
- ✅ Zookeeper (Kafka coordination)
- ✅ Kafka (Event streaming)
- ✅ Backend (Go API on port 8081)
- ✅ Frontend (React + Nginx on port 8082)

**Stop all services:**
```bash
docker-compose down
```

---

## 🧪 Testing

### Backend Tests

```bash
# Run all tests
go test ./...

# Test with verbose output
go test -v ./internal/game/

# Test with coverage
go test -cover ./...
```

### Frontend Tests

```bash
cd frontend

# Run tests
npm test

# Run with coverage
npm test -- --coverage
```

### API Testing

```bash
# Health check
curl http://localhost:8081/api/health

# Get leaderboard
curl http://localhost:8081/api/leaderboard

# Get player stats
curl http://localhost:8081/api/stats/player1
