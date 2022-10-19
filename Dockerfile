FROM golang:1.19-alpine AS base
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
COPY ./main.go .
RUN go build -o /usr/local/bin

FROM builder AS deploy
# RUN apk add --update --no-cache libc6-compat ca-certificates tzdata
CMD badger-api

FROM base AS devcontainer
RUN go install golang.org/x/tools/cmd/goimports@latest
RUN go install github.com/rogpeppe/godef

# FROM base AS builder
# WORKDIR /go
# COPY ./proto /go/proto
# COPY ./main.go /go/src/main.go
# COPY ./go.mod /go/go.mod
# COPY ./go.sum /go/go.sum
# RUN cd ./proto && buf generate && cd ..
# RUN go build