package activity

import (
	"badger-api/pkg/server"
	"context"
	"time"

	"cloud.google.com/go/civil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/genproto/googleapis/type/datetime"
)

type mongoActivityT struct {
	Id        string              `bson:"_id"`
	VideoUrl  string              `bson:"name"`
	Score     uint32              `bson:"description"`
	Timestamp primitive.Timestamp `bson:"description"`
}

type Activity struct {
	Id        string
	VideoUrl  string
	Score     uint32
	Timestamp time.Time
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
	return a.Timestamp
}

func (a *Activity) GetTimestampGoogleFormat() datetime.DateTime {
	civilDateTime := civil.DateTimeOf(a.GetTimestamp())
	return datetime.DateTime{
		Year:    int32(civilDateTime.Date.Year),
		Month:   int32(civilDateTime.Date.Month),
		Day:     int32(civilDateTime.Date.Day),
		Hours:   int32(civilDateTime.Time.Hour),
		Minutes: int32(civilDateTime.Time.Minute),
		Seconds: int32(civilDateTime.Time.Second),
		Nanos:   int32(civilDateTime.Time.Nanosecond),
	}
}

func GetActivities(s *server.ServerContext) []Activity {
	col := s.GetCollection("activities")

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

func GetActivity(s *server.ServerContext, hexId string) (*Activity, error) {
	col := s.GetCollection("activities")

	objectId, idErr := primitive.ObjectIDFromHex(hexId)

	if idErr != nil {
		panic(idErr)
	}

	query := bson.D{{Key: "_id", Value: objectId}}

	var mongoActivity mongoActivityT
	err := col.FindOne(s.GetMongoContext(), query).Decode(&mongoActivity)

	activity := Activity{
		Id:        mongoActivity.Id,
		VideoUrl:  mongoActivity.VideoUrl,
		Score:     mongoActivity.Score,
		Timestamp: time.Unix(int64(mongoActivity.Timestamp.T), 0),
	}

	if err == mongo.ErrNoDocuments {
		return nil, ErrNotFound
	}

	if err != nil {
		panic(err)
	}

	return &activity, nil
}
