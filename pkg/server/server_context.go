package server

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ServerContext struct {
	db_client *mongo.Client
	db_ctx    context.Context
	db_cancel context.CancelFunc
}

func NewServerContext(db_conn_uri string) *ServerContext {
	db_client, err_client := mongo.NewClient(options.Client().ApplyURI(db_conn_uri))

	if err_client != nil {
		log.Fatal(err_client)
	}

	db_ctx, db_ctx_cancel := context.WithTimeout(context.Background(), 10*time.Second)

	err_conn := db_client.Connect(db_ctx)

	if err_conn != nil {
		log.Fatal(err_conn)
	}

	return &ServerContext{
		db_client,
		db_ctx,
		db_ctx_cancel,
	}
}

func (s *ServerContext) Cleanup() {
	s.db_client.Disconnect(s.db_ctx)
}

func (s *ServerContext) GetMongoDbClient() *mongo.Client {
	return s.db_client
}

func (s *ServerContext) GetMongoContext() context.Context {
	return context.Background()
}
