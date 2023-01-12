package server

import (
	"net/http"

	batting_drillv1connect "badger-api/gen/batting_drill/v1/batting_drillv1connect"
	"badger-api/gen/catching_drill/v1/catching_drillv1connect"
	drill_v1connect "badger-api/gen/drill/v1/drillv1connect"
	drill_submission_v1connect "badger-api/gen/drill_submission/v1/drill_submissionv1connect"
	"badger-api/gen/leaderboard/v1/leaderboardv1connect"
	person_v1connect "badger-api/gen/person/v1/personv1connect"

	grpcreflect "github.com/bufbuild/connect-grpcreflect-go"
)

func RegisterReflector(mux *http.ServeMux) {
	reflector := grpcreflect.NewStaticReflector(
		drill_v1connect.DrillServiceName,
		person_v1connect.PersonServiceName,
		drill_submission_v1connect.DrillSubmissionServiceName,
		batting_drillv1connect.BattingDrillServiceName,
		catching_drillv1connect.CatchingDrillServiceName,
		leaderboardv1connect.LeaderboardServiceName,
	)

	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))
}
