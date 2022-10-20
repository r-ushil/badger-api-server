package server

import (
	"context"
	"net/http"

	"github.com/bufbuild/connect-go"

	activity_v1 "badger-api/gen/activity/v1"
	"badger-api/gen/activity/v1/activityv1connect"

	"badger-api/pkg/server"
)

type ActivityServer struct {
	ctx *server.ServerContext
}

func (s *ActivityServer) GetActivity(
	ctx context.Context,
	req *connect.Request[activity_v1.GetActivityRequest],
) (*connect.Response[activity_v1.GetActivityResponse], error) {
	res := connect.NewResponse(&activity_v1.GetActivityResponse{})
	res.Header().Set("Example-Version", "v1")

	return res, nil
}

func (s *ActivityServer) GetActivities(
	ctx context.Context,
	req *connect.Request[activity_v1.GetActivitiesRequest],
) (*connect.Response[activity_v1.GetActivitiesResponse], error) {
	res := connect.NewResponse(&activity_v1.GetActivitiesResponse{})
	res.Header().Set("Example-Version", "v1")

	return res, nil
}

func RegisterActivityService(mux *http.ServeMux, ctx *server.ServerContext) {
	server := &ActivityServer{
		ctx,
	}

	path, handler := activityv1connect.NewActivityServiceHandler(server)

	mux.Handle(path, handler)
}
