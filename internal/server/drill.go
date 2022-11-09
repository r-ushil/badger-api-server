package server

import (
	"context"
	"log"
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
	log.Println("Getting drill with ID:", req.Msg.DrillId)

	d, err := drill.GetDrill(s.ctx, req.Msg.DrillId)

	if err != nil {
		// TODO: Handle properly
		log.Println(err)
		return nil, connect.NewError(connect.CodeUnimplemented, err)
	}

	res := connect.NewResponse(&drill_v1.GetDrillResponse{
		Drill: &drill_v1.Drill{
			DrillId:          d.GetId(),
			DrillName:        d.GetName(),
			DrillDescription: d.GetDescription(),
			Instructions:     d.GetInstructions(),
			ThumbnailUrl:     d.GetThumbnailUrl(),
			Skills:           d.GetSkills(),
			VideoUrl:         d.GetVideoUrl(),
		},
	})

	return res, nil
}

func (s *DrillServer) GetDrills(
	ctx context.Context,
	req *connect.Request[drill_v1.GetDrillsRequest],
) (*connect.Response[drill_v1.GetDrillsResponse], error) {
	ds := drill.GetDrills(s.ctx)

	drills := make([]*drill_v1.Drill, 0, len(ds))

	for _, drill := range ds {
		drills = append(drills, &drill_v1.Drill{
			DrillId:          drill.GetId(),
			DrillName:        drill.GetName(),
			DrillDescription: drill.GetDescription(),
			Instructions:     drill.GetInstructions(),
			ThumbnailUrl:     drill.GetThumbnailUrl(),
			Skills:           drill.GetSkills(),
			VideoUrl:         drill.GetVideoUrl(),
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
