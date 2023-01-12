package leaderboard

import (
	"badger-api/pkg/drill"
	"badger-api/pkg/server"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const LeaderboardOverallScoreCollection = "leaderboard_scores"

type LeaderboardPlayerDoc struct {
	Id     string `bson:"_id"`
	UserId string `bson:"user_id"`
	Name   string `bson:"name"`
	Score  uint32 `bson:"score,truncate"`
}

type LeaderboardPlayer struct {
	UserId     string
	Name       string
	TotalScore uint32
	Breakdown  PlayerScore
}

type PlayerScore struct {
	BattingScore  uint32
	CatchingScore uint32
	BowlingScore  uint32

	TotalBattingSubmissions  uint32
	TotalCatchingSubmissions uint32
	TotalBowlingSubmissions  uint32
}

func GetPlayerScore(s *server.ServerContext, userId string) PlayerScore {
	return PlayerScore{
		BattingScore:  drill.ComputeBattingScoreForUser(s, userId),
		CatchingScore: drill.ComputeCatchingScoreForUser(s, userId),
		BowlingScore:  0,

		TotalBattingSubmissions:  drill.CountBattingSubmissionsByUser(s, userId),
		TotalCatchingSubmissions: drill.CountCatchingSubmissionsByUser(s, userId),
		TotalBowlingSubmissions:  0,
	}
}

func GetTopPlayers(s *server.ServerContext, count int) []LeaderboardPlayer {
	col := s.GetCollection(LeaderboardOverallScoreCollection)

	filter := bson.D{}
	opts := options.Find().SetSort(bson.D{{Key: "score", Value: -1}})

	cursor, err := col.Find(s.GetMongoContext(), filter, opts)

	var results []LeaderboardPlayerDoc
	if err = cursor.All(s.GetMongoContext(), &results); err != nil {
		panic(err)
	}

	var leaderboard []LeaderboardPlayer
	for i, result := range results {
		if i == count {
			break
		}

		leaderboard = append(leaderboard, LeaderboardPlayer{
			UserId:     result.UserId,
			Name:       result.Name,
			TotalScore: result.Score,
			Breakdown:  GetPlayerScore(s, result.UserId),
		})
	}

	return leaderboard
}

func GetPlayerPublicName(s *server.ServerContext, userId string) string {
	col := s.GetCollection(LeaderboardOverallScoreCollection)

	var result LeaderboardPlayerDoc
	err := col.FindOne(s.GetMongoContext(), bson.M{"user_id": userId}).Decode(&result)

	if err != nil {
		return "Anon"
	}

	return result.Name
}

func UpdatePlayerLeaderboardScore(s *server.ServerContext, userId string) {
	score := GetPlayerScore(s, userId)
	overallScore := (score.BattingScore + score.BowlingScore + score.CatchingScore) / 3

	col := s.GetCollection(LeaderboardOverallScoreCollection)

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "score", Value: uint32(overallScore)}}}}
	opts := options.Update().SetUpsert(true)

	col.UpdateOne(s.GetMongoContext(), bson.M{"user_id": userId}, update, opts)
}

func UpdatePlayerLeaderboardName(s *server.ServerContext, userId, name string) {
	col := s.GetCollection(LeaderboardOverallScoreCollection)

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "name", Value: name}}}}
	opts := options.Update().SetUpsert(true)

	col.UpdateByID(s.GetMongoContext(), bson.M{"user_id": userId}, update, opts)
}
