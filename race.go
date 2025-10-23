package main

import (
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

func startRace(numRacers, trackLength int, updates chan<- RaceUpdate,
	results chan<- RaceResult, wg *sync.WaitGroup) {

	start := time.Now()

	go func() {
		wg.Wait()
		// closes channels after all racers finish
		close(updates)
		close(results)
	}()

	for id := 0; id < numRacers; id++ {
		wg.Add(1)
		go func(racerID int) {
			defer wg.Done()

			// pushes racer, when finished, to give RaceResult
			pos := 0
			base := 60 + rand.Intn(90)

			for pos < trackLength {
				time.Sleep(time.Duration(base+rand.Intn(60)) * time.Millisecond)
				step := 1 + rand.Intn(2)
				pos += step
				if pos > trackLength {
					pos = trackLength
				}

				percent := int(float64(pos) * 100.0 / float64(trackLength))
				if percent > 100 {
					percent = 100
				}

				updates <- RaceUpdate{
					id:       racerID,
					position: percent,
					elapsed:  time.Since(start),
					finished: pos >= trackLength,
				}

				if pos >= trackLength {

					results <- RaceResult{
						id:         racerID,
						finishTime: time.Since(start),
					}
					return
				}
			}
		}(id)
	}
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
