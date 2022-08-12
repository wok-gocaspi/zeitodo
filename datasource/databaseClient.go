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
	InsertOne(ctx context.Context, documents interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (cur *mongo.Cursor, err error)
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
func (c Client) CreatTimeEntryById(te model.TimeEntry) (interface{}, error) {
	result, err := c.TimeEntries.InsertOne(context.TODO(), te)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (c Client) GetTimeEntryByUserID(id string) []model.TimeEntry {

	filter := bson.M{"userId": id}
	var timeEntries []model.TimeEntry
	courser, err := c.TimeEntries.Find(context.TODO(), filter)
	if err != nil {
		return nil
	}
	for courser.Next(context.TODO()) {
		var timeEntry model.TimeEntry
		err := courser.Decode(&timeEntry)
		if err != nil {
			return timeEntries
		}
		timeEntries = append(timeEntries, timeEntry)
	}
	return timeEntries
}

func (c Client) GetProposals(id string) ([]model.Proposal, error) {
	var proposalArr []model.Proposal
	filter := bson.M{"userId": id}
	cur, err := c.Proposals.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		// To decode into a struct, use cursor.Decode()
		var proposal model.Proposal
		err = cur.Decode(&proposal)
		if err != nil {
			return nil, err
		}
		proposalArr = append(proposalArr, proposal)
	}
	if err = cur.Err(); err != nil {
		return nil, err
	}
	return proposalArr, nil
}

func (c Client) SaveProposals(docs []interface{}) (interface{}, error) {
	results, err := c.Proposals.InsertMany(context.TODO(), docs)
	if err != nil {
		log.Println("database error")
		return nil, err
	}
	return results.InsertedIDs, nil
}

func (c Client) DeleteProposalByIdAndDate(id string, date string) (*mongo.DeleteResult, error) {

	filter := bson.M{"startDate": date, "userId": id}

	results, err := c.Proposals.DeleteOne(context.TODO(), filter)

	if err != nil {

		return nil, err
	}
	if results.DeletedCount == 0 {
		deleterror := errors.New("the Employee id is not existing")
		return nil, deleterror
	}
	return results, nil
}

func (c Client) UpdateProposal(update model.Proposal, date string) (*mongo.UpdateResult, error) {
	filter := bson.M{"userId": update.UserId, "startDate": date}
	// datensatz zur id auslesen
	// check doc geschnitten datensatzen
	// change update
	if update.UserId == "" {
		IdMissing := "an userId has to be provided"
		return nil, errors.New(IdMissing)
	}
	courser := c.Proposals.FindOne(context.TODO(), filter)
	var proposal model.Proposal
	err := courser.Decode(&proposal)
	if proposal.UserId == "" {
		IdWrong := "an userId has to be provided"
		return nil, errors.New(IdWrong)
	}
	fmt.Println(update)
	var setElements bson.D
	if update.StartDate != "" {
		setElements = append(setElements, bson.E{Key: "startDate", Value: update.StartDate})
	}
	if update.EndDate != "" {
		setElements = append(setElements, bson.E{Key: "endDate", Value: update.EndDate})
	}
	if update.Type != "" {
		setElements = append(setElements, bson.E{Key: "type", Value: update.Type})
	}

	setMap := bson.D{
		{"$set", setElements},
	}
	result, err := c.Proposals.UpdateOne(context.TODO(), filter, setMap)
	if err != nil {
		return nil, err

	}
	return result, nil
}
