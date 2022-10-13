package utilities

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"example-project/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/retailify/go-interval"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"net/http"
	"os"
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

func CreateTimeObject(start, end time.Time) (model.ProposalTimeObject, error) {

	Interval, err := interval.MakeTimeInterval(&start, &end)
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
		startObj, err := time.Parse("2006-Jan-02", p.StartDate)
		if err != nil {
			return proposals, err
		}
		endObj, err := time.Parse("2006-Jan-02", p.EndDate)
		if err != nil {
			return proposals, err
		}
		obj, err := CreateTimeObject(startObj, endObj)
		var pStatus = "pending"
		if p.Type == "sickness" {
			pStatus = "approved"
		}

		newProposal := model.Proposal{
			UserId:     p.UserId,
			StartDate:  startObj,
			EndDate:    endObj,
			Status:     pStatus,
			Type:       p.Type,
			TimeObject: obj,
			Timestamp:  time.Now(),
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
		case user.HoursPerWeek != 0:
			setElements = append(setElements, bson.E{Key: "totalWorkingDays", Value: user.HoursPerWeek})
			user.HoursPerWeek = 0
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

func GetWeekdaysBetween(start, end time.Time) int {
	days := 0
	for end.After(start) {

		if start.Weekday().String() != "Saturday" && start.Weekday().String() != "Sunday" {
			days++
		}
		start = start.Add(time.Hour * 24)

	}
	return days
}

func FormGetAllProposalsFilter(user model.UserPayload, ctx *gin.Context) (bson.M, bson.D) {
	filter := bson.M{}
	var filterTimeArray []bson.M
	sort := bson.D{{"timestamp", 1}}
	typeQuery, typeOK := ctx.GetQuery("type")
	var startTime time.Time
	var endTime time.Time
	filter["userId"] = user.ID.Hex()
	if typeOK {
		filter["type"] = typeQuery
	}
	statusQuery, statusOK := ctx.GetQuery("status")
	if statusOK {
		filter["status"] = statusQuery
	}

	startQuery, startOK := ctx.GetQuery("start")
	endQuery, endOK := ctx.GetQuery("end")
	if startOK {
		sTime, err := time.Parse(time.RFC3339, startQuery)
		if err == nil {
			startTime = sTime
		}

	} else {
		startTime = user.EntryTime
	}

	if endOK {
		eTime, err := time.Parse(time.RFC3339, endQuery)
		if err == nil {
			endTime = eTime
		}
	} else {
		endTime = time.Now()
	}

	filterTimeArray = []bson.M{
		{
			"$and": []bson.M{
				{
					"startDate": bson.M{"$gt": startTime},
				},
				{
					"endDate": bson.M{"$lt": endTime},
				},
			},
		},
		{
			"$and": []bson.M{
				{
					"startDate": bson.M{"$gt": startTime},
				},
				{
					"startDate": bson.M{"$lt": endTime},
				},
			},
		},
		{
			"$and": []bson.M{
				{
					"endDate": bson.M{"$gt": startTime},
				},
				{
					"endDate": bson.M{"$lt": endTime},
				},
			},
		},
		{
			"$and": []bson.M{
				{
					"startDate": bson.M{"$lt": startTime},
				},
				{
					"endDate": bson.M{"$gt": endTime},
				},
			},
		},
		{
			"startDate": startTime,
		},
		{
			"endDate": endTime,
		},
		{
			"endDate": startTime,
		},
		{
			"startDate": endTime,
		},
	}
	filter["$or"] = filterTimeArray

	sortingQuery, sortingOK := ctx.GetQuery("sort")
	if sortingOK {
		if sortingQuery == "desc" {
			sort = bson.D{{"timestamp", -1}}
		}
		if sortingQuery == "asce" {
			sort = bson.D{{"timestamp", 1}}
		}
	}
	return filter, sort
}

func FormGetTimeEntryFilter(ctx *gin.Context) (bson.M, error) {
	filter := bson.M{}
	userID, userIDOK := ctx.GetQuery("userid")
	if !userIDOK {
		return filter, errors.New("no userid have been supplied as a query")
	}
	filter["userId"] = userID
	start, startOK := ctx.GetQuery("start")
	end, endOK := ctx.GetQuery("end")

	if startOK {
		startTime, err := time.Parse(time.RFC3339, start)
		if err != nil {
			return filter, err
		}
		filter["start"] = bson.M{"$gt": startTime}
	}
	if endOK {
		endTime, err := time.Parse(time.RFC3339, end)
		if err != nil {
			return filter, err
		}
		filter["end"] = bson.M{"$lt": endTime}
	}
	if startOK && endOK {
		startTime, err := time.Parse(time.RFC3339, start)
		if err != nil {
			return filter, err
		}
		endTime, err := time.Parse(time.RFC3339, end)
		if err != nil {
			return filter, err
		}
		filter["start"] = bson.M{"$gt": startTime}
		filter["end"] = bson.M{"$lt": endTime}
	}
	return filter, nil
}

func CalculateRequiredWorkingHours(user model.UserPayload, proposals []model.Proposal, ctx *gin.Context) (float64, error) {
	var hoursPerDay = user.HoursPerWeek / 5
	var startTime time.Time
	var endTime time.Time
	var totalProposalDays float64 = 0
	start, startOK := ctx.GetQuery("start")
	if startOK {
		var err error
		startTime, err = time.Parse(time.RFC3339, start)
		if err != nil {
			return 0, err
		}
	} else {
		startTime = user.EntryTime
	}
	end, endOK := ctx.GetQuery("end")
	if endOK {
		var err error
		endTime, err = time.Parse(time.RFC3339, end)
		if err != nil {
			return 0, err
		}
	} else {
		endTime = time.Now()
	}
	var totalHours = float64(GetWeekdaysBetween(startTime, endTime)+1) * hoursPerDay
	totalHolidays, err := GetPublicHolidays(startTime, endTime)
	if err != nil {
		return 0, err
	}
	totalHours = totalHours - (float64(totalHolidays) * hoursPerDay)
	for _, proposal := range proposals {
		var proposalTotalDays float64 = 0
		var proposalStartOffset float64 = 0
		var proposalEndOffset float64 = 0
		proposalTotalDays = float64(GetWeekdaysBetween(proposal.StartDate, proposal.EndDate))
		if proposal.StartDate.Before(startTime) {
			proposalStartOffset = float64(GetWeekdaysBetween(proposal.StartDate, startTime))
		} else if (proposal.StartDate.Equal(endTime) || proposal.EndDate.Equal(startTime)) && !proposal.StartDate.Before(startTime) {
			proposalStartOffset += proposalTotalDays - 1
		}
		if proposal.EndDate.After(endTime) {
			proposalEndOffset = float64(GetWeekdaysBetween(proposal.StartDate, endTime))
		}
		totalProposalDays = totalProposalDays + (proposalTotalDays - (proposalStartOffset + proposalEndOffset)) + 1
	}
	var totalProposalHours = totalProposalDays * hoursPerDay
	return totalHours - totalProposalHours, nil
}

func GetPublicHolidays(start, end time.Time) (int, error) {
	var totalDays int = 0
	year := time.Now().Year()
	url := fmt.Sprintf("https://get.api-feiertage.de/?years=%v&states=hb", year)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	body := resp.Body
	rawData, err := ioutil.ReadAll(body)
	if err != nil {
		return 0, err
	}
	var responseInterface map[string]interface{}
	json.Unmarshal(rawData, &responseInterface)
	holidays := responseInterface["feiertage"].([]interface{})
	for _, hol := range holidays {
		obj := hol.(map[string]interface{})
		time, err := time.Parse("2006-01-02", fmt.Sprint(obj["date"]))
		if err != nil {
			return 0, err
		}
		if time.After(start) && time.Before(end) {
			totalDays++
		}
	}
	return totalDays, nil
}
