package server

import (
	"net/http"

	activity_v1connect "badger-api/gen/activity/v1/activityv1connect"
	drill_v1connect "badger-api/gen/drill/v1/drillv1connect"
	person_v1connect "badger-api/gen/person/v1/personv1connect"

	grpcreflect "github.com/bufbuild/connect-grpcreflect-go"
)

func RegisterReflector(mux *http.ServeMux) {
	reflector := grpcreflect.NewStaticReflector(
		activity_v1connect.ActivityServiceName,
		drill_v1connect.DrillServiceName,
		person_v1connect.PersonServiceName,
	)

	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))
}
