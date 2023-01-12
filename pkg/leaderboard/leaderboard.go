package leaderboard

import (
	"badger-api/pkg/drill"
	"badger-api/pkg/server"
)

type LeaderboardPlayer struct {
	userId string
	name   string
	score  PlayerScore
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
		CatchingScore: 0,
		BowlingScore:  0,

		TotalBattingSubmissions:  drill.CountBattingSubmissionsByUser(s, userId),
		TotalCatchingSubmissions: 0,
		TotalBowlingSubmissions:  0,
	}
}

func GetTopPlayers(s *server.ServerContext, count uint64) []LeaderboardPlayer {
	return []LeaderboardPlayer{
		{
			userId: "user-one",
			name:   "Anon",
			score:  GetPlayerScore(s, "user-one"),
		},
	}
}
