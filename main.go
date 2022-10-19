package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/bufbuild/connect-go"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	examplev1 "badger-api/gen/example/v1"
	"badger-api/gen/example/v1/examplev1connect"

	grpcreflect "github.com/bufbuild/connect-grpcreflect-go"
)

type ExampleServer struct{}

func (s *ExampleServer) Example(
	ctx context.Context,
	req *connect.Request[examplev1.ExampleRequest],
) (*connect.Response[examplev1.ExampleResponse], error) {

	log.Println("Request headers: ", req.Header())

	res := connect.NewResponse(&examplev1.ExampleResponse{
		ExampleText: fmt.Sprintf("Hello, %s!", req.Msg.Name),
	})

	res.Header().Set("Example-Version", "v1")

	return res, nil
}

func main() {
	exampleServer := &ExampleServer{}
	mux := http.NewServeMux()

	reflector := grpcreflect.NewStaticReflector(
		examplev1connect.ExampleServiceName,
	)

	path, handler := examplev1connect.NewExampleServiceHandler(exampleServer)
	mux.Handle(path, handler)

	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	log.Println("Server running.")

	http.ListenAndServe(
		"0.0.0.0:3000",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
