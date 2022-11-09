package drill_test

import (
	activity_v1 "badger-api/gen/activity/v1"
	activity_v1_connect "badger-api/gen/activity/v1/activityv1connect"
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

	client := activity_v1_connect.NewActivityServiceClient(http.DefaultClient, "http://0.0.0.0:3000", connect.WithGRPC())

	const GOOD_ACTIVITY_ID = "636b744a02f2bd8a671effc1"
	const BAD_ACTIVITY_ID = "42"

	g.Describe("ActivityServer", func() {
		g.It("Should be able to retrieve all activities", func() {
			req := connect.NewRequest(&activity_v1.GetActivitiesRequest{})
			res, err := client.GetActivities(context.Background(), req)
			Expect(res).NotTo(BeNil())
			Expect(len(res.Msg.GetActivities())).NotTo(Equal(0))
			Expect(err).To(BeNil())
		})

		g.It("Should be able to retrieve an activity from a valid activity id", func() {
			req := connect.NewRequest(&activity_v1.GetActivityRequest{ActivityId: GOOD_ACTIVITY_ID})
			res, err := client.GetActivity(context.Background(), req)
			Expect(err).To(BeNil())
			Expect(res).NotTo(BeNil())
			Expect(res.Msg.GetActivity().GetActivityId()).To(Equal(GOOD_ACTIVITY_ID))
		})

		g.It("Should have an error when retrieving a drill from an invalid drill id", func() {
			req := connect.NewRequest(&activity_v1.GetActivityRequest{ActivityId: BAD_ACTIVITY_ID})
			res, err := client.GetActivity(context.Background(), req)
			Expect(err).ToNot(BeNil())
			Expect(res).To(BeNil())
		})
	})
}
