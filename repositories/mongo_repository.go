package repositories

import (
	"context"
	"errors"
	"math/rand"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	client *mongo.Client
	ctx    context.Context
}

type MatchMakingPlayer struct {
	PlayerId         string `bson:"player_id"`
	MatchMakingGroup string `bson:"match_making_group"`
}

func (repo MongoRepository) QueuePlayer(playerId string) error {

	filter := bson.D{
		{Key: "player_id", Value: playerId},
		{Key: "match_making_group", Value: ""},
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "player_id", Value: playerId},
			{Key: "match_making_group", Value: ""},
			{Key: "last_activity_on", Value: time.Now().Unix()},
		}},
	}

	t := true
	opt := options.UpdateOptions{
		Upsert: &t,
	}
	_, err := repo.client.Database("match_management").Collection("players").UpdateOne(repo.ctx, filter, update, &opt)

	return err
}

func (repo MongoRepository) TouchQueue(playerId string) error {

	filter := bson.D{
		{Key: "player_id", Value: playerId},
		{Key: "match_making_group", Value: ""},
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "last_activity_on", Value: time.Now().Unix()},
		}},
	}

	_, err := repo.client.Database("match_management").Collection("players").UpdateOne(repo.ctx, filter, update)

	return err
}

func (repo MongoRepository) UnqueuePlayer(playerId string) error {

	filter := bson.D{
		{Key: "player_id", Value: playerId},
		{Key: "match_making_group", Value: ""},
	}
	res, err := repo.client.Database("match_management").Collection("players").DeleteOne(repo.ctx, filter)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("cannot remove from queue")
	}

	return err
}

func (repo MongoRepository) GetMatchMakingGroup(matchMakingGroup string) ([]MatchMakingPlayer, error) {
	filter := bson.D{
		{Key: "match_making_group", Value: matchMakingGroup},
	}
	cursor, err := repo.client.Database("match_management").Collection("players").Find(repo.ctx, filter)

	if err != nil {
		return []MatchMakingPlayer{}, err
	}
	var results []MatchMakingPlayer
	if err = cursor.All(repo.ctx, &results); err != nil {
		return []MatchMakingPlayer{}, err
	}
	return results, nil
}

func (repo MongoRepository) SetRandomMatchMakingGroup() (string, error) {
	randomGroup := strconv.Itoa(int(rand.Int31()))

	filter := bson.D{
		{Key: "match_making_group", Value: ""},
		{Key: "last_activity_on", Value: bson.D{
			{Key: "$gte", Value: time.Now().Unix() - 60},
		}},
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "match_making_group", Value: randomGroup},
		}},
	}
	_, err := repo.client.Database("match_management").Collection("players").UpdateMany(repo.ctx, filter, update)
	return randomGroup, err
}

func (repo MongoRepository) DeleteMatchMakingGroup(matchMakingGroup string) error {
	filter := bson.D{
		{Key: "match_making_group", Value: matchMakingGroup},
	}
	_, err := repo.client.Database("match_management").Collection("players").DeleteMany(repo.ctx, filter)
	return err
}

func (repo MongoRepository) ClearPlayerMatchMakingGroup(playerId string) error {
	filter := bson.D{
		{Key: "player_id", Value: playerId},
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "match_making_group", Value: ""},
		}},
	}
	_, err := repo.client.Database("match_management").Collection("players").UpdateOne(repo.ctx, filter, update)
	return err
}
