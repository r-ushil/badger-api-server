package server

import (
	"context"
	"log"
	"net/http"

	"github.com/bufbuild/connect-go"

	person_v1 "badger-api/gen/person/v1"
	"badger-api/gen/person/v1/personv1connect"
	"badger-api/pkg/person"
	"badger-api/pkg/server"
)

type PersonServer struct {
	ctx *server.ServerContext
}

func (s *PersonServer) GetPerson(
	ctx context.Context,
	req *connect.Request[person_v1.GetPersonRequest],
) (*connect.Response[person_v1.GetPersonResponse], error) {
	log.Println("Getting person with ID:", req.Msg.PersonId)

	d, err := person.GetPerson(s.ctx, req.Msg.PersonId)

	if err != nil {
		// TODO: Handle properly
		log.Println(err)
		return nil, connect.NewError(connect.CodeUnimplemented, err)
	}

	res := connect.NewResponse(&person_v1.GetPersonResponse{
		Person: &person_v1.Person{
			UserId:     d.GetId(),
			UserScore:  d.GetScore(),
			FirebaseId: d.GetFirebaseId(),
		},
	})

	return res, nil
}

func (s *PersonServer) GetPeople(
	ctx context.Context,
	req *connect.Request[person_v1.GetPeopleRequest],
) (*connect.Response[person_v1.GetPeopleResponse], error) {
	ds := person.GetPeople(s.ctx)

	people := make([]*person_v1.Person, 0, len(ds))

	for _, person := range ds {
		people = append(people, &person_v1.Person{
			UserId:     person.GetId(),
			UserScore:  person.GetScore(),
			FirebaseId: person.GetFirebaseId(),
		})
	}

	res := connect.NewResponse(&person_v1.GetPeopleResponse{
		People: people,
	})

	return res, nil
}

func RegisterPersonService(mux *http.ServeMux, ctx *server.ServerContext) {
	server := &PersonServer{
		ctx,
	}

	path, handler := personv1connect.NewPersonServiceHandler(server)

	mux.Handle(path, handler)
}
