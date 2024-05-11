package event

import (
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type Processor struct {
	Club     *Club
	eventSrc Source
	eventDst io.Writer
}

func NewProcessor(eventSrc Source, eventDst io.Writer) *Processor {
	p := &Processor{
		eventSrc: eventSrc,
		eventDst: eventDst,
	}
	return p
}

func (p *Processor) EndProcessing() {
	event := ClubClosingEvent{
		ResultWriter: p.eventDst,
	}
	event.ProcessEvent(p.Club)
}

func (p *Processor) ProcessEvent() (bool, error) {
	eventData, err := p.eventSrc.GetEventData()
	if err != nil {
		return false, err
	}

	if eventData == "" {
		p.EndProcessing()
		return false, nil
	}

	eventParts := strings.Split(eventData, " ")

	eventTime, err := time.Parse(hhmm, eventParts[0])
	if err != nil {
		return false, err
	}

	eventID := eventParts[1]
	clientID := eventParts[2]
	var sideParam any
	if len(eventParts) == 4 {
		if eventID == CLIENT_TABLE_USING_IN {
			sideParam, err = strconv.Atoi(eventParts[3])
			if err != nil {
				return false, err
			}

		}
	}
	var event ClubEvent
	switch eventID {
	case CLIENT_ARRIVING_IN:
		event = &ClientArrivalEvent{
			ClientID:     clientID,
			ArrTime:      eventTime,
			ResultWriter: p.eventDst,
		}
	case CLIENT_TABLE_USING_IN:
		event = &ClientTablePickEvent{
			ClientID:     clientID,
			Table:        sideParam.(int),
			EventTime:    eventTime,
			ResultWriter: p.eventDst,
		}
	case CLIENT_WAITING_IN:
		event = &ClientWaitingEvent{
			ClientID:     clientID,
			EventTime:    eventTime,
			ResultWriter: p.eventDst,
		}
	case CLIENT_LEAVING_IN:
		event = &ClientLeavingEvent{
			ClientID:     clientID,
			LeavingTime:  eventTime,
			ResultWriter: p.eventDst,
		}
	}
	event.ProcessEvent(p.Club)
	io.WriteString(os.Stdout, "\n")
	return true, nil
}
