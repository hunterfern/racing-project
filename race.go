package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

type RacerStats struct {
	Speed   float64
	Stamina float64
	Burst   float64
	Gate    float64
}

type RacerLine struct {
	ID      int     `json:"id"`
	Prob    float64 `json:"prob"`
	Amer    int     `json:"amer"`
	FracNum int     `json:"fracN"`
	FracDen int     `json:"fracD"`
	Implied int     `json:"impPct"`
}

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
var currentStats []RacerStats
var statsMu sync.RWMutex

func setCurrentStats(ss []RacerStats) {
	statsMu.Lock()
	defer statsMu.Unlock()
	currentStats = make([]RacerStats, len(ss))
	copy(currentStats, ss)
}

func getStat(i int) RacerStats {
	statsMu.RLock()
	defer statsMu.RUnlock()
	if i >= 0 && i < len(currentStats) {
		return currentStats[i]
	}
	return RacerStats{
		Speed:   0.55,
		Stamina: 0.55,
		Burst:   0.5,
		Gate:    0.5,
	}
}

const (
	wSpeed   = 0.50
	wStamina = 0.30
	wBurst   = 0.15
	wGate    = 0.05
)

func clamp(x, lo, hi float64) float64 {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}

func softmax(xs []float64) []float64 {
	if len(xs) == 0 {
		return nil
	}

	max := xs[0]
	for _, v := range xs[1:] {
		if v > max {
			max = v
		}
	}
	sum := 0.0
	out := make([]float64, len(xs))
	for i, v := range xs {
		e := math.Exp(v - max)
		out[i] = e
		sum += e
	}
	for i := range out {
		out[i] /= sum
	}
	return out
}

func genRacerStats(num int) []RacerStats {
	stats := make([]RacerStats, num)
	for i := 0; i < num; i++ {
		stats[i] = RacerStats{
			Speed:   clamp(rand.NormFloat64()*0.12+0.55, 0.25, 0.95),
			Stamina: clamp(rand.NormFloat64()*0.12+0.55, 0.25, 0.95),
			Burst:   clamp(rand.NormFloat64()*0.12+0.50, 0.20, 0.95),
			Gate:    clamp(rand.NormFloat64()*0.10+0.50, 0.20, 0.95),
		}
	}
	return stats
}

func statsToProbs(ss []RacerStats) []float64 {
	if len(ss) == 0 {
		return nil
	}
	scores := make([]float64, len(ss))
	for i, s := range ss {
		score := wSpeed*s.Speed + wStamina*s.Stamina + wBurst*s.Burst + wGate*s.Gate
		score += rand.NormFloat64() * 0.03
		scores[i] = score
	}
	ps := softmax(scores)

	const eps = 1e-6
	sum := 0.0
	for i, p := range ps {
		p = clamp(p, eps, 1.0)
		ps[i] = p
		sum += p
	}

	for i := range ps {
		ps[i] /= sum
	}
	return ps
}

func buildLines(probs []float64) []RacerLine {
	lines := make([]RacerLine, len(probs))
	for i, p := range probs {
		amer := probToAmerican(p)
		fn, fd := probToFraction(p)
		lines[i] = RacerLine{
			ID:      i,
			Prob:    p,
			Amer:    amer,
			FracNum: fn,
			FracDen: fd,
			Implied: int(math.Round(p * 100.0)),
		}
	}
	return lines
}

func GetRaceLines(numRacers int) ([]RacerLine, []RacerStats) {
	ss := genRacerStats(numRacers)
	ps := statsToProbs(ss)
	return buildLines(ps), ss
}

func probToAmerican(p float64) int {
	if p <= 0 {
		return 0
	}
	if p >= 1 {
		return -100000
	}
	if p >= 0.5 {
		return int(-100 * p / (1 - p))
	}
	return int((1 - p) * 100 / p)
}

func probToFraction(p float64) (n, d int) {
	if p <= 0 {
		return 0, 1
	}
	imp := (1.0 / p) - 1.0
	half := math.Round(imp*2.0) / 2.0
	n = int(math.Round(half * 2.0))
	d = 2
	g := gcd(n, d)
	return n / g, d / g
}
func gcd(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	if a == 0 {
		return 1
	}
	return a
}

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

			s := getStat(racerID)

			// map stats to movement parameters (tuned for your existing visuals)
			baseStep := 0.80 + s.Speed*0.60        // average forward step size
			stepVar := 0.10 + (1.0-s.Stamina)*0.20 // low stamina -> more wobble
			fadeAfter := 0.45 + s.Stamina*0.35     // stamina pushes fade later
			fadeAmt := 0.55 - s.Stamina*0.35       // stamina reduces fade severity
			if fadeAmt < 0.15 {
				fadeAmt = 0.15
			}
			if fadeAmt > 0.65 {
				fadeAmt = 0.65
			}

			burstChance := 0.02 + s.Burst*0.06 // burstier horse procs more often
			burstBoost := 0.30 + s.Burst*0.80
			burstCDTicks := 22 - int(s.Burst*10.0) // more burst -> shorter cooldown
			if burstCDTicks < 8 {
				burstCDTicks = 8
			}

			tickMin := time.Duration(35) * time.Millisecond
			tickMax := time.Duration(70) * time.Millisecond

			// gate skill: small slow start if Gate is weak
			gatePenaltyTicks := 6 - int(s.Gate*4.0) // up to ~6 early “hesitant” ticks
			if gatePenaltyTicks < 0 {
				gatePenaltyTicks = 0
			}
			hesitantTicks := gatePenaltyTicks

			/*baseStep := 1.00 + rand.Float64()*0.40
			stepVar := 0.25
			fadeAfter := 0.60
			fadeAmt := 0.35
			burstChance := 0.06
			burstBoost := 0.80
			burstCDTicks := 18

			tickMin := 35 * time.Millisecond
			tickMax := 70 * time.Millisecond */

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
				if hesitantTicks > 0 {
					step *= 0.65 + 0.10*rand.Float64()
					hesitantTicks--
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
