package main

import (
	"fmt"
	"os"
	"sort"
	"time"
)


func logRaceResults(all []RaceResult, winner RaceResult, total time.Duration) {
	f, err := os.OpenFile("race_log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
    		return
	}
	defer f.Close()

	cp := make([]RaceResult, len(all))
		copy(cp, all)
		sort.Slice(cp, func(i, j int) bool {
			if cp[i].finishTime == cp[j].finishTime {
				return cp[i].id < cp[j].id
			}
			return cp[i].finishTime < cp[j].finishTime
		})

		ts := time.Now().Format(time.RFC3339)
		fmt.Fprintf(f, "=== Race @ %s ===\n", ts)
		fmt.Fprintf(f, "Total racers: %d\n", len(all))
		fmt.Fprintf(f, "Winner: Racer %d (%.2fs)\n", winner.id, winner.finishTime.Seconds())
		fmt.Fprintf(f, "Total wall time: %.2fs\n", total.Seconds())
		fmt.Fprintln(f, "Placing:")
		for place, r := range cp {
			fmt.Fprintf(f, " %2d) Racer %d %.2fs\n", place + 1, r.id, r.finishTime.Seconds())
		}
		fmt.Fprintln(f)
}
