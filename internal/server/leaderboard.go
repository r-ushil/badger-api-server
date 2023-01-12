package server

import (
	"context"
	"net/http"

	"github.com/bufbuild/connect-go"

	leaderboardv1 "badger-api/gen/leaderboard/v1"
	"badger-api/gen/leaderboard/v1/leaderboardv1connect"
	"badger-api/pkg/auth"
	"badger-api/pkg/leaderboard"
	"badger-api/pkg/server"
)

type LeaderboardServer struct {
	ctx *server.ServerContext
}

func (s *LeaderboardServer) GetTopPlayers(
	ctx context.Context,
	req *connect.Request[leaderboardv1.GetTopPlayersRequest],
) (*connect.Response[leaderboardv1.GetTopPlayersResponse], error) {
	topPlayers := leaderboard.GetTopPlayers(s.ctx, int(req.Msg.Count))

	var players []*leaderboardv1.Player

	for _, topPlayer := range topPlayers {
		players = append(players, &leaderboardv1.Player{
			Id:    topPlayer.UserId,
			Name:  topPlayer.Name,
			Score: topPlayer.TotalScore,
		})
	}

	res := connect.NewResponse(&leaderboardv1.GetTopPlayersResponse{
		TopPlayers: players,
	})

	return res, nil
}

func (s *LeaderboardServer) GetMyScore(
	ctx context.Context,
	req *connect.Request[leaderboardv1.GetMyScoreRequest],
) (*connect.Response[leaderboardv1.GetMyScoreResponse], error) {
	authHeader := req.Header().Get("authorization")
	userId, err := auth.ParseAuthHeader(s.ctx, authHeader)

	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	score := leaderboard.GetPlayerScore(s.ctx, userId)

	res := connect.NewResponse(&leaderboardv1.GetMyScoreResponse{
		BattingScore:  score.BattingScore,
		CatchingScore: score.CatchingScore,
		BowlingScore:  score.BowlingScore,

		TotalBattingSubmissions:  score.TotalBattingSubmissions,
		TotalCatchingSubmissions: score.TotalCatchingSubmissions,
		TotalBowlingSubmissions:  score.TotalBowlingSubmissions,
	})

	return res, nil
}

func (s *LeaderboardServer) GetMyPublicName(
	ctx context.Context,
	req *connect.Request[leaderboardv1.GetMyPublicNameRequest],
) (*connect.Response[leaderboardv1.GetMyPublicNameResponse], error) {
	authHeader := req.Header().Get("authorization")
	userId, err := auth.ParseAuthHeader(s.ctx, authHeader)

	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	res := connect.NewResponse(&leaderboardv1.GetMyPublicNameResponse{
		Name: leaderboard.GetPlayerPublicName(s.ctx, userId),
	})

	return res, nil
}

func (s *LeaderboardServer) SetMyPublicName(
	ctx context.Context,
	req *connect.Request[leaderboardv1.SetMyPublicNameRequest],
) (*connect.Response[leaderboardv1.SetMyPublicNameResponse], error) {
	authHeader := req.Header().Get("authorization")
	userId, err := auth.ParseAuthHeader(s.ctx, authHeader)

	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	leaderboard.UpdatePlayerLeaderboardName(s.ctx, userId, req.Msg.Name)

	res := connect.NewResponse(&leaderboardv1.SetMyPublicNameResponse{})

	return res, nil
}

func RegisterLeaderboardService(mux *http.ServeMux, ctx *server.ServerContext) {
	server := &LeaderboardServer{
		ctx,
	}

	path, handler := leaderboardv1connect.NewLeaderboardServiceHandler(server)

	mux.Handle(path, handler)
}
