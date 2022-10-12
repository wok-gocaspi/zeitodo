package model

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/retailify/go-interval"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"time"
)

type DbConfig struct {
	URL      string
	Database string
}

type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Username     string             `json:"username" bson:"username"`
	Password     [32]byte           `json:"password" bson:"password"`
	FirstName    string             `json:"firstname" bson:"firstname"`
	LastName     string             `json:"lastname" bson:"lastname"`
	Email        string             `json:"email" bson:"email"`
	Team         string             `json:"team" bson:"team"`
	Projects     []string           `json:"projects" bson:"projects"`
	HoursPerWeek float64            `json:"hoursPerWeek" bson:"hoursPerWeek"`
	VacationDays int                `json:"vacationDays" bson:"vacationDays"`
	Group        string             `json:"group" bson:"group"`
	EntryTime    time.Time          `json:"entryTime" bson:"entryTime"`
}

type UpdateUserPayload struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Username     string             `json:"username" bson:"username"`
	Password     string             `json:"password"`
	FirstName    string             `json:"firstname" bson:"firstname"`
	LastName     string             `json:"lastname" bson:"lastname"`
	Email        string             `json:"email" bson:"email"`
	Team         string             `json:"team" bson:"team"`
	Projects     []string           `json:"projects" bson:"projects"`
	HoursPerWeek float64            `json:"hoursPerWeek" bson:"hoursPerWeek"`
	VacationDays int                `json:"vacationDays" bson:"vacationDays"`
	Group        string             `json:"group" bson:"group"`
}

type UserSignupPayload struct {
	Username  string `json:"username" bson:"username" key:"required"`
	Password  string `json:"password" bson:"password" key:"required"`
	FirstName string `json:"firstname" bson:"firstname" key:"required"`
	LastName  string `json:"lastname" bson:"lastname" key:"required"`
	Email     string `json:"email" bson:"email" key:"required"`
}

type UserSignup struct {
	Username     string    `json:"username" bson:"username"`
	Password     [32]byte  `json:"password" bson:"password"`
	FirstName    string    `json:"firstname" bson:"firstname"`
	LastName     string    `json:"lastname" bson:"lastname"`
	Email        string    `json:"email" bson:"email"`
	Group        string    `json:"group" bson:"group"`
	EntryTime    time.Time `json:"entryTime" bson:"entryTime"`
	HoursPerWeek float64   `json:"hoursPerWeek" bson:"hoursPerWeek"`
	VacationDays int       `json:"vacationDays" bson:"vacationDays"`
}

type UserSignupResult struct {
	Username  string             `json:"username" bson:"username"`
	FirstName string             `json:"firstname" bson:"firstname"`
	LastName  string             `json:"lastname" bson:"lastname"`
	Email     string             `json:"email" bson:"email"`
	ID        primitive.ObjectID `json:"id"`
}

type UserPayload struct {
	Username     string             `json:"username" bson:"username"`
	FirstName    string             `json:"firstname" bson:"firstname"`
	LastName     string             `json:"lastname" bson:"lastname"`
	Email        string             `json:"email" bson:"email"`
	Team         string             `json:"team" bson:"team"`
	Projects     []string           `json:"projects" bson:"projects"`
	HoursPerWeek float64            `json:"hoursPerWeek" bson:"hoursPerWeek"`
	VacationDays int                `json:"vacationDays" bson:"vacationDays"`
	Group        string             `json:"group" bson:"group"`
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	EntryTime    time.Time          `json:"entryTime" bson:"entryTime"`
}

type UserUpdateResult struct {
	UpdateResult *mongo.UpdateResult `json:"result"`
	ID           primitive.ObjectID  `json:"id"`
	Success      bool                `json:"success"`
	Error        string              `json:"error"`
}

type UserAuthPayload struct {
	Username string `json:"username"`
	Password string `json:"password" key:"required"`
}

type TeamMember struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Username  string             `json:"username" bson:"username"`
	FirstName string             `json:"firstname" bson:"firstname"`
	LastName  string             `json:"lastname" bson:"lastname"`
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
	Status    string    `json:"status" bson:"status"`
	Type      string    `json:"type" bson:"type"`
	//	ID         primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	TimeObject ProposalTimeObject `json:"timeObject" bson:"timeObject"`
	Timestamp  time.Time          `json:"timestamp" bson:"timestamp"`
}

type ProposalsByUser struct {
	Userid            primitive.ObjectID `json:"userid"`
	FirstName         string             `json:"firstname" bson:"firstname"`
	LastName          string             `json:"lastname" bson:"lastname"`
	Username          string             `json:"username" bson:"username"`
	Email             string             `json:"email" bson:"email"`
	SicknessProposals []Proposal         `json:"sicknessProposals"`
	VacationProposals []Proposal         `json:"vacationProposals"`
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

type AbsenceObject struct {
	VacationDays      int `json:"vacation"`
	TotalVacationDays int `json:"totalVacation"`
	SicknessDays      int `json:"sickness"`
}

type Permission struct {
	Whitelist   []string
	Uri         string
	Methods     []string
	Group       string
	GetSameUser bool
}
type PermissionList struct {
	Permissions []Permission
}
type WorkingHoursPayload struct {
	Projects map[string]float64 `json:"projects"`
	Required float64            `json:"required"`
	Actual   float64            `json:"actual"`
}

func (pl PermissionList) AddPermission(permission Permission) {
	pl.Permissions = append(pl.Permissions, permission)
}

func (pl PermissionList) CheckPolicy(ctx *gin.Context) (bool, error) {

	group := ctx.GetString("group")
	userid := ctx.GetString("userid")

	method := ctx.Request.Method
	url := ctx.Request.URL

	for _, p := range pl.Permissions {
		if strings.HasPrefix(url.String(), p.Uri) && group == p.Group {
			if contains(p.Whitelist, url.String()) {
				return true, nil
			}
			for _, pmethod := range p.Methods {

				if method == pmethod && ((method == "GET" || method == "DELETE" || method == "PUT" || method == "PATCH") && p.GetSameUser) {
					urlID, _ := ctx.Params.Get("id")

					if urlID == userid {
						return true, nil
					} else {
						return false, errors.New("requesting user data of other users is not allowed")
					}

				} else if method == pmethod {

					return true, nil
				}

			}

		}
	}
	return false, errors.New("url dosent match to any permission, deny request...")
}

//***************************************

func (pl PermissionList) IsSameUser(ctx *gin.Context, userID string) bool {

	req := ctx.GetString("userid")

	if userID == req {

		return true
	} else {
		return false
	}

}

//********************************************

func contains(s []string, str string) bool {

	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false

}
