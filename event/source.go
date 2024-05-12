package event

import (
	"bufio"
	"errors"
	"github.com/Paincake/yadro/event/club"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type Source interface {
	initSource() error
	GetEventData() (string, error)
}

type ClubFileSource struct {
	Club   *club.Club
	reader *bufio.Reader
}

func NewClubFileSource(file *os.File) (*ClubFileSource, error) {
	reader := bufio.NewReader(file)
	fs := ClubFileSource{
		reader: reader,
	}
	err := fs.initSource()
	return &fs, err
}

func (s *ClubFileSource) initSource() error {
	tableData, err := s.reader.ReadString('\n')
	if err != nil {
		return err
	}
	tableData, _, _ = strings.Cut(tableData, "\r\n")
	tables, err := strconv.Atoi(tableData)
	if err != nil {
		return err
	}

	openClosedTime, err := s.reader.ReadString('\n')
	if err != nil {
		return err
	}
	openClosedTime, _, _ = strings.Cut(openClosedTime, "\r\n")
	times := strings.Split(openClosedTime, " ")
	openTime, err := time.Parse(club.TimeFormat_hhmm, times[0])
	if err != nil {
		return err
	}
	closeTime, err := time.Parse(club.TimeFormat_hhmm, times[1])
	if err != nil {
		return err
	}

	hourCostData, err := s.reader.ReadString('\n')
	if err != nil {
		return err
	}
	hourCostData, _, _ = strings.Cut(hourCostData, "\r\n")
	hourCost, err := strconv.Atoi(hourCostData)
	if err != nil {
		return err
	}

	s.Club = club.NewClub(tables, hourCost, openTime, closeTime)
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
	event, _, _ = strings.Cut(event, "\r\n")
	return event, nil
}
