package event

import (
	"container/list"
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
	TableMap       map[int]table
}

func NewClub(tables int, hourCost int, openTime time.Time, closeTime time.Time) *Club {
	tableMap := make(map[int]table)
	for i := 0; i < tables; i++ {
		tableMap[i+1] = table{
			isInUse: false,
		}
	}
	return &Club{
		tablesFree:     tables,
		TableMap:       tableMap,
		HourCost:       hourCost,
		OpenTime:       openTime,
		CloseTime:      closeTime,
		clients:        make(map[string]client),
		queuedVisitors: list.New(),
	}
}

func (c *Club) GetRemainingClientsSorted() []string {
	clientIDS := make([]string, 0)
	for k, v := range c.clients {
		if !v.left {
			clientIDS = append(clientIDS, k)
		}
	}
	sort.Strings(clientIDS)
	return clientIDS
}

func (c *Club) AddClient(clientID string, arrTime time.Time) {
	c.clients[clientID] = client{clientID: clientID, arrTime: arrTime, left: false}
}

func (c *Club) ClientExists(clientID string) bool {
	client, ok := c.clients[clientID]
	if !ok {
		return false
	}
	return !client.left
}

func (c *Club) PickTable(clientID string, tableNum int, pickTime time.Time) int {
	tableToPick := c.TableMap[tableNum]
	if tableToPick.isInUse {
		return 0
	}
	client := c.clients[clientID]
	if client.table == 0 {
		c.tablesFree--
	}
	client.table = tableNum

	tableToPick.isInUse = true
	tableToPick.pickTime = pickTime

	c.TableMap[tableNum] = tableToPick
	c.clients[clientID] = client
	return client.table
}

type table struct {
	isInUse      bool
	pickTime     time.Time
	minutesInUse int
}

func (c *Club) IsBusy() bool {
	return c.tablesFree == 0
}

func (c *Club) EnqueueClient(clientID string) error {
	if c.queuedVisitors.Len() >= len(c.TableMap) {
		return QueueOverflowError{}
	}
	c.queuedVisitors.PushBack(clientID)
	return nil
}

func (c *Club) DequeueClient(tableNum int, eventTime time.Time) string {
	clientID := ""
	if c.queuedVisitors.Len() > 0 {
		clientID = c.queuedVisitors.Remove(c.queuedVisitors.Front()).(string)
		client := c.clients[clientID]
		client.table = tableNum

		pickedTable := c.TableMap[tableNum]
		pickedTable.isInUse = true
		pickedTable.pickTime = eventTime
		c.TableMap[tableNum] = pickedTable

		c.tablesFree--
		c.clients[clientID] = client
	}
	return clientID
}

func (c *Club) RemoveClient(clientID string, leavingTime time.Time) int {
	client := c.clients[clientID]
	client.leaveTime = leavingTime

	pickedTable, ok := c.TableMap[client.table]
	if ok {
		pickedTable.isInUse = false
		pickedTable.minutesInUse += int(leavingTime.Sub(pickedTable.pickTime).Minutes())
		c.TableMap[client.table] = pickedTable
		c.tablesFree++
	}

	var queueElem *list.Element
	for e := c.queuedVisitors.Front(); e != nil; e = e.Next() {
		v := e.Value.(string)
		if v == clientID {
			queueElem = e
			break
		}
	}
	if queueElem != nil {
		c.queuedVisitors.Remove(queueElem)
	}

	table := client.table
	client.left = true
	c.clients[clientID] = client
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

func (c *Club) CountProfit() []TableProfit {
	profitInfo := make([]TableProfit, 0, len(c.TableMap))
	for i := 1; i <= len(c.TableMap); i++ {
		table := c.TableMap[i]
		var profit int
		if table.minutesInUse != 0 {
			profit = c.HourCost * (table.minutesInUse/60 + 1)
		}
		profitInfo = append(profitInfo, TableProfit{
			TableNum:     i,
			Profit:       profit,
			MinutesInUse: table.minutesInUse,
		})
	}
	return profitInfo
}

type TableProfit struct {
	TableNum     int
	Profit       int
	MinutesInUse int
}

type client struct {
	clientID  string
	arrTime   time.Time
	leaveTime time.Time
	table     int
	left      bool
}
