package server

import (
	"context"
	"log"
	"os"
	"time"

	firebase "firebase.google.com/go"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ServerContext struct {
	db_client *mongo.Client
	db_ctx    context.Context
	db_cancel context.CancelFunc
	db        *mongo.Database

	firebase_app *firebase.App
}

func NewServerContext(db_conn_uri string) *ServerContext {
	log.Println("Creating new mongo client. ")
	db_client, err_client := mongo.NewClient(options.Client().ApplyURI(db_conn_uri).SetLoadBalanced(true))

	if err_client != nil {
		log.Fatal(err_client)
	}
	log.Println("Creating new mongo client done. ")

	db_ctx, db_ctx_cancel := context.WithTimeout(context.Background(), 10*time.Second)
	log.Println("Connecting to database client. ")
	err_conn := db_client.Connect(db_ctx)

	if err_conn != nil {
		log.Fatal(err_conn)
	}
	log.Println("Connecting to database client done. ")

	log.Println("Connecting to database. ")
	db := db_client.Database(os.Getenv("MONGO_DB_NAME"))
	log.Println("Connecting to database done. ")

	firebase_app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("Error initializing Firebase app: %v\n", err)
	}

	return &ServerContext{
		db_client,
		db_ctx,
		db_ctx_cancel,
		db,

		firebase_app,
	}
}

func (s *ServerContext) Cleanup() {
	s.db_client.Disconnect(s.db_ctx)
}

func (s *ServerContext) GetMongoDbClient() *mongo.Client {
	return s.db_client
}

func (s *ServerContext) GetCollection(collectionName string) *mongo.Collection {
	return s.db.Collection(collectionName)
}

func (s *ServerContext) GetMongoContext() context.Context {
	return context.Background()
}
