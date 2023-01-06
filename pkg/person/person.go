package person

import (
	"badger-api/pkg/server"
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Person struct {
	Id           string `bson:"_id"`
	Score        uint32 `bson:"score"`
	FirebaseId   string `bson:"firebase_id"`
	PowerScore   uint32 `bson:"power_score"`
	TimingScore  uint32 `bson:"timing_score"`
	AgilityScore uint32 `bson:"agility_score"`
}

func (d *Person) GetId() string {
	return d.Id
}

func (d *Person) GetScore() uint32 {
	return d.Score
}

func (d *Person) GetFirebaseId() string {
	return d.FirebaseId
}

func (d *Person) GetPowerScore() uint32 {
	return d.PowerScore
}

func (d *Person) GetTimingScore() uint32 {
	return d.TimingScore
}

func (d *Person) GetAgilityScore() uint32 {
	return d.AgilityScore
}

func GetPeople(s *server.ServerContext) []Person {
	col := s.GetCollection("people")

	cur, err_find := col.Find(s.GetMongoContext(), bson.D{})

	if err_find != nil {
		panic(err_find)
	}

	var people []Person
	err_cur := cur.All(context.TODO(), &people)

	if err_cur != nil {
		panic(err_cur)
	}

	return people
}

func GetPerson(s *server.ServerContext, hexId string) (*Person, error) {
	log.Println("Getting person collection. ")
	col := s.GetCollection("people")
	log.Println("Getting person collection done. ")

	query := bson.D{{Key: "firebaseid", Value: hexId}}

	var person Person
	log.Println("Getting person document. ")
	err := col.FindOne(s.GetMongoContext(), query).Decode(&person)
	log.Println("Getting person document done. ")

	if err == mongo.ErrNoDocuments {
		return nil, errors.New("ErrNotFound")
	}

	if err != nil {
		panic(err)
	}

	log.Println("All good, returning person. ")
	log.Println(person)

	return &person, nil
}

func InsertNewUser(s *server.ServerContext, firebase_id string) string {
	col := s.GetCollection("people")

	data := Person{
		FirebaseId: firebase_id,
		Score: 0,
		PowerScore: 0,
		AgilityScore: 0,
		TimingScore: 0,
	}

	result, err := col.InsertOne(s.GetMongoContext(), data)

	if err != nil {
		panic(err)
	}
	print(result.InsertedID.(primitive.ObjectID).Hex())
	return result.InsertedID.(primitive.ObjectID).Hex()
}
