package views

import (
	"calendar/models"
	"calendar/postgresql/queries"
	"calendar/restapi/operations"
	"calendar/utils"
	"database/sql"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"log"
	"sort"
	"time"
)

const scanningIntervalMonths = 3 // month

type Interval struct {
	start, end time.Time
	repeat     string
}

func FindSlot(i *[]Interval, slot int64) (*time.Time, *time.Time, bool) {
	maxEnd := (*i)[0].end
	for _, interval := range *i {
		if interval.start.Sub(maxEnd).Minutes() >= float64(slot) {
			return &maxEnd, &interval.start, true
		}
		if interval.end.After(maxEnd) {
			maxEnd = interval.end
		}
	}
	return nil, nil, false
}

func NewUsersFreeSlotHandler(dbClient *sqlx.DB) operations.PostUsersFreeSlotHandlerFunc {
	return func(params operations.PostUsersFreeSlotParams) middleware.Responder {
		from := time.Time(*params.Body.From)
		scanningEnd := from.AddDate(0, scanningIntervalMonths, 0)
		slotInterval := *params.Body.SlotIntervalMin

		intervals := []Interval{}

		rows, err := dbClient.Query(queries.UsersFreeSlotSelect,
			pq.Array(params.Body.UserIds), strfmt.DateTime(from), scanningEnd,
		)
		if err != nil {
			log.Print("Error while fetching user events: ", err.Error())
			return operations.NewGetUserEventsInternalServerError()
		}
		for rows.Next() {
			var st, end strfmt.DateTime
			var repeat sql.NullString
			if err := rows.Scan(&st, &end, &repeat); err != nil {
				log.Print("Error while scanning intervals: ", err.Error())
				return operations.NewGetUserEventsInternalServerError()
			}
			intervals = append(intervals,
				Interval{time.Time(st), time.Time(end), repeat.String})
		}
		sort.SliceStable(intervals, func(i, j int) bool {
			return intervals[i].start.Before(intervals[j].start)
		})

		for start := from; start.Before(scanningEnd); start = start.AddDate(0, 0, 1) {
			end := start.AddDate(0, 0, 1)
			dayIntervals := []Interval{{start, start, ""}}
			for _, interval := range intervals {
				if utils.CheckEvent(interval.start, interval.end, start, end, interval.repeat) {
					dayIntervals = append(dayIntervals, interval)
				}
			}
			dayIntervals = append(dayIntervals, Interval{end, end, ""})
			slotSt, slotEnd, ok := FindSlot(&dayIntervals, slotInterval)
			if ok {
				slotStRet := strfmt.DateTime(*slotSt)
				slotEndRet := strfmt.DateTime(*slotEnd)
				return &operations.PostUsersFreeSlotOK{Payload: &models.UsersFreeSlotResponse{TimeStart: &slotStRet, TimeEnd: &slotEndRet}}
			}
		}
		return &operations.PostUsersFreeSlotNotFound{}
	}
}
