package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	quake_log_parser "github.com/camopy/quake-log-parser"
)

func main() {
	outputFilePtr := flag.String("output-file", "", "output file path")

	flag.Parse()

	inputLogPath := flag.Arg(0)
	if inputLogPath == "" {
		log.Fatalln("input log path is required")
		return
	}

	p := quake_log_parser.NewParser(inputLogPath)
	games, err := p.Parse()
	if err != nil {
		log.Fatalln(err)
	}

	output, err := json.MarshalIndent(games.GenerateReport(), "", "  ")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(output))

	if *outputFilePtr != "" {
		if err := os.WriteFile(*outputFilePtr, output, 0644); err != nil {
			log.Fatalln(err)
		}
	}
}
