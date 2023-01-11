package server

import (
	"context"
	"net/http"

	"github.com/bufbuild/connect-go"

	catching_drillv1 "badger-api/gen/catching_drill/v1"
	"badger-api/gen/catching_drill/v1/catching_drillv1connect"
	"badger-api/pkg/server"
)

type CatchingDrillServer struct {
	ctx *server.ServerContext
}

func (s *CatchingDrillServer) SubmitCatchingDrill(
	ctx context.Context,
	req *connect.Request[catching_drillv1.SubmitCatchingDrillRequest],
) (*connect.Response[catching_drillv1.SubmitCatchingDrillResponse], error) {
	res := connect.NewResponse(&catching_drillv1.SubmitCatchingDrillResponse{
		SubmissionId: "Mock ID",
	})

	return res, nil
}

func (s *CatchingDrillServer) OnCatchingDrillProcessingComplete(
	ctx context.Context,
	req *connect.Request[catching_drillv1.OnCatchingDrillProcessingCompleteRequest],
) (*connect.Response[catching_drillv1.OnCatchingDrillProcessingCompleteResponse], error) {
	res := connect.NewResponse(&catching_drillv1.OnCatchingDrillProcessingCompleteResponse{})

	return res, nil
}

func RegisterCatchingDrillService(mux *http.ServeMux, ctx *server.ServerContext) {
	server := &CatchingDrillServer{
		ctx,
	}

	path, handler := catching_drillv1connect.NewCatchingDrillServiceHandler(server)

	mux.Handle(path, handler)
}
