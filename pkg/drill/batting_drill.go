package drill

import (
	"context"
	"time"

	"badger-api/pkg/server"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const BattingDrillSubmissionCollection = "submissions_batting_drill"

type BattingDrillSubmissionDoc struct {
	Id               string           `bson:"_id,omitempty"`
	UserId           string           `bson:"user_id"`
	VideoObjectName  string           `bson:"video_obj_name"`
	Timestamp        time.Time        `bson:"timestamp"`
	ProcessingStatus ProcessingStatus `bson:"processing_status"`
	Score            uint32           `bson:"score"`
}

func SubmitBattingDrill(s *server.ServerContext, videoObjectName string, userId string) string {
	col := s.GetCollection(BattingDrillSubmissionCollection)

	data := BattingDrillSubmissionDoc{
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

func RegisterBattingDrillResults(s *server.ServerContext, submissionId string, score uint32) {
	col := s.GetCollection(BattingDrillSubmissionCollection)

	col.UpdateByID(s.GetMongoContext(), submissionId, BattingDrillSubmissionDoc{
		Score: score,
	})
}

func CountBattingSubmissionsByUser(s *server.ServerContext, userId string) uint32 {
	col := s.GetCollection(BattingDrillSubmissionCollection)

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

func ComputeBattingScoreForUser(s *server.ServerContext, userId string) uint32 {
	col := s.GetCollection(BattingDrillSubmissionCollection)

	match_stage := bson.D{{"$match", bson.D{{"user_id", userId}}}}
	group_stage := bson.D{
		{"$group",
			bson.D{
				{"_id", "$user_id"},
				{"score", bson.D{{"$avg", "$score"}}},
			},
		},
	}

	cursor, err := col.Aggregate(s.GetMongoContext(), mongo.Pipeline{
		match_stage,
		group_stage,
	})

	if err != nil {
		panic(err)
	}

	var results []struct {
		Id    string `bson:"_id"`
		Score uint32 `bson:"score"`
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
