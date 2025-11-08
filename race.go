package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type RaceUpdate struct {
	id       int
	position int
	elapsed  time.Duration
	finished bool
}

type RaceResult struct {
	id         int
	finishTime time.Duration
}

var raceMu sync.Mutex
var raceActive bool
var stopRace chan struct{} // signal to stop current race safely

func startRace(numRacers, trackLength int, updates chan<- RaceUpdate,
	results chan<- RaceResult, wg *sync.WaitGroup) {

	raceMu.Lock()
	if raceActive {
		fmt.Println("Race already active — skipping new start.")
		raceMu.Unlock()
		return
	}
	raceActive = true
	stopRace = make(chan struct{})
	raceMu.Unlock()

	for id := 0; id < numRacers; id++ {
		wg.Add(1)
		go func(racerID int) {
			defer wg.Done()

			baseStep := 1.00 + rand.Float64()*0.40
			stepVar := 0.25
			fadeAfter := 0.60
			fadeAmt := 0.35
			burstChance := 0.06
			burstBoost := 0.80
			burstCDTicks := 18

			tickMin := 35 * time.Millisecond
			tickMax := 70 * time.Millisecond

			pos := 0.0
			ticksSinceBurst := burstCDTicks
			start := time.Now()

			for {
				select {
				case <-stopRace:
					return
				default:
				}

				jitter := time.Duration(rand.Intn(int(tickMax-tickMin+1))) + tickMin
				time.Sleep(jitter)

				prog := pos / float64(trackLength)
				fade := 1.0
				if prog > fadeAfter {
					pct := (prog - fadeAfter) / (1.0 - fadeAfter)
					fade = 1.0 - fadeAmt*pct
				}

				noise := (rand.Float64()*2 - 1) * stepVar

				burst := 0.0
				if ticksSinceBurst >= burstCDTicks && rand.Float64() < burstChance {
					burst = burstBoost * (0.5 + rand.Float64()*0.5)
					ticksSinceBurst = 0
				} else {
					ticksSinceBurst++
				}

				packNudge := (rand.Float64()*2 - 1) * 0.10

				step := (baseStep + noise + burst + packNudge) * fade
				if step < 0.2 {
					step = 0.2
				}

				pos += step
				if pos > float64(trackLength) {
					pos = float64(trackLength)
				}

				percent := int((pos / float64(trackLength)) * 100.0)
				if percent > 100 {
					percent = 100
				} else if percent < 0 {
					percent = 0
				}

				updates <- RaceUpdate{
					id:       racerID,
					position: percent,
					elapsed:  time.Since(start),
					finished: pos >= float64(trackLength),
				}

				if pos >= float64(trackLength) {
					results <- RaceResult{
						id:         racerID,
						finishTime: time.Since(start),
					}
					return
				}
			}

		}(id)
	}

	// Monitor completion
	go func() {
		wg.Wait()
		raceMu.Lock()
		defer raceMu.Unlock()
		raceActive = false
		close(stopRace) // signal any late goroutines to stop
		fmt.Println("Race completed — ready for another race!")
	}()
}

func resetRace() {
	raceMu.Lock()
	defer raceMu.Unlock()
	if !raceActive {
		return
	}
	fmt.Println("Resetting race...")
	select {
	case <-stopRace:
	default:
		close(stopRace) // stop current race
	}
	raceActive = false
}

func determineWinner(all []RaceResult) RaceResult {
	if len(all) == 0 {
		return RaceResult{id: -1}
	}
	winner := all[0]
	for i := 1; i < len(all); i++ {
		r := all[i]
		if r.finishTime < winner.finishTime ||
			(r.finishTime == winner.finishTime && r.id < winner.id) {
			winner = r
		}
	}
	return winner
}
