package drill

import (
	"badger-api/pkg/server"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Drill struct {
	Id          string `bson:"_id"`
	Name        string `bson:"name"`
	Description string `bson:"description"`
}

func (d *Drill) GetId() string {
	return d.Id
}

func (d *Drill) GetName() string {
	return d.Name
}

func (d *Drill) GetDescription() string {
	return d.Description
}

func GetDrills(s *server.ServerContext) []Drill {
	client := s.GetMongoDbClient()

	db := client.Database("badger_db")
	col := db.Collection("drills")

	cur, err_find := col.Find(s.GetMongoContext(), bson.D{})

	if err_find != nil {
		panic(err_find)
	}

	var drills []Drill
	err_cur := cur.All(context.TODO(), &drills)

	if err_cur != nil {
		panic(err_cur)
	}

	return drills
}

var ErrNotFound error

func GetDrill(s *server.ServerContext, id string) (*Drill, error) {
	client := s.GetMongoDbClient()

	db := client.Database("badger_db")
	col := db.Collection("drills")

	query := bson.D{{Key: "_id", Value: id}}

	var drill Drill
	err := col.FindOne(s.GetMongoContext(), query).Decode(&drill)

	if err == mongo.ErrNoDocuments {
		return nil, ErrNotFound
	}

	if err != nil {
		panic(err)
	}

	return &drill, nil
}
