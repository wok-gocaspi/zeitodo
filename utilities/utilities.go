package utilities

import (
	"example-project/model"
	"github.com/retailify/go-interval"
	"strings"
)

func ProposalTimeIntersectsProposals(proposal model.Proposal, Arr []model.Proposal) bool {
	for _, p := range Arr {
		p.TimeObject, _ = CreateTimeObject(p.StartDate, p.EndDate)

		if (*p.TimeObject.Interval.Start() == *proposal.TimeObject.Interval.Start()) || (*p.TimeObject.Interval.End() == *proposal.TimeObject.Interval.End()) {
			return true
		}

		if customOverlaps(p, proposal) {
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

func customOverlaps(p1 model.Proposal, p2 model.Proposal) bool {
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
