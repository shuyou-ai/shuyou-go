package database

import (
	"context"
	"fmt"

	"github.com/shuyou-ai/shuyou-go/internal/config"
	"github.com/shuyou-ai/shuyou-go/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type Client struct {
	client   *mongo.Client
	database *mongo.Database
}

func New(cfg config.DatabaseConfig) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectTimeout)
	defer cancel()

	clientOpts := options.Client().
		ApplyURI(cfg.URI).
		SetMaxPoolSize(cfg.MaxPoolSize)

	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return nil, fmt.Errorf("connect mongodb: %w", err)
	}

	c := &Client{
		client:   client,
		database: client.Database(cfg.Database),
	}

	if err := c.Ping(ctx); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, err
	}

	return c, nil
}

func (c *Client) DB() *mongo.Database {
	return c.database
}

func (c *Client) Ping(ctx context.Context) error {
	return c.client.Ping(ctx, readpref.Primary())
}

func (c *Client) Close(ctx context.Context) error {
	if c.client == nil {
		return nil
	}
	return c.client.Disconnect(ctx)
}

func (c *Client) EnsureIndexes(ctx context.Context) error {
	coll := c.database.Collection(model.UserCollectionName)
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	if _, err := coll.Indexes().CreateMany(ctx, indexes); err != nil {
		return fmt.Errorf("create user indexes: %w", err)
	}

	return nil
}
