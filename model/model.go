package model

import (
	"github.com/retailify/go-interval"
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
}

type UserSignupPayload struct {
	Username  string `json:"username" bson:"username" key:"required"`
	Password  string `json:"password" bson:"password" key:"required"`
	FirstName string `json:"firstname" bson:"firstname" key:"required"`
	LastName  string `json:"lastname" bson:"lastname" key:"required"`
	Email     string `json:"email" bson:"email" key:"required"`
}

type UserSignup struct {
	Username  string   `json:"username" bson:"username"`
	Password  [32]byte `json:"password" bson:"password"`
	FirstName string   `json:"firstname" bson:"firstname"`
	LastName  string   `json:"lastname" bson:"lastname"`
	Email     string   `json:"email" bson:"email"`
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
}

type UserUpdateResult struct {
	UpdateResult *mongo.UpdateResult `json:"result"`
	User         User                `json:"user"`
	Success      bool                `json:"success"`
}

type TeamMember struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Username  string             `json:"username" bson:"username"`
	FirstName string             `json:"first_name" bson:"first_name"`
	LastName  string             `json:"last_name" bson:"lastName"`
	Email     string             `json:"email" bson:"email"`
}

type TimeEntry struct {
	UserId     string `json:"userId" bson:"userId"`
	Start      string `json:"start" bson:"start"`
	End        string `json:"end" bson:"end"`
	BreakStart string `json:"breakStart" bson:"breakStart"`
	BreakEnd   string `json:"breakEnd" bson:"breakEnd"`
	Project    string `json:"project" bson:"project"`
}

type Proposal struct {
	UserId    string `json:"userId" bson:"userId"`
	StartDate string `json:"startDate" bson:"startDate"`
	EndDate   string `json:"endDate" bson:"endDate"`
	Approved  bool   `json:"approved" bson:"approved"`
	Type      string `json:"type" bson:"type"`
	//	ID         primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	TimeObject ProposalTimeObject `json:"timeObject" bson:"timeObject"`
}

type ProposalTimeObject struct {
	Duration time.Duration
	Interval *interval.TimeInterval
	//	Err      error
}

type ProposalPayload struct {
	UserId    string `json:"userId" bson:"userId"`
	StartDate string `json:"startDate" bson:"startDate"`
	EndDate   string `json:"endDate" bson:"endDate"`
	Type      string `json:"type" bson:"type"`
}

type ProposalTimeStringObject struct {
	Duration time.Duration
	Interval interval.TimeInterval
	//	Err      error
}
