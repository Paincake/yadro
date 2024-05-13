package event

import (
	"bytes"
	"github.com/Paincake/yadro/event/club"
	"testing"
	"time"
)

func TestClientArrivalEvent_Valid(t *testing.T) {
	bBuf := make([]byte, 0)
	buf := bytes.NewBuffer(bBuf)

	baseTime, _ := time.Parse(club.TimeFormat_hhmm, "09:00")
	c := club.NewClub(1, 10, baseTime.Add(time.Hour), baseTime.Add(time.Hour*3))
	event := ClientArrivalEvent{
		ClientID:     "client1",
		ArrTime:      baseTime,
		ResultWriter: buf,
	}
	event.ProcessEvent(c)
	event = ClientArrivalEvent{
		ClientID:     "client1",
		ArrTime:      baseTime.Add(time.Hour),
		ResultWriter: buf,
	}
	event.ProcessEvent(c)
	event = ClientArrivalEvent{
		ClientID:     "client1",
		ArrTime:      baseTime.Add(time.Hour * 2),
		ResultWriter: buf,
	}
	event.ProcessEvent(c)

	expected := "09:00 1 client1\n" +
		"09:00 13 NotOpenYet\n" +
		"10:00 1 client1\n" +
		"11:00 1 client1\n" +
		"11:00 13 YouShallNotPass\n"
	actual := buf.String()
	if actual != expected {
		t.Fatalf("\nExpected:\n%s\nGot:\n%s", expected, actual)
	}
}

func TestClientTablePickEvent_ProcessEvent(t *testing.T) {
	bBuf := make([]byte, 0)
	buf := bytes.NewBuffer(bBuf)

	baseTime, _ := time.Parse(club.TimeFormat_hhmm, "09:00")
	c := club.NewClub(1, 10, baseTime.Add(time.Hour), baseTime.Add(time.Hour*3))
	arrEvent1 := ClientArrivalEvent{
		ClientID:     "client1",
		ArrTime:      baseTime.Add(time.Hour),
		ResultWriter: buf,
	}
	arrEvent1.ProcessEvent(c)
	arrEvent2 := ClientArrivalEvent{
		ClientID:     "client2",
		ArrTime:      baseTime.Add(time.Hour),
		ResultWriter: buf,
	}
	arrEvent2.ProcessEvent(c)

	pickEvent1 := ClientTablePickEvent{
		ClientID:     "client1",
		Table:        1,
		EventTime:    baseTime.Add(time.Hour * 2),
		ResultWriter: buf,
	}
	pickEvent1.ProcessEvent(c)
	pickEvent2 := ClientTablePickEvent{
		ClientID:     "client2",
		Table:        1,
		EventTime:    baseTime.Add(time.Hour * 2),
		ResultWriter: buf,
	}
	pickEvent2.ProcessEvent(c)
	expected := "10:00 1 client1\n" +
		"10:00 1 client2\n" +
		"11:00 2 client1 1\n" +
		"11:00 2 client2 1\n" +
		"11:00 13 PlaceIsBusy\n"
	if buf.String() != expected {
		t.Fatalf("\nExpected:\n%s\nGot:\n%s", expected, buf.String())
	}
}

func TestClientWaitingEvent_ProcessEvent(t *testing.T) {
	bBuf := make([]byte, 0)
	buf := bytes.NewBuffer(bBuf)

	baseTime, _ := time.Parse(club.TimeFormat_hhmm, "09:00")
	c := club.NewClub(1, 10, baseTime.Add(time.Hour), baseTime.Add(time.Hour*4))

	var event ClubEvent
	event = &ClientArrivalEvent{
		ClientID:     "client1",
		ArrTime:      baseTime.Add(time.Hour),
		ResultWriter: buf,
	}
	event.ProcessEvent(c)

	event = &ClientArrivalEvent{
		ClientID:     "client2",
		ArrTime:      baseTime.Add(time.Hour * 2),
		ResultWriter: buf,
	}
	event.ProcessEvent(c)

	event = &ClientWaitingEvent{
		ClientID:     "client2",
		EventTime:    baseTime.Add(time.Hour * 2),
		ResultWriter: buf,
	}
	event.ProcessEvent(c)

	pickEvent1 := &ClientTablePickEvent{
		ClientID:     "client1",
		Table:        1,
		EventTime:    baseTime.Add(time.Hour * 2),
		ResultWriter: buf,
	}
	pickEvent1.ProcessEvent(c)

	event = &ClientWaitingEvent{
		ClientID:     "client2",
		EventTime:    baseTime.Add(time.Hour * 2),
		ResultWriter: buf,
	}
	event.ProcessEvent(c)

	event = &ClientArrivalEvent{
		ClientID:     "client3",
		ArrTime:      baseTime.Add(time.Hour * 3),
		ResultWriter: buf,
	}
	event.ProcessEvent(c)

	event = &ClientWaitingEvent{
		ClientID:     "client3",
		EventTime:    baseTime.Add(time.Hour * 3),
		ResultWriter: buf,
	}
	event.ProcessEvent(c)

	expected := "10:00 1 client1\n" +
		"11:00 1 client2\n" +
		"11:00 3 client2\n" +
		"11:00 13 ICanWaitNoLonger\n" +
		"11:00 2 client1 1\n" +
		"11:00 3 client2\n" +
		"12:00 1 client3\n" +
		"12:00 3 client3\n" +
		"12:00 11 client3\n"
	actual := buf.String()

	if buf.String() != expected {
		t.Fatalf("\nExpected:\n%s\nGot:\n%s", expected, actual)
	}
}

func TestClientLeavingEvent_ProcessEvent(t *testing.T) {
	bBuf := make([]byte, 0)
	buf := bytes.NewBuffer(bBuf)

	baseTime, _ := time.Parse(club.TimeFormat_hhmm, "09:00")
	c := club.NewClub(1, 10, baseTime.Add(time.Hour), baseTime.Add(time.Hour*4))

	var event ClubEvent
	event = &ClientArrivalEvent{
		ClientID:     "client1",
		ArrTime:      baseTime.Add(time.Hour),
		ResultWriter: buf,
	}
	event.ProcessEvent(c)

	event = &ClientArrivalEvent{
		ClientID:     "client2",
		ArrTime:      baseTime.Add(time.Hour * 2),
		ResultWriter: buf,
	}
	event.ProcessEvent(c)

	pickEvent1 := &ClientTablePickEvent{
		ClientID:     "client1",
		Table:        1,
		EventTime:    baseTime.Add(time.Hour * 2),
		ResultWriter: buf,
	}
	pickEvent1.ProcessEvent(c)

	event = &ClientWaitingEvent{
		ClientID:     "client2",
		EventTime:    baseTime.Add(time.Hour * 2),
		ResultWriter: buf,
	}
	event.ProcessEvent(c)

	event = &ClientLeavingEvent{
		ClientID:     "client1",
		LeavingTime:  baseTime.Add(time.Hour * 2),
		ResultWriter: buf,
	}
	event.ProcessEvent(c)

	event = &ClientArrivalEvent{
		ClientID:     "client3",
		ArrTime:      baseTime.Add(time.Hour * 3),
		ResultWriter: buf,
	}
	event.ProcessEvent(c)

	event = &ClientWaitingEvent{
		ClientID:     "client3",
		EventTime:    baseTime.Add(time.Hour * 3),
		ResultWriter: buf,
	}
	event.ProcessEvent(c)

	event = &ClientLeavingEvent{
		ClientID:     "client3",
		LeavingTime:  baseTime.Add(time.Hour * 3),
		ResultWriter: buf,
	}
	event.ProcessEvent(c)

	event = &ClientLeavingEvent{
		ClientID:     "client2",
		LeavingTime:  baseTime.Add(time.Hour * 3),
		ResultWriter: buf,
	}
	event.ProcessEvent(c)

	expected := "10:00 1 client1\n" +
		"11:00 1 client2\n" +
		"11:00 2 client1 1\n" +
		"11:00 3 client2\n" +
		"11:00 4 client1\n" +
		"11:00 12 client2 1\n" +
		"12:00 1 client3\n" +
		"12:00 3 client3\n" +
		"12:00 4 client3\n" +
		"12:00 4 client2\n"

	actual := buf.String()
	if buf.String() != expected {
		t.Fatalf("\nExpected:\n%s\nGot:\n%s", expected, actual)
	}
}
