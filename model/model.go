package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Employee struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type Payload struct {
	Employees []Employee `json:"employees"`
}

type DbConfig struct {
	URL      string
	Database string
}

type User struct {
	ID                primitive.ObjectID `json:"id" bson:"_id"`
	Username          string             `json:"username" bson:"username"`
	Password          [32]byte           `json:"password" bson:"password"`
	FirstName         string             `json:"firstname" bson:"firstname"`
	LastName          string             `json:"lastname" bson:"lastname"`
	Email             string             `json:"email" bson:"email"`
	Team              string             `json:"team" bson:"team"`
	Projects          []string           `json:"projects" bson:"projects"`
	TotalWorkingHours float32            `json:"totalWorkingHours" bson:"totalWorkingHours"`
	VacationDays      int                `json:"vacationDays" bson:"vacationDays"`
	Permission        string             `json:"permission" bson:"permission"`
}

type UserSignupPayload struct {
	Username  string `json:"username" bson:"username" key:"required"`
	Password  string `json:"password" bson:"password" key:"required"`
	FirstName string `json:"firstname" bson:"firstname" key:"required"`
	LastName  string `json:"lastname" bson:"lastname" key:"required"`
	Email     string `json:"email" bson:"email" key:"required"`
}

type UserSignup struct {
	Username   string   `json:"username" bson:"username"`
	Password   [32]byte `json:"password" bson:"password"`
	FirstName  string   `json:"firstname" bson:"firstname"`
	LastName   string   `json:"lastname" bson:"lastname"`
	Email      string   `json:"email" bson:"email"`
	Permission string   `json:"permission" bson:"permission"`
}

type UserPayload struct {
	Username          string   `json:"username" bson:"username"`
	FirstName         string   `json:"firstname" bson:"firstname"`
	LastName          string   `json:"lastname" bson:"lastname"`
	Email             string   `json:"email" bson:"email"`
	Team              string   `json:"team" bson:"team"`
	Projects          []string `json:"projects" bson:"projects"`
	TotalWorkingHours float32  `json:"totalWorkingHours" bson:"totalWorkingHours"`
	VacationDays      int      `json:"vacationDays" bson:"vacationDays"`
	Permission        string   `json:"permission" bson:"permission"`
}

type UserUpdateResult struct {
	UpdateResult *mongo.UpdateResult `json:"result"`
	User         User                `json:"user"`
	Success      bool                `json:"success"`
}

type UserAuthPayload struct {
	Username string `json:"username"`
	Password string `json:"password" key:"required"`
}

type TeamMember struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Username  string             `json:"username" bson:"username"`
	FirstName string             `json:"first_name" bson:"first_name"`
	LastName  string             `json:"last_name" bson:"lastName"`
	Email     string             `json:"email" bson:"email"`
}

type TimeEntry struct {
	UserId     string    `json:"userId" bson:"userId"`
	Start      time.Time `json:"start" bson:"start"`
	End        time.Time `json:"end" bson:"end"`
	BreakStart time.Time `json:"breakStart" bson:"breakStart"`
	BreakEnd   time.Time `json:"breakEnd" bson:"breakEnd"`
	Project    string    `json:"project" bson:"project"`
}

type Proposal struct {
	UserId    string    `json:"userId" bson:"userId"`
	StartDate time.Time `json:"startDate" bson:"startDate"`
	EndDate   time.Time `json:"endDate" bson:"endDate"`
	Approved  bool      `json:"approved" bson:"approved"`
	Type      string    `json:"type" bson:"type"`
}
