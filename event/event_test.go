package event

import (
	"fmt"
	"testing"
	"time"
)

func TestEventHandler_HandleInboundEvent(t *testing.T) {
	time1, err := time.Parse("15:04", "12:10")
	if err != nil {
		fmt.Println(err)
	}
	time2, err := time.Parse("15:04", "12:09")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(time2.Sub(time1).Hours())

}
