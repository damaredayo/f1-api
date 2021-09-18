package main

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/apognu/gocal"
)

var ErrNilCalander = errors.New("calander is nil")
var ErrNoRace = errors.New("no race")

const F1CALANDER = "http://www.formula1.com/calendar/Formula_1_Official_Calendar.ics"

// We want to cache the .ics file for an hour instead of redownloading it every time a request is made, we can do this by
// implementing a struct which has the calander bytes and the Time object with when it expires

type F1CalanderCache struct {
	Calander *gocal.Gocal
	Expires  time.Time
}

var f1Cache *F1CalanderCache

func updateCalander() {
	resp, err := http.Get(F1CALANDER)
	// If it fails for whatever reason, we can continue and simply reattempt a download next time the cache is checked (on any request)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Parse calander
	c := gocal.NewParser(resp.Body)

	c.Parse()

	// Now that we have the file, we can update the cache
	f1Cache = &F1CalanderCache{Calander: c, Expires: time.Now().Add(time.Hour)}
}

func GetNextRace() (*Race, error) {
	// Check if the cache has expired or not initalized, if it has then attempt to update it
	if f1Cache == nil {
		updateCalander()
	}

	if time.Now().After(f1Cache.Expires) {
		updateCalander()
	}

	// It shouldn't be nil after this but if it is then of course we can't continue.
	if f1Cache.Calander == nil {
		return nil, ErrNilCalander
	}

	// Get the closest event
	var closestEvent gocal.Event
	var closestEventDifference time.Duration
	for _, e := range f1Cache.Calander.Events {
		// If it happened in the past then its obviously not the next race
		if e.Start.Before(time.Now()) {
			continue
		}

		if closestEvent.Summary == "" {
			closestEvent = e
			closestEventDifference = time.Until(*e.Start)
		}

		if time.Until(*e.Start) < closestEventDifference {
			closestEvent = e
			closestEventDifference = time.Until(*e.Start)
		}
	}

	if closestEvent.Summary == "" {
		return nil, ErrNoRace
	}

	splitSummary := strings.Split(closestEvent.Summary, "-")
	raceType := strings.TrimSpace(splitSummary[len(splitSummary)-1])
	raceTitle := strings.TrimSpace(splitSummary[len(splitSummary)-2])

	race := &Race{
		Status:    closestEvent.Status,
		Name:      raceTitle,
		Type:      raceType,
		StartDate: closestEvent.Start.Format(time.RFC3339),
		EndDate:   closestEvent.End.Format(time.RFC3339),
	}

	return race, nil
}

func GetCurrentRace() (*Race, error) {
	// Check if the cache has expired or not initalized, if it has then attempt to update it
	if f1Cache == nil {
		updateCalander()
	}

	if time.Now().After(f1Cache.Expires) {
		updateCalander()
	}

	// It shouldn't be nil after this but if it is then of course we can't continue.
	if f1Cache.Calander == nil {
		return nil, ErrNilCalander
	}

	// Get the closest event
	var currentEvent gocal.Event

	for _, e := range f1Cache.Calander.Events {
		// If it happened in the past then its obviously not the next race
		if e.Start.Before(time.Now()) {
			continue
		}

		currentTime := time.Now()
		if e.Start.Before(currentTime) && e.End.After(currentTime) {
			currentEvent = e
			break
		}
	}

	if currentEvent.Summary == "" {
		return nil, ErrNoRace
	}

	splitSummary := strings.Split(currentEvent.Summary, "-")
	raceType := strings.TrimSpace(splitSummary[len(splitSummary)-1])
	raceTitle := strings.TrimSpace(splitSummary[len(splitSummary)-2])

	race := &Race{
		Status:    currentEvent.Status,
		Name:      raceTitle,
		Type:      raceType,
		StartDate: currentEvent.Start.Format(time.RFC3339),
		EndDate:   currentEvent.End.Format(time.RFC3339),
	}

	return race, nil
}
