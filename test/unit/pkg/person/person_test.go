package person_test

import (
	person_v1 "badger-api/gen/person/v1"
	person_v1_connect "badger-api/gen/person/v1/personv1connect"
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

	client := person_v1_connect.NewPersonServiceClient(http.DefaultClient, "http://0.0.0.0:3000", connect.WithGRPC())

	const GOOD_PERSON_ID = "636acbb84eedb04302f88d72"
	const BAD_PERSON_ID = "42"

	g.Describe("PersonServer", func() {
		g.It("Should be able to retrieve all people", func() {
			req := connect.NewRequest(&person_v1.GetPeopleRequest{})
			res, err := client.GetPeople(context.Background(), req)
			Expect(res).NotTo(BeNil())
			Expect(len(res.Msg.GetPeople())).NotTo(Equal(0))
			Expect(err).To(BeNil())
		})

		g.It("Should be able to retrieve a person from a valid person id", func() {
			req := connect.NewRequest(&person_v1.GetPersonRequest{PersonId: GOOD_PERSON_ID})
			res, err := client.GetPerson(context.Background(), req)
			Expect(err).To(BeNil())
			Expect(res).NotTo(BeNil())
			Expect(res.Msg.GetPerson().GetUserId()).To(Equal(GOOD_PERSON_ID))
		})

		g.It("Should have an error when retrieving a person from an invalid person id", func() {
			req := connect.NewRequest(&person_v1.GetPersonRequest{PersonId: BAD_PERSON_ID})
			res, err := client.GetPerson(context.Background(), req)
			Expect(err).NotTo(BeNil())
			Expect(res).To(BeNil())
		})
	})

}
