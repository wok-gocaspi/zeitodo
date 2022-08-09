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
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (cur *mongo.Cursor, err error)
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
		IdMissing := fmt.Sprintf("User %v got no ID", update.UserId)
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
