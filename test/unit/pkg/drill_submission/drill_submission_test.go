package drill_test

import (
	drill_submission_v1 "badger-api/gen/drill_submission/v1"
	drill_submission_v1_connect "badger-api/gen/drill_submission/v1/drill_submissionv1connect"
	"context"
	"net/http"
	"testing"

	"github.com/bufbuild/connect-go"
	goblin "github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func TestActivity(t *testing.T) {
	g := goblin.Goblin(t)

	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	client := drill_submission_v1_connect.NewDrillSubmissionServiceClient(http.DefaultClient, "http://0.0.0.0:3000", connect.WithGRPC())

	const GOOD_DRILL_SUBMISSION_ID = "637177971a12146c88e3158f"
	const BAD_DRILL_SUBMISSION_ID = "42"

	g.Describe("DrillServer", func() {
		g.It("Should be able to retrieve all drill submissions", func() {
			req := connect.NewRequest(&drill_submission_v1.GetDrillSubmissionsRequest{})
			res, err := client.GetDrillSubmissions(context.Background(), req)
			Expect(err).To(BeNil())
			Expect(res).NotTo(BeNil())
			Expect(len(res.Msg.GetDrillSubmissions())).NotTo(Equal(0))
		})

		g.It("Should be able to retrieve a drill submission from a valid drill submission id", func() {
			req := connect.NewRequest(&drill_submission_v1.GetDrillSubmissionRequest{DrillSubmissionId: GOOD_DRILL_SUBMISSION_ID})
			res, err := client.GetDrillSubmission(context.Background(), req)
			Expect(err).To(BeNil())
			Expect(res).NotTo(BeNil())
			Expect(res.Msg.DrillSubmission.GetDrillSubmissionId()).To(Equal(GOOD_DRILL_SUBMISSION_ID))
		})

		g.It("Should have an error when retrieving a drill submission from an invalid drill submission id", func() {
			req := connect.NewRequest(&drill_submission_v1.GetDrillSubmissionRequest{DrillSubmissionId: BAD_DRILL_SUBMISSION_ID})
			res, err := client.GetDrillSubmission(context.Background(), req)
			Expect(err).NotTo(BeNil())
			Expect(res).To(BeNil())
		})
	})

}
