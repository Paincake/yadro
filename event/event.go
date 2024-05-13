package event

import (
	"fmt"
	"github.com/Paincake/yadro/event/club"
	"io"
	"time"
)

type ClubEvent interface {
	ProcessEvent(club *club.Club)
	logEvent()
}

type ClientArrivalEvent struct {
	ClientID     string
	ArrTime      time.Time
	ResultWriter io.Writer
}

func (e *ClientArrivalEvent) logEvent() {
	defer io.WriteString(e.ResultWriter, fmt.Sprintf("%s %s %s\n", e.ArrTime.Format(club.TimeFormat_hhmm), club.CLIENT_ARRIVING_IN, e.ClientID))
}

func (e *ClientArrivalEvent) ProcessEvent(club *club.Club) {
	if club.ClientExists(e.ClientID) {
		//ID 13 YouShallNotPass
		ev := ErrorEvent{
			Error:        ClientPresentError{},
			EventTime:    e.ArrTime,
			ClientID:     e.ClientID,
			ResultWriter: e.ResultWriter,
		}
		e.logEvent()
		ev.ProcessEvent(club)
		return
	}
	if e.ArrTime.Sub(club.OpenTime) < 0 {
		//ID 13 NotOpenYet
		ev := ErrorEvent{
			Error:        ClubClosedError{},
			EventTime:    e.ArrTime,
			ClientID:     e.ClientID,
			ResultWriter: e.ResultWriter,
		}
		e.logEvent()
		ev.ProcessEvent(club)
		return
	}
	club.AddClient(e.ClientID, e.ArrTime)
	e.logEvent()
}

type ClientTablePickEvent struct {
	ClientID     string
	Table        int
	EventTime    time.Time
	ResultWriter io.Writer
}

func (e *ClientTablePickEvent) logEvent() {
	io.WriteString(e.ResultWriter, fmt.Sprintf("%s %s %s %d\n", e.EventTime.Format(club.TimeFormat_hhmm), club.CLIENT_TABLE_USING_IN, e.ClientID, e.Table))
}

func (e *ClientTablePickEvent) ProcessEvent(club *club.Club) {
	if !club.ClientExists(e.ClientID) {
		//ID 13 ClientUnknown
		ev := ErrorEvent{
			Error:        ClientNotPresentError{},
			EventTime:    e.EventTime,
			ClientID:     e.ClientID,
			ResultWriter: e.ResultWriter,
		}
		e.logEvent()
		ev.ProcessEvent(club)
		return
	}
	table := club.PickTable(e.ClientID, e.Table, e.EventTime)
	if table == 0 {
		//ID 13 PlaceIsBusy
		ev := ErrorEvent{
			Error:        NoTableAvailableError{},
			EventTime:    e.EventTime,
			ClientID:     e.ClientID,
			ResultWriter: e.ResultWriter,
		}
		e.logEvent()
		ev.ProcessEvent(club)
		return
	}
	e.logEvent()
}

type ClientWaitingEvent struct {
	ClientID     string
	EventTime    time.Time
	ResultWriter io.Writer
}

func (e *ClientWaitingEvent) logEvent() {
	io.WriteString(e.ResultWriter, fmt.Sprintf("%s %s %s\n", e.EventTime.Format(club.TimeFormat_hhmm), club.CLIENT_WAITING_IN, e.ClientID))
}

func (e *ClientWaitingEvent) ProcessEvent(club *club.Club) {
	if !club.ClientExists(e.ClientID) {
		//ID 13 ClientUnknown
		ev := ErrorEvent{
			Error:        ClientNotPresentError{},
			EventTime:    e.EventTime,
			ClientID:     e.ClientID,
			ResultWriter: e.ResultWriter,
		}
		e.logEvent()
		ev.ProcessEvent(club)
		return
	}
	if !club.IsBusy() {
		//ID 13 ICanWaitNoLonger
		ev := ErrorEvent{
			Error:        WaitingError{},
			EventTime:    e.EventTime,
			ClientID:     e.ClientID,
			ResultWriter: e.ResultWriter,
		}
		e.logEvent()
		ev.ProcessEvent(club)
		return
	} else {
		err := club.EnqueueClient(e.ClientID)
		if err != nil {
			// ID 11 event
			ev := ClientLeftEvent{
				ClientID:     e.ClientID,
				LeavingTime:  e.EventTime,
				ResultWriter: e.ResultWriter,
			}
			e.logEvent()
			ev.ProcessEvent(club)
			return
		}
		e.logEvent()
	}
}

type ClientLeavingEvent struct {
	ClientID     string
	LeavingTime  time.Time
	ResultWriter io.Writer
}

func (e *ClientLeavingEvent) logEvent() {
	io.WriteString(e.ResultWriter, fmt.Sprintf("%s %s %s\n", e.LeavingTime.Format(club.TimeFormat_hhmm), club.CLIENT_LEAVING_IN, e.ClientID))
}

func (e *ClientLeavingEvent) ProcessEvent(club *club.Club) {
	if !club.ClientExists(e.ClientID) {
		//ID 13 ClientUnknown
		ev := ErrorEvent{
			Error:        ClientNotPresentError{},
			EventTime:    e.LeavingTime,
			ClientID:     e.ClientID,
			ResultWriter: e.ResultWriter,
		}
		e.logEvent()
		ev.ProcessEvent(club)
		return
	}

	table := club.RemoveClient(e.ClientID, e.LeavingTime)

	//ID 12 event
	if table != 0 {
		ev := ClientDequeuedEvent{
			EventTime:    e.LeavingTime,
			Table:        table,
			ResultWriter: e.ResultWriter,
		}
		e.logEvent()
		ev.ProcessEvent(club)
		return
	}
	e.logEvent()
}

type ClientLeftEvent struct {
	ClientID     string
	LeavingTime  time.Time
	ResultWriter io.Writer
}

func (e *ClientLeftEvent) logEvent() {
	io.WriteString(e.ResultWriter, fmt.Sprintf("%s %s %s\n", e.LeavingTime.Format(club.TimeFormat_hhmm), club.CLIENT_LEAVING_OUT, e.ClientID))
}

func (e *ClientLeftEvent) ProcessEvent(club *club.Club) {
	club.RemoveClient(e.ClientID, e.LeavingTime)
	e.logEvent()
}

type ClientDequeuedEvent struct {
	EventTime    time.Time
	Table        int
	ResultWriter io.Writer
	clientID     string
}

func (e *ClientDequeuedEvent) logEvent() {
	if e.clientID != "" {
		io.WriteString(e.ResultWriter, fmt.Sprintf("%s %s %s %d\n", e.EventTime.Format(club.TimeFormat_hhmm), club.CLIENT_TABLE_USING_OUT, e.clientID, e.Table))
	}
}

func (e *ClientDequeuedEvent) ProcessEvent(club *club.Club) {
	defer e.logEvent()
	e.clientID = club.DequeueClient(e.Table, e.EventTime)
}

type ErrorEvent struct {
	Error        error
	EventTime    time.Time
	ClientID     string
	ResultWriter io.Writer
}

func (e *ErrorEvent) logEvent() {
	io.WriteString(e.ResultWriter, fmt.Sprintf("%s %s %s\n", e.EventTime.Format(club.TimeFormat_hhmm), club.ERROR, e.Error.Error()))
}

func (e *ErrorEvent) ProcessEvent(club *club.Club) {
	e.logEvent()
}

type ClubClosingEvent struct {
	closeTime    time.Time
	profits      []club.TableProfit
	ResultWriter io.Writer
}

func (e *ClubClosingEvent) ProcessEvent(club *club.Club) {
	clients := club.GetRemainingClientsSorted()
	for _, c := range clients {
		ev := ClientLeftEvent{
			ClientID:     c,
			LeavingTime:  club.CloseTime,
			ResultWriter: e.ResultWriter,
		}
		ev.ProcessEvent(club)
	}
	e.profits = club.CountProfit()
	e.logEvent()
}

func (e *ClubClosingEvent) logEvent() {
	io.WriteString(e.ResultWriter, fmt.Sprintf("%s\n", e.closeTime.Format(club.TimeFormat_hhmm)))
	for _, profit := range e.profits {
		hours := profit.MinutesInUse / 60
		minutes := profit.MinutesInUse % 60
		timeInUse := fmt.Sprintf("%.2d:%.2d", hours, minutes)
		io.WriteString(e.ResultWriter, fmt.Sprintf("%d %d %s\n", profit.TableNum, profit.Profit, timeInUse))
	}
}

type ClubOpeningEvent struct {
	ResultWriter io.Writer
	openTime     time.Time
}

func (e *ClubOpeningEvent) ProcessEvent(c *club.Club) {
	e.logEvent()
}

func (e *ClubOpeningEvent) logEvent() {
	io.WriteString(e.ResultWriter, fmt.Sprintf("%s\n", e.openTime.Format(club.TimeFormat_hhmm)))
}
