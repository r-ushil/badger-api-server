package drill

import (
	"badger-api/pkg/server"
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Drill struct {
	Id           string   `bson:"_id"`
	Name         string   `bson:"name"`
	Description  string   `bson:"description"`
	Instructions string   `bson:"instructions"`
	ThumbnailUrl string   `bson:"thumbnail_url"`
	Skills       []string `bson:"skills"`
	VideoUrl     string   `bson:"video_url"`
	Duration     int32    `bson:"duration"`
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

func (d *Drill) GetInstructions() string {
	return d.Instructions
}

func (d *Drill) GetThumbnailUrl() string {
	return d.ThumbnailUrl
}

func (d *Drill) GetSkills() []string {
	return d.Skills
}

func (d *Drill) GetVideoUrl() string {
	return d.VideoUrl
}

func (d *Drill) GetDuration() int32 {
	return d.Duration
}

func GetDrills(s *server.ServerContext) []Drill {
	col := s.GetCollection("drills")

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

func GetDrill(s *server.ServerContext, hexId string) (*Drill, error) {
	log.Println("Getting drill collection. ")
	col := s.GetCollection("drills")
	log.Println("Getting drill collection done. ")

	objectId, idErr := primitive.ObjectIDFromHex(hexId)

	if idErr != nil {
		panic(idErr)
	}

	query := bson.D{{Key: "_id", Value: objectId}}

	var drill Drill
	log.Println("Getting drill document. ")
	err := col.FindOne(s.GetMongoContext(), query).Decode(&drill)
	log.Println("Getting drill document done. ")

	if err == mongo.ErrNoDocuments {
		return nil, errors.New("ErrNotFound")
	}

	if err != nil {
		panic(err)
	}

	log.Println("All good, returning drill. ")
	log.Println(drill)

	return &drill, nil
}
