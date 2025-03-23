package database

import (
	"context"
	"time"

	"github.com/smokersaim/droqsic/cmd/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func ConnectMongoDB(cfg *config.Config, log *zap.Logger) (*MongoDB, error) {
	clientOptions := options.Client().
		ApplyURI(cfg.Database.URI).
		SetMinPoolSize(10).
		SetMaxPoolSize(100).
		SetConnectTimeout(10 * time.Second).
		SetSocketTimeout(30 * time.Second).
		SetServerSelectionTimeout(5 * time.Second).
		SetMaxConnIdleTime(5 * time.Minute).
		SetRetryWrites(true)

	var client *mongo.Client
	var err error

	for i := 1; i <= 3; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		client, err = mongo.Connect(ctx, clientOptions)
		if err == nil {
			if err = client.Ping(ctx, nil); err == nil {
				log.Info("MongoDB connection established successfully")
				return &MongoDB{
					Client:   client,
					Database: client.Database(cfg.Database.Database),
				}, nil
			}
		}

		log.Warn("MongoDB connection attempt failed", zap.Int("attempt", i), zap.Error(err))
		time.Sleep(2 * time.Second)
	}

	log.Error("MongoDB connection failed after multiple attempts", zap.Error(err))
	return nil, err
}

func (m *MongoDB) Disconnect(log *zap.Logger) {
	if m.Client != nil {
		if err := m.Client.Disconnect(context.Background()); err != nil {
			log.Error("Error while disconnecting from MongoDB", zap.Error(err))
		} else {
			log.Info("MongoDB disconnected successfully")
		}
	}
}
