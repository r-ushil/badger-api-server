package server

import (
	"context"
	"log"
	"net/http"

	"github.com/bufbuild/connect-go"

	activity_v1 "badger-api/gen/activity/v1"
	"badger-api/gen/activity/v1/activityv1connect"

	"badger-api/pkg/activity"
	"badger-api/pkg/server"
)

type ActivityServer struct {
	ctx *server.ServerContext
}

func (s *ActivityServer) GetActivity(
	ctx context.Context,
	req *connect.Request[activity_v1.GetActivityRequest],
) (*connect.Response[activity_v1.GetActivityResponse], error) {

	a, err := activity.GetActivity(s.ctx, req.Msg.ActivityId)

	if err != nil {
		// TODO: Handle properly
		log.Println(err)
		return nil, connect.NewError(connect.CodeUnimplemented, err)
	}

	timestampGoogleFormat := a.GetTimestampGoogleFormat()
	res := connect.NewResponse(&activity_v1.GetActivityResponse{
		Activity: &activity_v1.Activity{
			ActivityId:           a.GetId(),
			ActivityThumbnailUrl: a.GetThumbnailUrl(),
			ActivityScore:        a.GetScore(),
			ActivityTimestamp:    &timestampGoogleFormat,
		},
	})

	return res, nil
}

func (s *ActivityServer) GetActivities(
	ctx context.Context,
	req *connect.Request[activity_v1.GetActivitiesRequest],
) (*connect.Response[activity_v1.GetActivitiesResponse], error) {
	ds := activity.GetActivities(s.ctx)

	activities := make([]*activity_v1.ActivityOverview, 0, len(ds))

	for _, a := range ds {

		timestampGoogleFormat := a.GetTimestampGoogleFormat()
		activities = append(activities, &activity_v1.ActivityOverview{
			ActivityId:           a.GetId(),
			ActivityThumbnailUrl: a.GetThumbnailUrl(),
			ActivityScore:        a.GetScore(),
			ActivityTimestamp:    &timestampGoogleFormat,
		})
	}

	res := connect.NewResponse(&activity_v1.GetActivitiesResponse{
		Activities: activities,
	})

	return res, nil
}

func RegisterActivityService(mux *http.ServeMux, ctx *server.ServerContext) {
	server := &ActivityServer{
		ctx,
	}

	path, handler := activityv1connect.NewActivityServiceHandler(server)

	mux.Handle(path, handler)
}
