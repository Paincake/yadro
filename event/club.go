package event

import (
	"container/list"
	"fmt"
	"sort"
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

type Club struct {
	HourCost       int
	tablesFree     int
	OpenTime       time.Time
	CloseTime      time.Time
	queuedVisitors *list.List
	clients        map[string]client
	tableMap       map[int]bool
}

func NewClub(tables int, hourCost int, openTime time.Time, closeTime time.Time) *Club {
	tableMap := make(map[int]bool)
	for i := 0; i < tables; i++ {
		tableMap[i+1] = false
	}
	return &Club{
		tableMap:       make(map[int]bool),
		HourCost:       hourCost,
		OpenTime:       openTime,
		CloseTime:      closeTime,
		clients:        make(map[string]client),
		queuedVisitors: list.New(),
	}
}

func (c *Club) AddClient(clientID string, arrTime time.Time) {
	c.clients[clientID] = client{clientID: clientID, arrTime: arrTime}
}

func (c *Club) ClientExists(clientID string) bool {
	_, ok := c.clients[clientID]
	return ok
}

func (c *Club) PickTable(clientID string, table int) int {
	if c.tableMap[table] {
		return 0
	}
	client := c.clients[clientID]
	if client.table == 0 {
		c.tablesFree--
	}
	client.table = table
	return client.table
}

func (c *Club) IsBusy() bool {
	return c.tablesFree == 0
}

func (c *Club) EnqueueClient(clientID string) error {
	if c.queuedVisitors.Len() >= len(c.tableMap) {
		return QueueOverflowError{}
	}
	c.queuedVisitors.PushBack(clientID)
	return nil
}

func (c *Club) DequeueClient(table int) string {
	clientID := c.queuedVisitors.Remove(c.queuedVisitors.Front()).(string)
	client := c.clients[clientID]
	client.table = table
	c.tableMap[client.table] = true
	c.tablesFree--
	return clientID
}

func (c *Club) RemoveClient(clientID string, leavingTime time.Time) int {
	client := c.clients[clientID]
	client.leaveTime = leavingTime
	c.tableMap[client.table] = false
	table := client.table
	client.table = 0
	c.tablesFree++
	return table
}

func (c *Club) GetClientsSorted() []string {
	clients := make([]string, 0, len(c.clients))
	for k, _ := range c.clients {
		clients = append(clients, k)
	}
	sort.Strings(clients)
	return clients
}

func (c *Club) CountProfit() []ClientProfit {
	clientIDS := c.GetClientsSorted()
	profits := make([]ClientProfit, 0, len(clientIDS))
	for _, id := range clientIDS {
		client := c.clients[id]
		spentTime := client.leaveTime.Sub(client.arrTime)
		hours := int(spentTime.Hours())
		minutes := int(spentTime.Minutes())
		profit := c.HourCost * hours
		profits = append(profits, ClientProfit{
			ClientID: id,
			Time:     fmt.Sprintf("%d:%d", hours, minutes),
			Profit:   profit,
		})
	}
	return profits
}

type ClientProfit struct {
	ClientID string
	Time     string
	Profit   int
}

type client struct {
	clientID  string
	arrTime   time.Time
	leaveTime time.Time
	table     int
}
