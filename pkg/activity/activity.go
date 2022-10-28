package activity

import (
	"badger-api/pkg/server"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Activity struct {
	Id        string              `bson:"_id"`
	VideoUrl  string              `bson:"name"`
	Score     uint32              `bson:"description"`
	Timestamp primitive.Timestamp `bson:"description"`
}

func (a *Activity) GetId() string {
	return a.Id
}

func (a *Activity) GetVideoUrl() string {
	return a.VideoUrl
}

func (a *Activity) GetScore() uint32 {
	return a.Score
}

func (a *Activity) GetTimestamp() time.Time {
	return time.Unix(int64(a.Timestamp.T), 0)
}

func GetActivities(s *server.ServerContext) []Activity {
	client := s.GetMongoDbClient()

	db := client.Database("badger_db")
	col := db.Collection("activities")

	cur, err_find := col.Find(s.GetMongoContext(), bson.D{})

	if err_find != nil {
		panic(err_find)
	}

	var activities []Activity
	err_cur := cur.All(context.TODO(), &activities)

	if err_cur != nil {
		panic(err_cur)
	}

	return activities
}

var ErrNotFound error

func GetActivity(s *server.ServerContext, id string) (*Activity, error) {
	client := s.GetMongoDbClient()

	db := client.Database("badger_db")
	col := db.Collection("activities")

	query := bson.D{{Key: "_id", Value: id}}

	var activity Activity
	err := col.FindOne(s.GetMongoContext(), query).Decode(&activity)

	if err == mongo.ErrNoDocuments {
		return nil, ErrNotFound
	}

	if err != nil {
		panic(err)
	}

	return &activity, nil
}
