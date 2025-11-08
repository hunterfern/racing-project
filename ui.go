package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var uiNumRacers int

func displayRaceUIWithWS(updates <-chan RaceUpdate, hub *Hub, commands chan<- Command) {
	var mu sync.RWMutex
	progress := make([]int, uiNumRacers)
	finishTimes := make(map[int]time.Duration)
	var winnerSent bool

	//Process live race updates
	go func() {
		for u := range updates {
			if u.position < 0 {
				u.position = 0
			}
			if u.position > 100 {
				u.position = 100
			}

			mu.Lock()
			for len(progress) <= u.id {
				progress = append(progress, 0)
			}
			progress[u.id] = u.position

			if u.finished {
				if _, ok := finishTimes[u.id]; !ok {
					finishTimes[u.id] = u.elapsed
				}
			}

			// Determine winner when all racers finish
			if !winnerSent && uiNumRacers > 0 && len(finishTimes) >= uiNumRacers {
				winnerID := -1
				var best time.Duration
				for id, t := range finishTimes {
					if winnerID == -1 || t < best || (t == best && id < winnerID) {
						winnerID = id
						best = t
					}
				}
				payload := map[string]interface{}{
					"id":       winnerID,
					"finishMs": best.Milliseconds(),
				}
				hub.broadcast <- mustJSON(WSMessage{Type: "winner", Data: payload})
				winnerSent = true

				// Automatically reset 2s after race
				go func() {
					time.Sleep(2 * time.Second)
					resetRaceState(&mu, hub, &progress, &finishTimes, &winnerSent)
				}()
			}
			mu.Unlock()
		}
	}()

	setupRoutesOnce(hub, commands, &mu, &progress, &finishTimes, &winnerSent)
}

var setupOnce sync.Once

func setupRoutesOnce(hub *Hub, commands chan<- Command, mu *sync.RWMutex, progress *[]int, finishTimes *map[int]time.Duration, winnerSent *bool) {
	setupOnce.Do(func() {
		http.HandleFunc("/ws", hub.WSHandler(commands))

		http.HandleFunc("/reset-race", func(w http.ResponseWriter, r *http.Request) {
			resetRaceState(mu, hub, progress, finishTimes, winnerSent)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Race reset"))
		})

		// Landing page
		http.HandleFunc("/landing", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprintf(w, `
			<!doctype html>
			<meta charset="utf-8" />
			<title>Race Simulator</title>
			<style>
				body {
					font-family: system-ui, sans-serif;
					display: flex;
					flex-direction: column;
					align-items: center;
					justify-content: center;
					height: 100vh;
					background: linear-gradient(135deg, #74ABE2, #5563DE);
					color: white;
				}
				h1 {
					font-size: 2.5em;
					margin-bottom: 30px;
				}
				button {
					background: white;
					color: #333;
					border: none;
					padding: 12px 24px;
					margin: 10px;
					border-radius: 8px;
					font-size: 1.1em;
					cursor: pointer;
					transition: transform 0.15s, background 0.3s;
				}
				button:hover {
					background: #f0f0f0;
					transform: scale(1.05);
				}
			</style>
			<h1>Select a Race Track</h1>
			<div>
				<button onclick="location.href='/track'">🏁 Track 1</button>
				<button onclick="location.href='/track2'">🚗 Track 2</button>
			</div>
			`)
		})

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/landing", http.StatusSeeOther)
		})

		http.HandleFunc("/progress", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")

			ticker := time.NewTicker(100 * time.Millisecond)
			defer ticker.Stop()
			flusher, _ := w.(http.Flusher)

			for range ticker.C {
				mu.RLock()
				cp := make([]int, uiNumRacers)
				for i := 0; i < uiNumRacers && i < len(*progress); i++ {
					cp[i] = (*progress)[i]
				}
				mu.RUnlock()

				data, _ := json.Marshal(cp)
				fmt.Fprintf(w, "data: %s\n\n", data)
				if flusher != nil {
					flusher.Flush()
				}
			}
		})

		http.HandleFunc("/track", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			page := generateTrackPageHTML(uiNumRacers)
			w.Write([]byte(page))
		})

		http.HandleFunc("/track2", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			page := generateTrack2PageHTML(uiNumRacers)
			w.Write([]byte(page))
		})

		http.Handle("/racer_pictures/", http.StripPrefix("/racer_pictures/", http.FileServer(http.Dir("./racer_pictures"))))

		// Show progress over WS
		go func() {
			ticker := time.NewTicker(100 * time.Millisecond)
			defer ticker.Stop()
			for range ticker.C {
				mu.RLock()
				cp := make([]int, uiNumRacers)
				for i := 0; i < uiNumRacers && i < len(*progress); i++ {
					cp[i] = (*progress)[i]
				}
				mu.RUnlock()
				hub.broadcast <- mustJSON(WSMessage{Type: "progress", Data: cp})
			}
		}()

		resultsHandler()

		go func() {
			fmt.Println("Server started at http://localhost:8080")
			if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
				fmt.Println("Server error:", err)
			}
		}()
	})
}

func resetRaceState(mu *sync.RWMutex, hub *Hub, progress *[]int, finishTimes *map[int]time.Duration, winnerSent *bool) {
	mu.Lock()
	defer mu.Unlock()
	for i := range *progress {
		(*progress)[i] = 0
	}
	*finishTimes = make(map[int]time.Duration)
	*winnerSent = false
	hub.broadcast <- mustJSON(WSMessage{Type: "reset", Data: nil})
	fmt.Println("✅ Race reset — ready for next race")
}
