package server

import (
	"context"
	"net/http"

	"github.com/bufbuild/connect-go"

	leaderboardv1 "badger-api/gen/leaderboard/v1"
	"badger-api/gen/leaderboard/v1/leaderboardv1connect"
	"badger-api/pkg/auth"
	"badger-api/pkg/drill"
	"badger-api/pkg/server"
)

type LeaderboardServer struct {
	ctx *server.ServerContext
}

func (s *LeaderboardServer) GetTopPlayers(
	ctx context.Context,
	req *connect.Request[leaderboardv1.GetTopPlayersRequest],
) (*connect.Response[leaderboardv1.GetTopPlayersResponse], error) {
	res := connect.NewResponse(&leaderboardv1.GetTopPlayersResponse{
		TopPlayers: []*leaderboardv1.Player{},
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

	res := connect.NewResponse(&leaderboardv1.GetMyScoreResponse{
		BattingScore:  drill.ComputeBattingScoreForUser(s.ctx, userId),
		CatchingScore: 0,
		BowlingScore:  0,

		TotalBattingSubmissions:  drill.CountBattingSubmissionsByUser(s.ctx, userId),
		TotalCatchingSubmissions: 0,
		TotalBowlingSubmissions:  0,
	})

	return res, nil
}

func (s *LeaderboardServer) GetMyPublicName(
	ctx context.Context,
	req *connect.Request[leaderboardv1.GetMyPublicNameRequest],
) (*connect.Response[leaderboardv1.GetMyPublicNameResponse], error) {
	res := connect.NewResponse(&leaderboardv1.GetMyPublicNameResponse{
		Name: "Mock Name",
	})

	return res, nil
}

func (s *LeaderboardServer) SetMyPublicName(
	ctx context.Context,
	req *connect.Request[leaderboardv1.SetMyPublicNameRequest],
) (*connect.Response[leaderboardv1.SetMyPublicNameResponse], error) {
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
