package main

import (
	"context"
	"log"
	"net/http"

	"github.com/bufbuild/connect-go"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	activity_v1 "badger-api/gen/activity/v1"
	"badger-api/gen/activity/v1/activity_v1connect"

	drill_v1 "badger-api/gen/drill/v1"
	"badger-api/gen/drill/v1/drill_v1connect"

	grpcreflect "github.com/bufbuild/connect-grpcreflect-go"
)

type ActivityServer struct{}
type DrillServer struct{}

func (s *ActivityServer) GetActivity(
	ctx context.Context,
	req *connect.Request[activity_v1.GetActivityRequest],
) (*connect.Response[activity_v1.GetActivityResponse], error) {

	log.Println("Request headers: ", req.Header())

	res := connect.NewResponse(&activity_v1.GetActivityResponse{})

	res.Header().Set("Example-Version", "v1")

	return res, nil
}

func (s *ActivityServer) GetActivities(
	ctx context.Context,
	req *connect.Request[activity_v1.GetActivitiesRequest],
) (*connect.Response[activity_v1.GetActivitiesResponse], error) {

	log.Println("Request headers: ", req.Header())

	res := connect.NewResponse(&activity_v1.GetActivitiesResponse{})

	res.Header().Set("Example-Version", "v1")

	return res, nil
}

func (s *DrillServer) GetDrill(
	ctx context.Context,
	req *connect.Request[drill_v1.GetDrillRequest],
) (*connect.Response[drill_v1.GetDrillResponse], error) {

	log.Println("Request headers: ", req.Header())

	res := connect.NewResponse(&drill_v1.GetDrillResponse{})

	res.Header().Set("Example-Version", "v1")

	return res, nil
}

func (s *DrillServer) GetDrills(
	ctx context.Context,
	req *connect.Request[drill_v1.GetDrillsRequest],
) (*connect.Response[drill_v1.GetDrillsResponse], error) {

	log.Println("Request headers: ", req.Header())

	res := connect.NewResponse(&drill_v1.GetDrillsResponse{})

	res.Header().Set("Example-Version", "v1")

	return res, nil
}

func main() {
	mux := http.NewServeMux()

	activityServer := &ActivityServer{}
	drillServer := &DrillServer{}

	reflector := grpcreflect.NewStaticReflector(
		activity_v1connect.ActivityServiceName,
		drill_v1connect.DrillServiceName,
	)

	path, handler := activity_v1connect.NewActivityServiceHandler(activityServer)
	mux.Handle(path, handler)

	path2, handler2 := drill_v1connect.NewDrillServiceHandler(drillServer)
	mux.Handle(path2, handler2)

	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	log.Println("Server running.")

	http.ListenAndServe(
		"0.0.0.0:3000",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
