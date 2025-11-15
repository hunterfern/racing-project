package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// In-memory store for the most recent race(s)
var (
	recentRaces []string
	recentMu    sync.Mutex
)

// Call this whenever a race finishes to make it appear immediately
func AddRecentRace(raceText string) {
	recentMu.Lock()
	defer recentMu.Unlock()
	recentRaces = append(recentRaces, raceText)
}

// Serve the /results page
func resultsHandler() {
	http.HandleFunc("/results", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store, max-age=0")
		w.Header().Set("Pragma", "no-cache")

		// Read from file
		raw, err := os.ReadFile("race_log.txt")
		if err != nil {
			fmt.Fprintf(w, "<h1>Race Results</h1><p>Error reading race_log.txt: %v</p>", err)
			return
		}

		// Prepend any recent races that might not yet be flushed to the file
		recentMu.Lock()
		if len(recentRaces) > 0 {
			raw = []byte(strings.Join(recentRaces, "\n") + "\n" + string(raw))
		}
		recentMu.Unlock()

		// Split by race headings
		races := strings.Split(string(raw), "=== Race @")
		// Reverse slice to show newest races first
		for i, j := 0, len(races)-1; i < j; i, j = i+1, j-1 {
			races[i], races[j] = races[j], races[i]
		}

		var html strings.Builder

		for _, rblock := range races {
			rblock = strings.TrimSpace(rblock)
			if rblock == "" {
				continue
			}

			// Extract timestamp
			timeRe := regexp.MustCompile(`(.+?) ===`)
			timeMatch := timeRe.FindStringSubmatch(rblock)
			timestamp := "Unknown"
			if len(timeMatch) > 1 {
				timestamp = strings.TrimSpace(timeMatch[1])
			}

			// Extract winner
			winRe := regexp.MustCompile(`Winner:\s*Racer\s*(\d+)\s*\(([\d\.]+)s\)`)
			win := winRe.FindStringSubmatch(rblock)

			winnerID := "?"
			winnerTime := "?"
			if len(win) == 3 {
				rid, _ := strconv.Atoi(win[1])
				rid++ // start at 1
				winnerID = strconv.Itoa(rid)
				winnerTime = win[2]
			}

			// Extract placing lines
			placingRe := regexp.MustCompile(`\d+\)\s*Racer\s*(\d+)\s*([\d\.]+)s`)
			placings := placingRe.FindAllStringSubmatch(rblock, -1)

			// Start race card
			html.WriteString(fmt.Sprintf(`
				<div class="race-card">
					<h2>Race @ %s</h2>
					<p class="winner-line"><strong>Winner:</strong> Racer %s (%ss)</p>
					<table>
						<tr><th>Place</th><th>Racer</th><th>Time</th></tr>
			`, timestamp, winnerID, winnerTime))

			placeNum := 1
			for _, p := range placings {
				if len(p) != 3 {
					continue
				}

				racerNum, _ := strconv.Atoi(p[1])
				racerNum++ // Racer numbers start at 1
				time := p[2]

				color := ""
				switch placeNum {
				case 1:
					color = "gold"
				case 2:
					color = "gray"
				case 3:
					color = "#cd7f32" // bronze
				}

				html.WriteString(fmt.Sprintf(`
					<tr style="color:%s">
						<td>%d</td>
						<td>Racer %d</td>
						<td>%ss</td>
					</tr>
				`, color, placeNum, racerNum, time))

				placeNum++
			}

			html.WriteString("</table></div>")
		}

		// Render page
		fmt.Fprintf(w, `
<!doctype html>
<html>
<head>
<meta charset="utf-8">
<title>Race Results</title>

<style>
body { 
  font-family: system-ui, sans-serif; 
  padding: 20px; 
  background: linear-gradient(135deg, #74ABE2, #3d4de3); 
  color: white;
}

h1 { 
  text-align: center; 
  margin-bottom: 20px;
  font-size: 2.2rem;
  text-shadow: 0 2px 5px rgba(0,0,0,0.3);
}

.race-card {
  background: rgba(255,255,255,0.18);
  backdrop-filter: blur(6px);
  border-radius: 14px;
  padding: 18px 22px;
  margin: 25px auto;
  max-width: 650px;
  box-shadow: 0 4px 14px rgba(0,0,0,0.25);
}

.race-card h2 {
  margin-top: 0;
  text-align: center;
  margin-bottom: 12px;
}

.winner-line {
  text-align: center;
  margin-bottom: 12px;
  font-size: 1.1rem;
  color: #ffebc2;
}

table {
  width: 100%%;
  border-collapse: collapse;
}

th, td {
  padding: 10px;
  text-align: center;
  border-bottom: 1px solid rgba(255,255,255,0.3);
}

th {
  font-weight: bold;
  font-size: 1.05rem;
}

button { 
  margin-bottom: 12px; 
  background-color: #ffffff; 
  color: black; 
  border: none; 
  padding: 8px 14px; 
  border-radius: 8px; 
  cursor: pointer;
  font-size: 1rem;
  transition: 0.15s ease;
  box-shadow: 0 2px 6px rgba(0,0,0,0.3);
}
button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 10px rgba(0,0,0,0.35);
}
</style>
</head>

<body>

<h1>Race Results</h1>

<button onclick="window.location.href='/'">Home</button>
<button onclick="location.href='/track'">Track 1</button>
<button onclick="location.href='/track2'">Track 2</button>

%s

</body>
</html>
`, html.String())
	})
}
