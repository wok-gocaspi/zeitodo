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
	"log"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . MongoDBInterface
type MongoDBInterface interface {
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	InsertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
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

func (c Client) GetUserByUsername(username string) (model.UserPayload, error) {
	filter := bson.M{"username": username}
	courser := c.Users.FindOne(context.TODO(), filter)
	var user model.UserPayload
	err := courser.Decode(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (c Client) GetUserByEmail(email string) (model.UserPayload, error) {
	filter := bson.M{"email": email}
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

func (c Client) CreateUser(docs interface{}) (interface{}, error) {
	results, err := c.Users.InsertOne(context.TODO(), docs)
	if err != nil {
		return nil, err
	}
	return results.InsertedID, nil
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
func (c Client) CreatTimeEntryById(te model.TimeEntry) (interface{}, error) {
	result, err := c.TimeEntries.InsertOne(context.TODO(), te)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (c Client) UpdateTimeEntryById(update model.TimeEntry) (*mongo.UpdateResult, error) {
	filter := bson.M{"userId": update.UserId}

	if update.UserId == "" {
		IdMissing := fmt.Sprintf("ID %v got no User", update.UserId)
		return nil, errors.New(IdMissing)
	}
	courser := c.TimeEntries.FindOne(context.TODO(), filter)
	var User model.TimeEntry
	err := courser.Decode(&User)

	fmt.Println(update)
	var setElements bson.D

	if update.Project != "" {
		fmt.Sprintf(update.Project)
		setElements = append(setElements, bson.E{Key: "project", Value: update.Project})
	}

	if update.UserId != "" {
		fmt.Sprintf(update.UserId)
		setElements = append(setElements, bson.E{Key: "userId", Value: update.UserId})
	}

	if !update.Start.IsZero() {
		fmt.Sprintf(update.Start.String())
		setElements = append(setElements, bson.E{Key: "start", Value: update.Start})
	}

	if !update.BreakStart.IsZero() {
		fmt.Sprintf(update.BreakStart.String())
		setElements = append(setElements, bson.E{Key: "breakStart", Value: update.BreakStart})
	}

	if !update.BreakEnd.IsZero() {
		fmt.Sprintf(update.BreakEnd.String())
		setElements = append(setElements, bson.E{Key: "breakEnd", Value: update.BreakEnd})
	}
	if !update.End.IsZero() {
		fmt.Sprintf(update.End.String())
		setElements = append(setElements, bson.E{Key: "end", Value: update.End})
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
func (c Client) GetTimeEntryByID(id string) []model.TimeEntry {

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
