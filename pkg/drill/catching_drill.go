package drill

import (
	"context"
	"time"

	"badger-api/pkg/server"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const CatchingDrillSubmissionCollection = "submissions_catching_drill"

type CatchingDrillSubmissionDoc struct {
	Id               string           `bson:"_id,omitempty"`
	UserId           string           `bson:"user_id"`
	VideoObjectName  string           `bson:"video_obj_name"`
	Timestamp        time.Time        `bson:"timestamp"`
	ProcessingStatus ProcessingStatus `bson:"processing_status"`
	Score            uint32           `bson:"score,truncate"`
}

func SubmitCatchingDrill(s *server.ServerContext, videoObjectName string, userId string) string {
	col := s.GetCollection(CatchingDrillSubmissionCollection)

	data := CatchingDrillSubmissionDoc{
		UserId:           userId,
		VideoObjectName:  videoObjectName,
		Timestamp:        time.Now(),
		ProcessingStatus: PROCESSING_STATUS_PENDING,
	}

	result, err := col.InsertOne(s.GetMongoContext(), data)

	if err != nil {
		panic(err)
	}

	return result.InsertedID.(primitive.ObjectID).Hex()
}

func RegisterCatchingDrillResults(s *server.ServerContext, submissionId string, score uint32) {
	col := s.GetCollection(CatchingDrillSubmissionCollection)

	submissionObjId, err := primitive.ObjectIDFromHex(submissionId)

	if err != nil {
		panic(err)
	}

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "score", Value: uint32(score)}}}}
	col.UpdateByID(s.GetMongoContext(), submissionObjId, update)
}

func CountCatchingSubmissionsByUser(s *server.ServerContext, userId string) uint32 {
	col := s.GetCollection(CatchingDrillSubmissionCollection)

	filter := bson.D{{
		Key: "user_id",
		Value: bson.D{{
			Key:   "$eq",
			Value: userId,
		}},
	}}

	count, err := col.CountDocuments(s.GetMongoContext(), filter)

	if err != nil {
		panic(err)
	}

	if count < 0 {
		return 0
	}

	return uint32(count)
}

func ComputeCatchingScoreForUser(s *server.ServerContext, userId string) uint32 {
	col := s.GetCollection(CatchingDrillSubmissionCollection)

	match_stage := bson.D{{
		Key: "$match",
		Value: bson.D{{
			Key:   "user_id",
			Value: userId,
		}},
	}}
	group_stage := bson.D{{
		Key: "$group",
		Value: bson.D{
			{Key: "_id", Value: "$user_id"},
			{Key: "score", Value: bson.D{{Key: "$avg", Value: "$score"}}},
		},
	}}

	cursor, err := col.Aggregate(s.GetMongoContext(), mongo.Pipeline{
		match_stage,
		group_stage,
	})

	if err != nil {
		panic(err)
	}

	var results []struct {
		Id    string `bson:"_id"`
		Score uint32 `bson:"score,truncate"`
	}

	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	if len(results) == 0 {
		return 0
	}

	agg := results[0]

	return agg.Score
}
