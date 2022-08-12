package model

import "time"

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
	Username          string   `json:"username" bson:"username"`
	Password          string   `json:"password" bson:"password"`
	FirstName         string   `json:"first_name" bson:"first_name"`
	LastName          string   `json:"last_name" bson:"lastName"`
	Email             string   `json:"email" bson:"email"`
	Team              string   `json:"team" bson:"team"`
	Projects          []string `json:"projects" bson:"projects"`
	TotalWorkingHours float32  `json:"totalWorkingHours" bson:"totalWorkingHours"`
	VacationDays      int      `json:"vacationDays" bson:"vacationDays"`
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
	UserId    string    `json:"userId" bson:"userId"`
	StartDate time.Time `json:"startDate" bson:"startDate"`
	EndDate   time.Time `json:"endDate" bson:"endDate"`
	Approved  bool      `json:"approved" bson:"approved"`
	Type      string    `json:"type" bson:"type"`
}
