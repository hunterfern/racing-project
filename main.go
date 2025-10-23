package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	_ "github.com/gorilla/websocket"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	const numRacers = 8
	const trackLength = 120

	uiNumRacers = numRacers

	hub := NewHub()
	go hub.Run()
	commands := make(chan Command, 32)

	updates := make(chan RaceUpdate, 64)
	results := make(chan RaceResult, numRacers)

	go displayRaceUIWithWS(updates, hub, commands)

	for {
		cmd := <-commands
		if cmd.Kind == "START" {
			var wg sync.WaitGroup
			startTime := time.Now()

			startRace(numRacers, trackLength, updates, results, &wg)

			allResults := make([]RaceResult, 0, numRacers)
			for i := 0; i < numRacers; i++ {
				r := <-results
				allResults = append(allResults, r)
			}

			endTime := time.Since(startTime)
			winner := determineWinner(allResults)
			if winner.id >= 0 {
				fmt.Printf("\n Winner: Racer %d\n", winner.id+1)
			}
			fmt.Println(" Total time:", endTime.Seconds())
			logRaceResults(allResults, winner, endTime)

			break
		}
	}

	fmt.Printf("\n Race results have been logged to race_log.txt\n")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
