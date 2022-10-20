package server

import (
	"context"
	"net/http"

	"github.com/bufbuild/connect-go"

	drill_v1 "badger-api/gen/drill/v1"
	"badger-api/gen/drill/v1/drillv1connect"

	"badger-api/pkg/server"
)

type DrillServer struct {
	ctx *server.ServerContext
}

func (s *DrillServer) GetDrill(
	ctx context.Context,
	req *connect.Request[drill_v1.GetDrillRequest],
) (*connect.Response[drill_v1.GetDrillResponse], error) {
	res := connect.NewResponse(&drill_v1.GetDrillResponse{})
	res.Header().Set("Example-Version", "v1")

	return res, nil
}

func (s *DrillServer) GetDrillInstructions(
	ctx context.Context,
	req *connect.Request[drill_v1.GetDrillInstructionsRequest],
) (*connect.Response[drill_v1.GetDrillInstructionsResponse], error) {
	res := connect.NewResponse(&drill_v1.GetDrillInstructionsResponse{})
	res.Header().Set("Example-Version", "v1")

	return res, nil
}

func (s *DrillServer) GetDrills(
	ctx context.Context,
	req *connect.Request[drill_v1.GetDrillsRequest],
) (*connect.Response[drill_v1.GetDrillsResponse], error) {
	res := connect.NewResponse(&drill_v1.GetDrillsResponse{})
	res.Header().Set("Example-Version", "v1")

	return res, nil
}

func RegisterDrillService(mux *http.ServeMux, ctx *server.ServerContext) {
	server := &DrillServer{
		ctx,
	}

	path, handler := drillv1connect.NewDrillServiceHandler(server)

	mux.Handle(path, handler)
}
