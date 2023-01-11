package drill_submission

import (
	"badger-api/pkg/server"
	"context"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/civil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/genproto/googleapis/type/datetime"

	"io/ioutil"
	"net/http"

	drill_submission_v1 "badger-api/gen/drill_submission/v1"
)

type DrillSubmission struct {
	DrillSubmissionId string    `bson:"_id,omitempty"`
	UserId            string    `bson:"user_id"`
	DrillId           string    `bson:"drill_id"`
	BucketUrl         string    `bson:"bucket_url"`
	Timestamp         time.Time `bson:"timestamp"`
	ProcessingStatus  string    `bson:"processing_status"`
	DrillScore        uint32    `bson:"drill_score"`
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

func (d *DrillSubmission) GetTimestamp() time.Time {
	return d.Timestamp
}

func (d *DrillSubmission) GetTimestampGoogleFormat() datetime.DateTime {
	civilDateTime := civil.DateTimeOf(d.GetTimestamp())
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

func (d *DrillSubmission) GetProcessingStatus() string {
	return d.ProcessingStatus
}

func (d *DrillSubmission) GetDrillScore() uint32 {
	return d.DrillScore
}

func InsertDrillSubmission(s *server.ServerContext, drill_submission *drill_submission_v1.DrillSubmission) string {
	col := s.GetCollection("drill_submissions")

	data := DrillSubmission{
		UserId:    drill_submission.UserId,
		DrillId:   drill_submission.DrillId,
		BucketUrl: drill_submission.BucketUrl,
		Timestamp: time.Date(
			int(drill_submission.Timestamp.Year),
			time.Month(drill_submission.Timestamp.Month),
			int(drill_submission.Timestamp.Day),
			int(drill_submission.Timestamp.Hours),
			int(drill_submission.Timestamp.Minutes),
			int(drill_submission.Timestamp.Seconds),
			int(drill_submission.Timestamp.Nanos),
			time.UTC),
		ProcessingStatus: drill_submission.ProcessingStatus,
		DrillScore:       drill_submission.DrillScore,
	}

	result, err := col.InsertOne(s.GetMongoContext(), data)

	if err != nil {
		panic(err)
	}
	return result.InsertedID.(primitive.ObjectID).Hex()
}

func GetUserScores(s *server.ServerContext, userId string) (float32, float32) {

	var userSubmissions = GetUserDrillSubmissions(s, userId)
	var coverDriveScore = 0
	var coverDrives = 0
	for _, submission := range userSubmissions {
		if submission.DrillId == "6352414e50c7d61db5d52861" {
			coverDrives++
			coverDriveScore += int(submission.DrillScore)
		} 
	//TODO: katchet board
	}
	return float32(coverDriveScore/coverDrives), 90.00
}

func ProcessDrillSubmission(s *server.ServerContext, submissionId string, bucketUrl string) (uint32, string, string) {

	var requestUrl = "https://badger-cv-microservice-6la2hzpokq-ew.a.run.app/cover-drive-drill?video_object_name=" + bucketUrl
	response, get_err := http.Get(requestUrl)

	if get_err != nil {
		panic(get_err)
	}

	responseData, io_err := ioutil.ReadAll(response.Body)
	if io_err != nil {
		panic(io_err)
	}

	split_response_data := strings.Split(string(responseData), ",")
	score, atoi_err := strconv.Atoi(split_response_data[0])
	if atoi_err != nil {
		panic(atoi_err)
	}
	advice1 := split_response_data[1]
	advice2 := split_response_data[2]

	col := s.GetCollection("drill_submissions")
	id, id_err := primitive.ObjectIDFromHex(submissionId)
	if id_err != nil {
		panic(id_err)
	}

	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "drill_score", Value: uint32(score)}, {Key: "processing_status", Value: "Done"}}}}
	_, update_err := col.UpdateOne(context.TODO(), filter, update)
	if update_err != nil {
		panic(update_err)
	}

	return uint32(score), advice1, advice2
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

func GetUserDrillSubmissions(s *server.ServerContext, userId string) []DrillSubmission {
	col := s.GetCollection("drill_submissions")

	query := bson.D{{Key: "user_id", Value: userId}}

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
