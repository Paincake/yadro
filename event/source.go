package event

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
	"time"
)

type Source interface {
	InitSource() error
	GetEventData() (string, error)
}

type ClubFileSource struct {
	Club   *Club
	reader bufio.Reader
}

func (s *ClubFileSource) InitSource() error {
	tableData, err := s.reader.ReadString('\n')
	if err != nil {
		return err
	}
	tableData, _, _ = strings.Cut(tableData, "\n")
	tables, err := strconv.Atoi(tableData)
	if err != nil {
		return err
	}

	openClosedTime, err := s.reader.ReadString('\n')
	if err != nil {
		return err
	}
	times := strings.Split(openClosedTime, " ")
	openTime, err := time.Parse(hhmm, times[0])
	if err != nil {
		return err
	}
	closeTime, err := time.Parse(hhmm, times[0])
	if err != nil {
		return err
	}

	hourCostData, err := s.reader.ReadString('\n')
	if err != nil {
		return err
	}
	hourCost, err := strconv.Atoi(hourCostData)
	if err != nil {
		return err
	}

	s.Club = NewClub(tables, hourCost, openTime, closeTime)
	return nil
}

func (s *ClubFileSource) GetEventData() (string, error) {
	event, err := s.reader.ReadString('\n')
	if err != nil {
		if errors.Is(err, io.EOF) {
			return "", nil
		} else {
			return "", err
		}
	}
	event, _, _ = strings.Cut(event, "\n")
	return event, nil
}
