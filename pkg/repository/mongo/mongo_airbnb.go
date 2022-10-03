package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CRUDInstance struct {
	db *mongo.Collection
}

func NewAirbnbMongoInstance(db *mongo.Client) *CRUDInstance {
	return &CRUDInstance{db: db.Database("testDB").Collection("airbnbCollection")}
}

//get all instance by id
//filter by minimum nights
//filter by price
func (i *CRUDInstance) GetOfferById(id string) (bson.M, error) {
	var result bson.M
	//projection := bson.D{{"price", 1}}

	opts := options.Find().SetProjection(bson.D{{"_id", id}, {"price", 1}})
	cursor, err := i.db.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		panic(err)
	}
	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	for _, result := range results {
		fmt.Println(result)
	}
	return result, nil
}

func (i *CRUDInstance) FindByType(typename string) error {

	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"rating", bson.D{{"$gt", 7}}}},
				bson.D{{"rating", bson.D{{"$lte", 10}}}},
			}},
	}
	cursor, err := i.db.Find(context.TODO(), filter)
	if err != nil {
		return err
	}

	var results []bson.D

	err = cursor.All(context.TODO(), &results)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(results)

	return nil
}
