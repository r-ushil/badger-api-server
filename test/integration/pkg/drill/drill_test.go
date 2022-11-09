package drill_test

import (
	drill_v1 "badger-api/gen/drill/v1"
	drill_v1_connect "badger-api/gen/drill/v1/drillv1connect"
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

	client := drill_v1_connect.NewDrillServiceClient(http.DefaultClient, "http://0.0.0.0:3000", connect.WithGRPC())

	const GOOD_DRILL_ID = "6352414e50c7d61db5d52863"
	const BAD_DRILL_ID = "42"

	g.Describe("DrillServer", func() {
		g.It("Should be able to retrieve all drills", func() {
			req := connect.NewRequest(&drill_v1.GetDrillsRequest{})
			res, err := client.GetDrills(context.Background(), req)
			Expect(res).NotTo(BeNil())
			Expect(len(res.Msg.GetDrills())).NotTo(Equal(0))
			Expect(err).To(BeNil())
		})

		g.It("Should be able to retrieve a drill from a valid drill id", func() {
			req := connect.NewRequest(&drill_v1.GetDrillRequest{DrillId: GOOD_DRILL_ID})
			res, err := client.GetDrill(context.Background(), req)
			Expect(err).To(BeNil())
			Expect(res).NotTo(BeNil())
			Expect(res.Msg.Drill.GetDrillId()).To(Equal(GOOD_DRILL_ID))
		})

		g.It("Should have an error when retrieving a drill from an invalid drill id", func() {
			req := connect.NewRequest(&drill_v1.GetDrillRequest{DrillId: BAD_DRILL_ID})
			res, err := client.GetDrill(context.Background(), req)
			Expect(err).ToNot(BeNil())
			Expect(res).To(BeNil())
		})
	})

}
