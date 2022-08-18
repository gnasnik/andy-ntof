package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

func newDB() (*mongo.Client, error) {
	var err error
	ctx := context.Background()

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://admin:Qweasdzxc123@127.0.0.1:27037/ntof?authsource=admin"
	}

	//cs, err := connstring.Parse(uri)
	//if err != nil {
	//	return err
	//}
	//database := cs.Database

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func Players(c *mongo.Client) *mongo.Collection {
	return c.Database("ntof").Collection("players")
}

func Stats(c *mongo.Client) *mongo.Collection {
	return c.Database("ntof").Collection("stats")
}

func UpsertPlayer(ctx context.Context, c *mongo.Collection, players []*player) error {
	for _, player := range players {
		opt := options.FindOneAndUpdate().SetUpsert(true)
		f := bson.D{{"id", player.Id}, {"date", player.Date}}
		u := bson.D{{"$set", player}}
		player.CreatedAt = time.Now()
		rs := c.FindOneAndUpdate(ctx, f, u, opt)
		if rs.Err() != nil && rs.Err() != mongo.ErrNoDocuments {
			fmt.Println(rs.Err())
			return rs.Err()
		}
	}
	return nil
}

type stats struct {
	InitialCap float64
	MarketCap  float64
	Players    int
	GoodCount  int
	Date       string
	CreatedAt  time.Time
}

func UpsertStats(ctx context.Context, c *mongo.Collection, stats *stats) error {
	opt := options.FindOneAndUpdate().SetUpsert(true)
	f := bson.D{{"date", stats.Date}}
	u := bson.D{{"$set", stats}}
	rs := c.FindOneAndUpdate(ctx, f, u, opt)
	if rs.Err() != nil {
		return rs.Err()
	}
	return nil
}

func GetPlayers(ctx context.Context, c *mongo.Collection) ([]*player, error) {
	opt := options.Find().SetSort(bson.D{{"date", -1}, {"asset", -1}}).SetLimit(20)
	rs, err := c.Find(ctx, bson.D{}, opt)
	if err != nil {
		return nil, err
	}

	var players []*player
	if err := rs.All(ctx, &players); err != nil {
		return nil, err
	}

	return players, nil
}

func GetStats(ctx context.Context, c *mongo.Collection) ([]*stats, error) {
	opt := options.Find().SetSort(bson.D{{"asset", -1}}).SetLimit(20)
	rs, err := c.Find(ctx, bson.D{}, opt)
	if err != nil {
		return nil, err
	}

	var stats []*stats
	if err := rs.All(ctx, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

func GetPlayerListByName(ctx context.Context, name string, c *mongo.Collection) ([]*player, error) {
	opt := options.Find().SetSort(bson.D{{"date", -1}, {"asset", -1}}).SetLimit(20)
	rs, err := c.Find(ctx, bson.D{{"name", name}}, opt)
	if err != nil {
		return nil, err
	}

	var players []*player
	if err := rs.All(ctx, &players); err != nil {
		return nil, err
	}

	return players, nil
}
