package datasource

import (
	"context"
	"errors"
	"example-project/model"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . MongoDBInterface
type MongoDBInterface interface {
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	InsertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
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

func (c Client) GetUserByID(id primitive.ObjectID) (model.UserPayload, error) {
	filter := bson.M{"_id": id}
	courser := c.Users.FindOne(context.TODO(), filter)
	var user model.UserPayload
	err := courser.Decode(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (c Client) GetAllUser() ([]model.UserPayload, error) {
	filter := bson.M{}
	coursor, err := c.Users.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var users []model.UserPayload
	for coursor.Next(context.TODO()) {
		var user model.UserPayload
		err := coursor.Decode(&user)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	if len(users) == 0 {
		noUserError := errors.New("no Users exist")
		return users, noUserError
	}
	return users, nil
}

func (c Client) GetUserTeamMembersByID(userid primitive.ObjectID) (interface{}, error) {
	userFilter := bson.M{"_id": userid}
	UserCoursor := c.Users.FindOne(context.TODO(), userFilter)
	var matchingUser model.User
	err := UserCoursor.Decode(&matchingUser)
	if err != nil {
		return nil, err
	}
	teamFilter := bson.M{"team": matchingUser.Team}
	teamCoursor, err := c.Users.Find(context.TODO(), teamFilter)
	if err != nil {
		return nil, err
	}
	var teamMembers []model.TeamMember
	for teamCoursor.Next(context.TODO()) {
		var teamMember model.TeamMember
		err := teamCoursor.Decode(&teamMember)
		if err != nil {
			return nil, err
		}
		teamMembers = append(teamMembers, teamMember)

	}
	return teamMembers, nil

}

func (c Client) GetUserTeamMembersByName(team string) (interface{}, error) {
	teamFilter := bson.M{"team": team}
	teamCoursor, err := c.Users.Find(context.TODO(), teamFilter)
	if err != nil {
		return nil, err
	}
	var teamMembers []model.TeamMember
	for teamCoursor.Next(context.TODO()) {
		var teamMember model.TeamMember
		err := teamCoursor.Decode(&teamMember)
		if err != nil {
			return nil, err
		}
		teamMembers = append(teamMembers, teamMember)

	}
	return teamMembers, nil

}

func (c Client) CreateUser(docs []interface{}) (interface{}, error) {
	results, err := c.Users.InsertMany(context.TODO(), docs)
	if err != nil {
		return nil, err
	}
	return results.InsertedIDs, nil
}

func (c Client) UpdateManyUserByID(docs []model.User) []model.UserUpdateResult {
	var UMR []model.UserUpdateResult
	for _, user := range docs {
		var UpdateResult model.UserUpdateResult
		filter := bson.M{"_id": user.ID}
		if user.ID.String() == "" {
			UpdateResult.User = user
			UpdateResult.Success = false
			UMR = append(UMR, UpdateResult)
			continue
		}
		courser := c.Users.FindOne(context.TODO(), filter)
		var userDoc model.User
		err := courser.Decode(&userDoc)
		if err != nil || userDoc.ID.String() == "" {
			UpdateResult.Success = false
			UpdateResult.User = user
			UMR = append(UMR, UpdateResult)
			continue
		}

		var setElements bson.D
		if user.FirstName != "" {
			fmt.Println(user.FirstName)
			setElements = append(setElements, bson.E{Key: "firstname", Value: user.FirstName})
		}
		if user.LastName != "" {
			setElements = append(setElements, bson.E{Key: "lastname", Value: user.LastName})
		}
		if user.Email != "" {
			setElements = append(setElements, bson.E{Key: "email", Value: user.Email})
		}
		if user.Team != "" {
			setElements = append(setElements, bson.E{Key: "team", Value: user.Team})
		}
		if user.TotalWorkingHours != 0 {
			setElements = append(setElements, bson.E{Key: "totalWorkingHours", Value: user.TotalWorkingHours})
		}
		if user.VacationDays != 0 {
			setElements = append(setElements, bson.E{Key: "vacationDays", Value: user.VacationDays})
		}
		if user.Username != "" {
			setElements = append(setElements, bson.E{Key: "username", Value: user.Username})
		}
		fmt.Println(len(user.Projects))
		if len(user.Projects) != 0 {
			setElements = append(setElements, bson.E{Key: "projects", Value: user.Projects})
		}
		setMap := bson.D{
			{"$set", setElements},
		}
		result, err := c.Users.UpdateOne(context.TODO(), filter, setMap)
		if err != nil {
			UpdateResult.Success = false
			UpdateResult.User = user
			UMR = append(UMR, UpdateResult)
			continue
		}
		UpdateResult.Success = true
		UpdateResult.UpdateResult = result
		UpdateResult.User = user
		UMR = append(UMR, UpdateResult)
	}
	return UMR
}

func (c Client) DeleteUser(id primitive.ObjectID) (interface{}, error) {
	filter := bson.M{"_id": id}
	deleteResult, err := c.Users.DeleteOne(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	if deleteResult.DeletedCount == 0 {
		deleteError := errors.New("no user have been deleted, please check the id")
		return nil, deleteError
	}
	return deleteResult, nil
}
