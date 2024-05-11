package event

import (
	"fmt"
	"io"
	"time"
)

type ClubEvent interface {
	ProcessEvent(club *Club)
	logEvent()
}

type ClientArrivalEvent struct {
	ClientID     string
	ArrTime      time.Time
	ResultWriter io.Writer
}

func (e *ClientArrivalEvent) logEvent() {
	defer io.WriteString(e.ResultWriter, fmt.Sprintf("%s %s %s\n", e.ArrTime.Format(hhmm), CLIENT_ARRIVING_IN, e.ClientID))
}

func (e *ClientArrivalEvent) ProcessEvent(club *Club) {
	defer e.logEvent()
	if club.ClientExists(e.ClientID) {
		//ID 13 YouShallNotPass
		ev := ErrorEvent{
			Error:        ClientPresentError{},
			EventTime:    e.ArrTime,
			ClientID:     e.ClientID,
			ResultWriter: e.ResultWriter,
		}
		ev.ProcessEvent(club)
		return
	}
	if e.ArrTime.Sub(club.OpenTime) < 0 || e.ArrTime.Sub(club.CloseTime) > 0 {
		//ID 13 NotOpenYet
		ev := ErrorEvent{
			Error:        ClubClosedError{},
			EventTime:    e.ArrTime,
			ClientID:     e.ClientID,
			ResultWriter: e.ResultWriter,
		}
		ev.ProcessEvent(club)
		return
	}
	club.AddClient(e.ClientID, e.ArrTime)
}

type ClientTablePickEvent struct {
	ClientID     string
	Table        int
	EventTime    time.Time
	ResultWriter io.Writer
}

func (e *ClientTablePickEvent) logEvent() {
	io.WriteString(e.ResultWriter, fmt.Sprintf("%s %s %s %d\n", e.EventTime.Format(hhmm), CLIENT_TABLE_USING_IN, e.ClientID, e.Table))
}

func (e *ClientTablePickEvent) ProcessEvent(club *Club) {
	defer e.logEvent()
	if !club.ClientExists(e.ClientID) {
		//ID 13 ClientUnknown
		ev := ErrorEvent{
			Error:        ClientNotPresentError{},
			EventTime:    e.EventTime,
			ClientID:     e.ClientID,
			ResultWriter: e.ResultWriter,
		}
		ev.ProcessEvent(club)
		return
	}
	table := club.PickTable(e.ClientID, e.Table)
	if table == 0 {
		//ID 13 PlaceIsBusy
		ev := ErrorEvent{
			Error:        NoTableAvailableError{},
			EventTime:    e.EventTime,
			ClientID:     e.ClientID,
			ResultWriter: e.ResultWriter,
		}
		ev.ProcessEvent(club)
		return
	}
}

type ClientWaitingEvent struct {
	ClientID     string
	EventTime    time.Time
	ResultWriter io.Writer
}

func (e *ClientWaitingEvent) logEvent() {
	io.WriteString(e.ResultWriter, fmt.Sprintf("%s %s %s\n", e.EventTime.Format(hhmm), CLIENT_WAITING_IN, e.ClientID))
}

func (e *ClientWaitingEvent) ProcessEvent(club *Club) {
	defer e.logEvent()
	if !club.IsBusy() {
		//ID 13 ICanWaitNoLonger
		ev := ErrorEvent{
			Error:        WaitingError{},
			EventTime:    e.EventTime,
			ClientID:     e.ClientID,
			ResultWriter: e.ResultWriter,
		}
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
			ev.ProcessEvent(club)
		}
	}
}

type ClientLeavingEvent struct {
	ClientID     string
	LeavingTime  time.Time
	ResultWriter io.Writer
}

func (e *ClientLeavingEvent) logEvent() {
	io.WriteString(e.ResultWriter, fmt.Sprintf("%s %s %s\n", e.LeavingTime.Format(hhmm), CLIENT_LEAVING_OUT, e.ClientID))
}

func (e *ClientLeavingEvent) ProcessEvent(club *Club) {
	defer e.logEvent()
	if !club.ClientExists(e.ClientID) {
		//ID 13 ClientUnknown
		ev := ErrorEvent{
			Error:        ClientNotPresentError{},
			EventTime:    e.LeavingTime,
			ClientID:     e.ClientID,
			ResultWriter: e.ResultWriter,
		}
		ev.ProcessEvent(club)
		return
	}
	table := club.RemoveClient(e.ClientID, e.LeavingTime)

	//ID 12 event
	ev := ClientDequeuedEvent{
		EventTime:    e.LeavingTime,
		Table:        table,
		ResultWriter: e.ResultWriter,
	}
	ev.ProcessEvent(club)
}

type ClientLeftEvent struct {
	ClientID     string
	LeavingTime  time.Time
	ResultWriter io.Writer
}

func (e *ClientLeftEvent) logEvent() {
	io.WriteString(e.ResultWriter, fmt.Sprintf("%s %s %s\n", e.LeavingTime.Format(hhmm), CLIENT_LEAVING_OUT, e.ClientID))
}

func (e *ClientLeftEvent) ProcessEvent(club *Club) {
	defer e.logEvent()
	club.RemoveClient(e.ClientID, e.LeavingTime)
}

type ClientDequeuedEvent struct {
	EventTime    time.Time
	Table        int
	ResultWriter io.Writer
	clientID     string
}

func (e *ClientDequeuedEvent) logEvent() {
	if e.clientID != "" {
		io.WriteString(e.ResultWriter, fmt.Sprintf("%s %s %s %d\n", e.EventTime.Format(hhmm), CLIENT_TABLE_USING_OUT, e.clientID, e.Table))
	}
}

func (e *ClientDequeuedEvent) ProcessEvent(club *Club) {
	defer e.logEvent()
	e.clientID = club.DequeueClient(e.Table)
}

type ErrorEvent struct {
	Error        error
	EventTime    time.Time
	ClientID     string
	ResultWriter io.Writer
}

func (e *ErrorEvent) logEvent() {
	io.WriteString(e.ResultWriter, fmt.Sprintf("%s %s %s\n", e.EventTime.Format(hhmm), ERROR, e.Error.Error()))
}

func (e *ErrorEvent) ProcessEvent(club *Club) {
	e.logEvent()
}

type ClubClosingEvent struct {
	profits      []ClientProfit
	ResultWriter io.Writer
}

func (e *ClubClosingEvent) ProcessEvent(club *Club) {
	defer e.logEvent()
	clients := club.GetClientsSorted()
	for _, c := range clients {
		ev := ClientLeftEvent{
			ClientID:     c,
			LeavingTime:  club.CloseTime,
			ResultWriter: e.ResultWriter,
		}
		ev.ProcessEvent(club)
	}
	e.profits = club.CountProfit()
}

func (e *ClubClosingEvent) logEvent() {
	for _, profit := range e.profits {
		io.WriteString(e.ResultWriter, fmt.Sprintf("%s %s %d", profit.ClientID, profit.Profit, profit.Time))
	}
}
