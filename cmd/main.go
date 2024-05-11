package main

import (
	"github.com/Paincake/yadro/event"
	"os"
)

func main() {
	f, _ := os.Open("C:\\Users\\ryazh\\GolandProjects\\yadro\\test")
	src := event.NewClubFileSource(f)
	err := src.InitSource()
	if err != nil {
		panic(err)
	}

	processor := event.NewProcessor(src, os.Stdout)

	//TODO
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
