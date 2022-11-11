package drill_submission

import (
	"badger-api/pkg/server"
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type DrillSubmission struct {
	DrillSubmissionId string                 `bson:"_id"`
	UserId            string                 `bson:"user_id"`
	DrillId           string                 `bson:"drill_id"`
	BucketUrl         string                 `bson:"bucket_url"`
	Timestamp         *timestamppb.Timestamp `bson:"timestamp"`
	ProcessingStatus  string                 `bson:"processing_status"`
	DrillScore        uint32                 `bson:"drill_score"`
}

func (d *DrillSubmission) GetId() string {
	return d.DrillSubmissionId
}

func (d *DrillSubmission) GetUserId() string {
	return d.UserId
}

func (d *DrillSubmission) GetDrillId() string {
	return d.DrillId
}

func (d *DrillSubmission) GetBucketUrl() string {
	return d.BucketUrl
}

func (d *DrillSubmission) GetTimestamp() *timestamppb.Timestamp {
	return d.Timestamp
}

func (d *DrillSubmission) GetProcessingStatus() string {
	return d.ProcessingStatus
}

func (d *DrillSubmission) GetDrillScore() uint32 {
	return d.DrillScore
}

func GetDrillSubmissions(s *server.ServerContext) []DrillSubmission {
	col := s.GetCollection("drill_submissions")

	cur, err_find := col.Find(s.GetMongoContext(), bson.D{})

	if err_find != nil {
		panic(err_find)
	}

	var drill_submissions []DrillSubmission
	err_cur := cur.All(context.TODO(), &drill_submissions)

	if err_cur != nil {
		panic(err_cur)
	}

	return drill_submissions
}

func GetDrillSubmission(s *server.ServerContext, hexId string) (*DrillSubmission, error) {
	log.Println("Getting drill submissino collection. ")
	col := s.GetCollection("drill_submissions")
	log.Println("Getting drill submission collection done. ")

	objectId, idErr := primitive.ObjectIDFromHex(hexId)

	if idErr != nil {
		panic(idErr)
	}

	query := bson.D{{Key: "_id", Value: objectId}}

	var drill_submission DrillSubmission
	log.Println("Getting drill document. ")
	err := col.FindOne(s.GetMongoContext(), query).Decode(&drill_submission)
	log.Println("Getting drill document done. ")

	if err == mongo.ErrNoDocuments {
		return nil, errors.New("ErrNotFound")
	}

	if err != nil {
		panic(err)
	}

	log.Println("All good, returning drill. ")
	log.Println(drill_submission)

	return &drill_submission, nil
}

func GetUserDrillSubmissions(s *server.ServerContext, hexUserId string) []DrillSubmission {
	col := s.GetCollection("drill_submissions")

	objectId, idErr := primitive.ObjectIDFromHex(hexUserId)

	if idErr != nil {
		panic(idErr)
	}

	query := bson.D{{Key: "_id", Value: objectId}}

	cur, err_find := col.Find(s.GetMongoContext(), query)

	if err_find != nil {
		panic(err_find)
	}

	var drill_submissions []DrillSubmission
	err_cur := cur.All(context.TODO(), &drill_submissions)

	if err_cur != nil {
		panic(err_cur)
	}

	return drill_submissions
}
