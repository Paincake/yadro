package main

import (
	"fmt"
	"github.com/Paincake/yadro/event"
	"os"
)

//TODO testing
//TODO restructure
//TODO start and end processing printing
//TODO better processing loop

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: <task.exe> <filename>")
		return
	}
	filename := os.Args[1]
	f, _ := os.Open(filename)
	src, err := event.NewClubFileSource(f)
	if err != nil {
		panic(err)
	}

	processor := event.NewProcessor(src, os.Stdout)
	//TODO this is so bad
	processor.Club = src.Club

	for {
		processed, err := processor.ProcessEvent()
		if err != nil {
			panic(err)
		}
		if !processed {
			break
		}
	}
}
