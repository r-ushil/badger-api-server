package server

import (
	"context"
	"net/http"

	"github.com/bufbuild/connect-go"

	batting_drillv1 "badger-api/gen/batting_drill/v1"
	"badger-api/gen/batting_drill/v1/batting_drillv1connect"
	"badger-api/pkg/auth"
	"badger-api/pkg/drill"
	"badger-api/pkg/server"
)

type BattingDrillServer struct {
	ctx *server.ServerContext
}

func (s *BattingDrillServer) SubmitBattingDrill(
	ctx context.Context,
	req *connect.Request[batting_drillv1.SubmitBattingDrillRequest],
) (*connect.Response[batting_drillv1.SubmitBattingDrillResponse], error) {
	authHeader := req.Header().Get("authorization")

	userId, err := auth.ParseAuthHeader(s.ctx, authHeader)

	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	submissionId := drill.SubmitBattingDrill(s.ctx, req.Msg.VideoObjectName, userId)

	res := connect.NewResponse(&batting_drillv1.SubmitBattingDrillResponse{
		SubmissionId: submissionId,
	})

	return res, nil
}

func (s *BattingDrillServer) OnBattingDrillProcessingComplete(
	ctx context.Context,
	req *connect.Request[batting_drillv1.OnBattingDrillProcessingCompleteRequest],
) (*connect.Response[batting_drillv1.OnBattingDrillProcessingCompleteResponse], error) {
	res := connect.NewResponse(&batting_drillv1.OnBattingDrillProcessingCompleteResponse{})

	return res, nil
}

func RegisterBattingDrillService(mux *http.ServeMux, ctx *server.ServerContext) {
	server := &BattingDrillServer{
		ctx,
	}

	path, handler := batting_drillv1connect.NewBattingDrillServiceHandler(server)

	mux.Handle(path, handler)
}
