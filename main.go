package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	_ "github.com/gorilla/websocket"
)

// default values if client doesn't provide racers
const defaultNumRacers = 8
const trackLength = 120

func main() {
	rand.Seed(time.Now().UnixNano())

	// start with default; will be updated on START command if client asks
	uiNumRacers = defaultNumRacers

	hub := NewHub()
	go hub.Run()

	commands := make(chan Command, 32)
	updates := make(chan RaceUpdate, 256) // keep channels long-lived and shared
	results := make(chan RaceResult, 256)

	// run UI display (reads from updates channel)
	go displayRaceUIWithWS(updates, hub, commands)

	// main command loop: react to START commands from any connected client
	for cmd := range commands {
		if cmd.Kind == "START" {
			// parse optional racers arg (string) into int
			numRacers := defaultNumRacers
			if v, ok := cmd.Args["racers"]; ok && v != "" {
				if n, err := strconv.Atoi(v); err == nil {
					// clamp to allowed range 4..12
					if n < 4 {
						n = 4
					}
					if n > 12 {
						n = 12
					}
					numRacers = n
				}
			}

			// update global uiNumRacers *before* starting the race so UI and server agree
			uiNumRacers = numRacers
			fmt.Printf("Starting race with %d racers\n", numRacers)

			// create a fresh wait group for this race
			var wg sync.WaitGroup

			lines, stats := GetRaceLines(numRacers)
			setCurrentStats(stats)

			hub.broadcast <- mustJSON(WSMessage{Type: "lines", Data: lines})

			// start the race (race.go's startRace uses passed numRacers)
			startRace(numRacers, trackLength, updates, results, &wg)

			// collect results for exactly numRacers racers
			allResults := make([]RaceResult, 0, numRacers)
			for i := 0; i < numRacers; i++ {
				r := <-results
				allResults = append(allResults, r)
			}

			// winner calculation and logging
			endTime := time.Since(time.Now()) // not critical; keep for backward compatibility
			winner := determineWinner(allResults)
			if winner.id >= 0 {
				fmt.Printf("\n Winner: Racer %d\n", winner.id+1)
			}
			fmt.Println(" Total results collected:", len(allResults))
			logRaceResults(allResults, winner, endTime)

			fmt.Printf("\n Race results logged. Ready for another race!\n")
			// loop continues to accept more START commands
		}
	}

	// keep the binary alive in case commands channel closed (not expected)
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
