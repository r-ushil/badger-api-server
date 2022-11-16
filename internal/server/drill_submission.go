package server

import (
	"context"
	"log"
	"net/http"

	"github.com/bufbuild/connect-go"

	drill_submission_v1 "badger-api/gen/drill_submission/v1"
	"badger-api/gen/drill_submission/v1/drill_submissionv1connect"

	"badger-api/pkg/drill_submission"
	"badger-api/pkg/server"
)

type DrillSubmissionServer struct {
	ctx *server.ServerContext
}

func (s *DrillSubmissionServer) GetDrillSubmission(
	ctx context.Context,
	req *connect.Request[drill_submission_v1.GetDrillSubmissionRequest],
) (*connect.Response[drill_submission_v1.GetDrillSubmissionResponse], error) {
	log.Println("Getting drill submission with ID:", req.Msg.DrillSubmissionId)

	d, err := drill_submission.GetDrillSubmission(s.ctx, req.Msg.DrillSubmissionId)

	if err != nil {
		// TODO: Handle properly
		log.Println(err)
		return nil, connect.NewError(connect.CodeUnimplemented, err)
	}
	timestampGoogleFormat := d.GetTimestampGoogleFormat()
	res := connect.NewResponse(&drill_submission_v1.GetDrillSubmissionResponse{
		DrillSubmission: &drill_submission_v1.DrillSubmission{
			DrillSubmissionId: d.GetId(),
			UserId:            d.GetUserId(),
			DrillId:           d.GetDrillId(),
			BucketUrl:         d.GetBucketUrl(),
			Timestamp:         &timestampGoogleFormat,
			ProcessingStatus:  d.GetProcessingStatus(),
			DrillScore:        d.GetDrillScore(),
		},
	})

	return res, nil
}

func (s *DrillSubmissionServer) GetDrillSubmissions(
	ctx context.Context,
	req *connect.Request[drill_submission_v1.GetDrillSubmissionsRequest],
) (*connect.Response[drill_submission_v1.GetDrillSubmissionsResponse], error) {
	dss := drill_submission.GetDrillSubmissions(s.ctx)

	drill_submissions := make([]*drill_submission_v1.DrillSubmission, 0, len(dss))

	for _, d := range dss {
		timestampGoogleFormat := d.GetTimestampGoogleFormat()
		drill_submissions = append(drill_submissions, &drill_submission_v1.DrillSubmission{
			DrillSubmissionId: d.GetId(),
			UserId:            d.GetUserId(),
			DrillId:           d.GetDrillId(),
			BucketUrl:         d.GetBucketUrl(),
			Timestamp:         &timestampGoogleFormat,
			ProcessingStatus:  d.GetProcessingStatus(),
			DrillScore:        d.GetDrillScore(),
		})
	}

	res := connect.NewResponse(&drill_submission_v1.GetDrillSubmissionsResponse{
		DrillSubmissions: drill_submissions,
	})

	return res, nil
}
func (s *DrillSubmissionServer) InsertDrillSubmission(
	ctx context.Context,
	req *connect.Request[drill_submission_v1.InsertDrillSubmissionRequest],
) (*connect.Response[drill_submission_v1.InsertDrillSubmissionResponse], error) {
	hex_id := drill_submission.InsertDrillSubmission(s.ctx, req.Msg.DrillSubmission)
	res := connect.NewResponse(&drill_submission_v1.InsertDrillSubmissionResponse{
		HexId: hex_id,
	})
	return res, nil
}

func (s *DrillSubmissionServer) GetUserDrillSubmissions(
	ctx context.Context,
	req *connect.Request[drill_submission_v1.GetUserDrillSubmissionsRequest],
) (*connect.Response[drill_submission_v1.GetUserDrillSubmissionsResponse], error) {
	dss := drill_submission.GetUserDrillSubmissions(s.ctx, req.Msg.UserId)

	drill_submissions := make([]*drill_submission_v1.DrillSubmission, 0, len(dss))

	for _, d := range dss {
		timestampGoogleFormat := d.GetTimestampGoogleFormat()
		drill_submissions = append(drill_submissions, &drill_submission_v1.DrillSubmission{
			DrillSubmissionId: d.GetId(),
			UserId:            d.GetUserId(),
			DrillId:           d.GetDrillId(),
			BucketUrl:         d.GetBucketUrl(),
			Timestamp:         &timestampGoogleFormat,
			ProcessingStatus:  d.GetProcessingStatus(),
			DrillScore:        d.GetDrillScore(),
		})
	}

	res := connect.NewResponse(&drill_submission_v1.GetUserDrillSubmissionsResponse{
		DrillSubmissions: drill_submissions,
	})

	return res, nil
}

func RegisterDrillSubmissionService(mux *http.ServeMux, ctx *server.ServerContext) {
	server := &DrillSubmissionServer{
		ctx,
	}

	path, handler := drill_submissionv1connect.NewDrillSubmissionServiceHandler(server)

	mux.Handle(path, handler)
}
