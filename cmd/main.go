package main

import (
	"fmt"
	"github.com/Paincake/yadro/event"
	"io"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: <task.exe> <filename>")
		return
	}
	filename := os.Args[1]
	file, err := os.Open(filename)
	src, err := event.NewClubFileSource(file)
	if err != nil {
		panic(err)
	}
	processor := event.NewProcessor(src, os.Stdout, src.Club)
	err = processor.ProcessEvents()
	if err != nil {
		io.WriteString(os.Stderr, err.Error())
	}

}
