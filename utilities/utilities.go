package utilities

import (
	"crypto/sha256"
	"errors"
	"example-project/model"
	"github.com/golang-jwt/jwt/v4"
	"github.com/retailify/go-interval"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"strings"
	"time"
)

func ProposalTimeIntersectsProposals(proposal model.Proposal, Arr []model.Proposal) bool {
	for _, p := range Arr {
		p.TimeObject, _ = CreateTimeObject(p.StartDate, p.EndDate)

		if CustomOverlaps(p, proposal) {
			return true
		}
		if (*p.TimeObject.Interval.Start() == *proposal.TimeObject.Interval.Start()) || (*p.TimeObject.Interval.End() == *proposal.TimeObject.Interval.End()) {
			return true
		}

		if proposal.TimeObject.Interval.During(p.TimeObject.Interval) {
			return true
		}
		if p.TimeObject.Interval.During(proposal.TimeObject.Interval) {
			return true
		}
		if p.TimeObject.Interval.Equals(proposal.TimeObject.Interval) {
			return true
		}

	}
	return false
}

func StartDateExceedsEndDate(Arr []model.Proposal) bool {
	for _, p := range Arr {
		p.TimeObject, _ = CreateTimeObject(p.StartDate, p.EndDate)
		if p.TimeObject.Interval.End().Before(*p.TimeObject.Interval.Start()) {
			return true
		}
	}
	return false
}

func CustomOverlaps(p1 model.Proposal, p2 model.Proposal) bool {
	if (*p1.TimeObject.Interval.Start() == *p2.TimeObject.Interval.Start()) && (p1.TimeObject.Interval.End().Before(*p2.TimeObject.Interval.End())) {
		return true
	}

	if (*p2.TimeObject.Interval.Start() == *p1.TimeObject.Interval.Start()) && (p2.TimeObject.Interval.End().Before(*p1.TimeObject.Interval.End())) {
		return true
	}

	return false
}

func CreateTimeObject(start, end string) (model.ProposalTimeObject, error) {
	const shortForm = "2006-Jan-02"
	Start := strings.Split(start, " ")
	End := strings.Split(end, " ")

	Interval, err := interval.MakeTimeIntervalFromStrings(Start[0], End[0], shortForm)
	obj := model.ProposalTimeObject{
		Duration: Interval.Duration(),
		Interval: Interval,
		//		Err:      err,
	}
	return obj, err
}

func CraftProposalFromPayload(payload []model.ProposalPayload) ([]model.Proposal, error) {

	var proposals []model.Proposal
	for _, p := range payload {
		obj, err := CreateTimeObject(p.StartDate, p.EndDate)
		var pStatus = "pending"
		if p.Type == "sickness" {
			pStatus = "approved"
		}
		newProposal := model.Proposal{
			UserId:     p.UserId,
			StartDate:  p.StartDate,
			EndDate:    p.EndDate,
			Status:     pStatus,
			Type:       p.Type,
			TimeObject: obj,
		}
		if err != nil {
			timeIntervalErrMsg := errors.New("Error occured in building the time interval for a new proposal")
			return nil, timeIntervalErrMsg
		}
		proposals = append(proposals, newProposal)
	}

	return proposals, nil
}

func GenerateToken(userid primitive.ObjectID) string {
	claims := jwt.MapClaims{
		"exp":    time.Now().Add(time.Minute * 5).Unix(),
		"iat":    time.Now().Unix(),
		"userID": userid.Hex(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return t
}

func ValidateToken(token string) (*jwt.Token, jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	jwtoken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, nil, err
	}
	return jwtoken, claims, nil
}

func UserUpdateSetter(user model.UpdateUserPayload, userGroup string) (bson.D, error) {
	var setElements bson.D

userLoop:
	for {
		switch {
		case user.FirstName != "":
			setElements = append(setElements, bson.E{Key: "firstname", Value: user.FirstName})
			user.FirstName = ""
		case user.LastName != "":
			setElements = append(setElements, bson.E{Key: "lastname", Value: user.LastName})
			user.LastName = ""
		case user.Email != "":
			setElements = append(setElements, bson.E{Key: "email", Value: user.Email})
			user.Email = ""
		case user.Team != "":
			setElements = append(setElements, bson.E{Key: "team", Value: user.Team})
			user.Team = ""
		case user.TotalWorkingDays != 0:
			setElements = append(setElements, bson.E{Key: "totalWorkingDays", Value: user.TotalWorkingDays})
			user.TotalWorkingDays = 0
		case user.VacationDays != 0:
			setElements = append(setElements, bson.E{Key: "vacationDays", Value: user.VacationDays})
			user.VacationDays = 0
		case user.Username != "":
			setElements = append(setElements, bson.E{Key: "username", Value: user.Username})
			user.Username = ""
		case user.Password != "":
			setElements = append(setElements, bson.E{Key: "password", Value: sha256.Sum256([]byte(user.Password))})
			user.Password = ""
		case userGroup == "admin" && user.Group != "":
			setElements = append(setElements, bson.E{Key: "group", Value: user.Group})
			user.Group = ""
		case len(user.Projects) != 0:
			setElements = append(setElements, bson.E{Key: "projects", Value: user.Projects})
			user.Projects = []string{}
		default:
			break userLoop
		}

	}

	if len(setElements) == 0 {
		return nil, errors.New("no data changed on user")
	}
	setMap := bson.D{
		{"$set", setElements},
	}
	return setMap, nil
}
