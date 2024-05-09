package event

import (
	"container/list"
	"time"
)

const (
	CLIENT_ARRIVING_IN     = "1"
	CLIENT_TABLE_USING_IN  = "2"
	CLIENT_WAITING_IN      = "3"
	CLIENT_LEAVING_IN      = "4"
	CLIENT_LEAVING_OUT     = "11"
	CLIENT_TABLE_USING_OUT = "12"
	ERROR                  = "13"

	hhmm = "15:04"
)

type visitor struct {
	id   string
	time time.Time
}

type EventHandler struct {
	HourCost       int
	OpenTime       time.Time
	CloseTime      time.Time
	visitors       map[visitor]int
	queuedVisitors *list.List
	tableMap       map[int]bool
	tables         int
	tablesUsed     int
}

func NewEventHandler(tables int, hourCost int, openTime time.Time, closeTime time.Time) *EventHandler {
	return &EventHandler{
		tables:         tables,
		HourCost:       hourCost,
		OpenTime:       openTime,
		CloseTime:      closeTime,
		tableMap:       make(map[int]bool),
		visitors:       make(map[visitor]int),
		queuedVisitors: list.New(),
	}
}

func (eh *EventHandler) HandleInboundEvent(event []string) {
	//if len(event) < 3 {
	//	eh.HandleOutboundEvent(ERROR)
	//}
	//time, err := newWorkTime(event[0])
	//if err != nil {
	//	eh.HandleOutboundEvent(ERROR)
	//}
	//eventID := event[1]
	//clientID := event[2]
	//if len(event) == 4 {
	//	table, err := strconv.Atoi(event[3])
	//	if err != nil {
	//		eh.handleInputError(event)
	//	} else {
	//
	//	}
	//}
}

func (eh *EventHandler) handleInputError(event []string) {

}

func (eh *EventHandler) handleTableUse(client visitor, table int) error {
	if _, ok := eh.visitors[client]; !ok {
		return ClientNotPresentError{}
	}
	v, ok := eh.tableMap[table]
	if !ok {
		eh.tableMap[eh.visitors[client]] = false
		eh.tableMap[table] = true
		eh.visitors[client] = table
	}
	if v {
		return NoTableAvailableError{}
	}
	return nil
}

func (eh *EventHandler) handleClientArrival(client visitor) error {
	if _, ok := eh.visitors[client]; ok {
		return ClientPresentError{}
	}
	if client.time.Sub(eh.OpenTime).Hours() < 0 || client.time.Sub(eh.CloseTime) > 0 {
		return ClubClosedError{}
	}
	eh.visitors[client] = 0
	return nil
}

func (eh *EventHandler) handleClientWaiting(client visitor) error {
	if eh.tables > eh.tablesUsed {
		return WaitingError{}
	}
	if eh.queuedVisitors.Len() > eh.tables {
		return QueueOverflowError{}
	}
	eh.queuedVisitors.PushBack(client)
	return nil
}

func (eh *EventHandler) handleClientLeaving(client visitor) (int, error) {
	if _, ok := eh.visitors[client]; !ok {
		return 0, ClientNotPresentError{}
	}
	table := eh.visitors[client]
	delete(eh.visitors, client)
	return table, nil
}

func (eh *EventHandler) generateOutboundClientLeaving(client visitor) {
	delete(eh.visitors, client)
}

func (eh *EventHandler) generateOutboundTableUse(table int) {
	queuedClient := eh.queuedVisitors.Remove(eh.queuedVisitors.Front()).(visitor)
	eh.visitors[queuedClient] = table
}
