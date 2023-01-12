package drill

import (
	"time"

	"badger-api/pkg/server"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const BattingDrillSubmissionCollection = "submissions_batting_drill"

type BattingDrillSubmissionDoc struct {
	Id               string           `bson:"_id,omitempty"`
	UserId           string           `bson:"user_id"`
	VideoObjectName  string           `bson:"video_obj_name"`
	Timestamp        time.Time        `bson:"timestamp"`
	ProcessingStatus ProcessingStatus `bson:"processing_status"`
	Score            float64          `bson:"score"`
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

func RegisterBattingDrillResults(s *server.ServerContext, submissionId string, score float64) {
	col := s.GetCollection(BattingDrillSubmissionCollection)

	col.UpdateByID(s.GetMongoContext(), submissionId, BattingDrillSubmissionDoc{
		Score: score,
	})
}
