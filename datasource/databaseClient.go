package datasource

import (
	"context"
	"errors"
	"example-project/model"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . MongoDBInterface
type MongoDBInterface interface {
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	InsertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
}

type Client struct {
	Users       MongoDBInterface
	TimeEntries MongoDBInterface
	Proposals   MongoDBInterface
}

func NewDbClient(d model.DbConfig) Client {
	client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI(d.URL))
	db := client.Database(d.Database)
	return Client{
		Users:       db.Collection("Users"),
		TimeEntries: db.Collection("TimeEntries"),
		Proposals:   db.Collection("Proposals"),
	}
}

func (c Client) UpdateMany(docs []interface{}) interface{} {
	results, err := c.Users.InsertMany(context.TODO(), docs)
	if err != nil {
		log.Println("database error")
	}
	return results.InsertedIDs
}

func (c Client) GetByID(id string) model.Employee {
	filter := bson.M{"id": id}
	courser := c.Users.FindOne(context.TODO(), filter)
	var employee model.Employee
	err := courser.Decode(&employee)
	if err != nil {
		log.Println("error during data marshalling")
	}
	return employee
}

func (c Client) DeleteTimeEntryById(id string) (interface{}, error) {
	filter := bson.M{"id": id}

	results, err := c.TimeEntries.DeleteOne(context.TODO(), filter)

	if err != nil {

		return nil, err
	}
	if results.DeletedCount == 0 {
		deleterror := errors.New(" Time not existing")
		return nil, deleterror
	}
	return results.DeletedCount, nil
}
func (c Client) UpdateTimeEntryById(update model.TimeEntry) (*mongo.UpdateResult, error) {
	filter := bson.M{"id": update.Start}
	if update.UserId == "" {
		IdMissing := fmt.Sprintf("ID %v got no Time", update.UserId)
		return nil, errors.New(IdMissing)
	}
	courser := c.TimeEntries.FindOne(context.TODO(), filter)
	var employee model.Employee
	err := courser.Decode(&employee)
	if employee.ID == "" {
		IdWrong := fmt.Sprintf("User %v dosent exist", update.UserId)
		return nil, errors.New(IdWrong)
	}
	fmt.Println(update)
	var setElements bson.D
	if update.Project != "" {
		fmt.Sprintf(update.Project)
		setElements = append(setElements, bson.E{Key: "firstname", Value: update.Project})
	}
	if update.UserId != "" {
		fmt.Sprintf(update.UserId)
		setElements = append(setElements, bson.E{Key: "lastname", Value: update.UserId})
	}
	if update.UserId != "" {
		fmt.Sprintf(update.UserId)
		setElements = append(setElements, bson.E{Key: "email", Value: update.UserId})
	}
	setMap := bson.D{
		{"$set", setElements},
	}
	result, err := c.TimeEntries.UpdateOne(context.TODO(), filter, setMap)
	if err != nil {
		return nil, err

	}
	return result, nil
}
func (c Client) CreatTimeEntryById(id string) (interface{}, error) {
	filter := bson.M{"id": id}
	courser := c.Users.FindOne(context.TODO(), filter)
	var timee model.TimeEntry
	err := courser.Decode(&timee)
	if err != nil {
		log.Println("Time Entry")
	}
	return timee, nil
}
func (c Client) GetTimeEntryById(id string) model.TimeEntry {
	filter := bson.M{"id": id}
	courser := c.Users.FindOne(context.TODO(), filter)
	var timee model.TimeEntry
	err := courser.Decode(&timee)
	if err != nil {
		log.Println("error ")
	}
	return timee
}

func (c Client) GetAllTimeEntriesById(id string) model.TimeEntry {
	filter := bson.M{"id": id}
	courser := c.Users.FindOne(context.TODO(), filter)
	var Alltimee model.TimeEntry
	err := courser.Decode(&Alltimee)
	if err != nil {
		log.Println("error ")
	}
	return Alltimee
}
