package database

import (
	"context"
	"time"

	"github.com/AkshatPandey-2004/4-IN-A-ROW/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (db *DB) CreateOrGetPlayer(username string, playerID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.Database.Collection("players")

	player := bson.M{
		"_id":        playerID,
		"username":   username,
		"created_at": time.Now(),
	}

	opts := options.Update().SetUpsert(true)
	filter := bson.M{"username": username}
	update := bson.M{"$setOnInsert": player}

	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (db *DB) SaveGame(game *models.Game) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.Database.Collection("games")

	gameDoc := bson.M{
		"_id":         game.ID,
		"player1_id":  nil,
		"player2_id":  nil,
		"winner_id":   nil,
		"is_bot":      game.IsBot,
		"status":      game.Status,
		"created_at":  game.CreatedAt,
		"finished_at": game.FinishedAt,
	}

	if game.Player1 != nil {
		gameDoc["player1_id"] = game.Player1.ID
		gameDoc["player1_username"] = game.Player1.Username
	}

	if game.Player2 != nil {
		gameDoc["player2_id"] = game.Player2.ID
		gameDoc["player2_username"] = game.Player2.Username
	}

	if game.Winner != nil {
		gameDoc["winner_id"] = game.Winner.ID
		gameDoc["winner_username"] = game.Winner.Username
	}

	_, err := collection.InsertOne(ctx, gameDoc)
	return err
}

func (db *DB) UpdateGameStats(game *models.Game) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.Database.Collection("game_stats")

	if game.Winner != nil {
		// Update winner stats
		_, err := collection.UpdateOne(
			ctx,
			bson.M{"username": game.Winner.Username},
			bson.M{
				"$inc": bson.M{
					"wins":        1,
					"total_games": 1,
				},
				"$setOnInsert": bson.M{
					"losses": 0,
					"draws":  0,
				},
			},
			options.Update().SetUpsert(true),
		)
		if err != nil {
			return err
		}

		// Update loser stats
		loser := game.Player1
		if game.Winner.ID == game.Player1.ID {
			loser = game.Player2
		}

		if loser != nil {
			_, err = collection.UpdateOne(
				ctx,
				bson.M{"username": loser.Username},
				bson.M{
					"$inc": bson.M{
						"losses":      1,
						"total_games": 1,
					},
					"$setOnInsert": bson.M{
						"wins":  0,
						"draws": 0,
					},
				},
				options.Update().SetUpsert(true),
			)
			return err
		}
	} else {
		// Draw - update both players
		if game.Player1 != nil {
			collection.UpdateOne(
				ctx,
				bson.M{"username": game.Player1.Username},
				bson.M{
					"$inc": bson.M{
						"draws":       1,
						"total_games": 1,
					},
					"$setOnInsert": bson.M{
						"wins":   0,
						"losses": 0,
					},
				},
				options.Update().SetUpsert(true),
			)
		}

		if game.Player2 != nil {
			collection.UpdateOne(
				ctx,
				bson.M{"username": game.Player2.Username},
				bson.M{
					"$inc": bson.M{
						"draws":       1,
						"total_games": 1,
					},
					"$setOnInsert": bson.M{
						"wins":   0,
						"losses": 0,
					},
				},
				options.Update().SetUpsert(true),
			)
		}
	}

	return nil
}

func (db *DB) GetLeaderboard(limit int) ([]models.LeaderboardEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.Database.Collection("game_stats")

	opts := options.Find().
		SetSort(bson.D{{Key: "wins", Value: -1}, {Key: "total_games", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var leaderboard []models.LeaderboardEntry
	if err = cursor.All(ctx, &leaderboard); err != nil {
		return nil, err
	}

	return leaderboard, nil
}