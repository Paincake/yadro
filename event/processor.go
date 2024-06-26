package event

import (
	"github.com/Paincake/yadro/event/club"
	"io"
	"strconv"
	"strings"
	"time"
)

type Processor struct {
	Club     *club.Club
	eventSrc Source
	eventDst io.Writer
}

func NewProcessor(eventSrc Source, eventDst io.Writer, club *club.Club) *Processor {
	p := &Processor{
		Club:     club,
		eventSrc: eventSrc,
		eventDst: eventDst,
	}
	return p
}

func (p *Processor) endProcessing() {
	event := ClubClosingEvent{
		closeTime:    p.Club.CloseTime,
		ResultWriter: p.eventDst,
	}
	event.ProcessEvent(p.Club)
}

func (p *Processor) startProcessing() {
	event := ClubOpeningEvent{
		openTime:     p.Club.OpenTime,
		ResultWriter: p.eventDst,
	}
	event.ProcessEvent(p.Club)
}

func (p *Processor) ProcessEvents() error {
	p.startProcessing()
	for {
		processed, err := p.processEvent()
		if err != nil {
			return err
		}
		if !processed {
			break
		}
	}
	return nil
}

func (p *Processor) processEvent() (bool, error) {
	eventData, err := p.eventSrc.GetEventData()
	if err != nil {
		return false, err
	}

	if eventData == "" {
		p.endProcessing()
		return false, nil
	}

	eventParts := strings.Split(eventData, " ")

	eventTime, err := time.Parse(club.TimeFormat_hhmm, eventParts[0])
	if err != nil {
		return false, err
	}

	eventID := eventParts[1]
	clientID := eventParts[2]
	var sideParam any
	if len(eventParts) == 4 {
		if eventID == club.CLIENT_TABLE_USING_IN {
			sideParam, err = strconv.Atoi(eventParts[3])
			if err != nil {
				return false, err
			}

		}
	}
	var event ClubEvent
	switch eventID {
	case club.CLIENT_ARRIVING_IN:
		event = &ClientArrivalEvent{
			ClientID:     clientID,
			ArrTime:      eventTime,
			ResultWriter: p.eventDst,
		}
	case club.CLIENT_TABLE_USING_IN:
		event = &ClientTablePickEvent{
			ClientID:     clientID,
			Table:        sideParam.(int),
			EventTime:    eventTime,
			ResultWriter: p.eventDst,
		}
	case club.CLIENT_WAITING_IN:
		event = &ClientWaitingEvent{
			ClientID:     clientID,
			EventTime:    eventTime,
			ResultWriter: p.eventDst,
		}
	case club.CLIENT_LEAVING_IN:
		event = &ClientLeavingEvent{
			ClientID:     clientID,
			LeavingTime:  eventTime,
			ResultWriter: p.eventDst,
		}
	}
	event.ProcessEvent(p.Club)
	return true, nil
}
