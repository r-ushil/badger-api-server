package server

import (
	"context"
	"net/http"

	"github.com/bufbuild/connect-go"

	drill_v1 "badger-api/gen/drill/v1"
	"badger-api/gen/drill/v1/drillv1connect"

	"badger-api/pkg/drill"
	"badger-api/pkg/server"
)

type DrillServer struct {
	ctx *server.ServerContext
}

func (s *DrillServer) GetDrill(
	ctx context.Context,
	req *connect.Request[drill_v1.GetDrillRequest],
) (*connect.Response[drill_v1.GetDrillResponse], error) {
	d, err := drill.GetDrill(s.ctx, req.Msg.DrillId)

	if err == drill.ErrNotFound {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	res := connect.NewResponse(&drill_v1.GetDrillResponse{
		Drill: &drill_v1.Drill{
			DrillId:          d.GetId(),
			DrillName:        d.GetName(),
			DrillDescription: d.GetDescription(),
		},
	})

	return res, nil
}

func (s *DrillServer) GetDrillInstructions(
	ctx context.Context,
	req *connect.Request[drill_v1.GetDrillInstructionsRequest],
) (*connect.Response[drill_v1.GetDrillInstructionsResponse], error) {
	res := connect.NewResponse(&drill_v1.GetDrillInstructionsResponse{})

	return res, nil
}

func (s *DrillServer) GetDrills(
	ctx context.Context,
	req *connect.Request[drill_v1.GetDrillsRequest],
) (*connect.Response[drill_v1.GetDrillsResponse], error) {
	ds := drill.GetDrills(s.ctx)

	drills := make([]*drill_v1.DrillOverview, 0, len(ds))

	for _, drill := range ds {
		drills = append(drills, &drill_v1.DrillOverview{
			DrillId:          drill.GetId(),
			DrillName:        drill.GetName(),
			DrillDescription: drill.GetDescription(),
		})
	}

	res := connect.NewResponse(&drill_v1.GetDrillsResponse{
		Drills: drills,
	})

	return res, nil
}

func RegisterDrillService(mux *http.ServeMux, ctx *server.ServerContext) {
	server := &DrillServer{
		ctx,
	}

	path, handler := drillv1connect.NewDrillServiceHandler(server)

	mux.Handle(path, handler)
}
