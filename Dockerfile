FROM golang:1.19-bullseye AS base
WORKDIR /workspaces/badger-api-server
RUN go install github.com/bufbuild/buf/cmd/buf@latest
RUN go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install github.com/bufbuild/connect-go/cmd/protoc-gen-connect-go@latest

FROM base AS deps
COPY ./go.mod .
COPY ./go.sum .
RUN go mod download

FROM deps AS proto_builder
COPY ./buf.gen.yaml .
COPY ./buf.work.yaml .
COPY ./proto ./proto
RUN buf generate

FROM proto_builder AS builder
COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./pkg ./pkg
RUN go build -o /usr/local/bin/badger_api ./cmd

FROM builder AS deploy
CMD badger_api

FROM base AS devcontainer
RUN go install golang.org/x/tools/cmd/goimports@latest
RUN go install golang.org/x/tools/cmd/callgraph@latest
RUN go install golang.org/x/tools/cmd/digraph@latest
RUN go install golang.org/x/tools/cmd/stringer@latest
RUN go install golang.org/x/tools/cmd/toolstash@latest
RUN go install golang.org/x/tools/gopls@latest