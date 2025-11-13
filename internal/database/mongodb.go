package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewDatabase(mongoURI, dbName string) (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create MongoDB client
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("MongoDB connected successfully")

	database := client.Database(dbName)

	// Create indexes
	if err := createIndexes(ctx, database); err != nil {
		return nil, fmt.Errorf("failed to create indexes: %w", err)
	}

	return &DB{
		Client:   client,
		Database: database,
	}, nil
}

func createIndexes(ctx context.Context, db *mongo.Database) error {
	// Players collection indexes
	playersCollection := db.Collection("players")
	_, err := playersCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    map[string]interface{}{"username": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	// Games collection indexes
	gamesCollection := db.Collection("games")
	_, err = gamesCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: map[string]interface{}{"player1_id": 1}},
		{Keys: map[string]interface{}{"player2_id": 1}},
		{Keys: map[string]interface{}{"status": 1}},
		{Keys: map[string]interface{}{"created_at": -1}},
	})
	if err != nil {
		return err
	}

	// Game stats collection indexes
	statsCollection := db.Collection("game_stats")
	_, err = statsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    map[string]interface{}{"username": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	log.Println("MongoDB indexes created successfully")
	return nil
}

func (db *DB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return db.Client.Disconnect(ctx)
}