package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "os/signal"
    "strings"
    "syscall"
    "time"
    
    "github.com/AkshatPandey-2004/4-in-a-row/internal/game"
    "github.com/AkshatPandey-2004/4-in-a-row/internal/matchmaking"
    "github.com/AkshatPandey-2004/4-in-a-row/internal/websocket"
    "github.com/AkshatPandey-2004/4-in-a-row/internal/database"
    "github.com/AkshatPandey-2004/4-in-a-row/internal/kafka"
    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
    "github.com/rs/cors"
)

func main() {
    // Load environment variables
    godotenv.Load()

    // Initialize database
    db, err := database.NewDatabase(
        getEnv("DB_HOST", "localhost"),
        getEnv("DB_PORT", "5432"),
        getEnv("DB_USER", "postgres"),
        getEnv("DB_PASSWORD", "postgres"),
        getEnv("DB_NAME", "fourinarow"),
    )
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()

    // Initialize Kafka producer
    kafkaBrokers := strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ",")
    kafkaProducer, err := kafka.NewProducer(kafkaBrokers, "game-events")
    if err != nil {
        log.Fatal("Failed to create Kafka producer:", err)
    }
    defer kafkaProducer.Close()

    // Initialize Kafka consumer
    kafkaConsumer, err := kafka.NewConsumer(kafkaBrokers, "analytics-group", "game-events")
    if err != nil {
        log.Fatal("Failed to create Kafka consumer:", err)
    }
    defer kafkaConsumer.Close()

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    kafkaConsumer.Start(ctx)

    // Initialize components
    hub := websocket.NewHub()
    go hub.Run()

    gameManager := game.NewManager()
    matchmaker := matchmaking.NewMatchmaker()
    wsHandler := websocket.NewHandler(hub, gameManager, matchmaker, db, kafkaProducer)

    // Setup HTTP router
    router := mux.NewRouter()
    
    router.HandleFunc("/ws", wsHandler.HandleWebSocket)
    router.HandleFunc("/api/leaderboard", getLeaderboardHandler(db)).Methods("GET")
    router.HandleFunc("/api/health", healthCheckHandler).Methods("GET")

    // CORS
    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"*"},
        AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
        AllowedHeaders:   []string{"*"},
        AllowCredentials: true,
    })

    handler := c.Handler(router)

    // Start server
    port := getEnv("PORT", "8081")
    server := &http.Server{
        Addr:         ":" + port,
        Handler:      handler,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
    }

    // Graceful shutdown
    go func() {
        log.Printf("Server starting on port %s", port)
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal("Server failed:", err)
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Println("Shutting down server...")
    ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }

    log.Println("Server exited")
}

func getLeaderboardHandler(db *database.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        leaderboard, err := db.GetLeaderboard(10)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(leaderboard)
    }
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}