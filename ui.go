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

	// Track finish times and broadcast winner when all finish
	go func() {
		// Maps to store finish times. Access guarded by mu.
		finishTimes := make(map[int]time.Duration)
		winnerSent := false

		for u := range updates {
			if u.position < 0 {
				u.position = 0
			}
			if u.position > 100 {
				u.position = 100
			}

			// update progress slice
			mu.Lock()
			for len(progress) <= u.id {
				progress = append(progress, 0)
			}
			progress[u.id] = u.position

			// if this update marks the racer finished and we haven't recorded it yet,
			// record their elapsed time
			if u.finished {
				if _, ok := finishTimes[u.id]; !ok {
					finishTimes[u.id] = u.elapsed
				}
			}

			// if all racers have finished (and we have a non-zero expected count),
			// determine winner and broadcast (only once).
			if !winnerSent && uiNumRacers > 0 && len(finishTimes) >= uiNumRacers {
				// Find minimum elapsed time (tie-breaker: lower id)
				winnerID := -1
				var best time.Duration
				for id, t := range finishTimes {
					if winnerID == -1 || t < best || (t == best && id < winnerID) {
						winnerID = id
						best = t
					}
				}
				// broadcast winner over websocket with finish time in milliseconds
				payload := map[string]interface{}{
					"id":       winnerID,
					"finishMs": best.Milliseconds(),
				}
				hub.broadcast <- mustJSON(WSMessage{Type: "winner", Data: payload})
				winnerSent = true
			}
			mu.Unlock()
		}
	}()

	http.HandleFunc("/ws", hub.WSHandler(commands))

	http.HandleFunc("/ws-test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `
<!doctype html><meta charset="utf-8">
<h1>WS Test</h1>
<pre id="out"></pre>
<script>
 const out = document.getElementById('out');
 const proto = location.protocol === 'https:' ? 'wss://' : 'ws://';
 const ws = new WebSocket(proto + location.host + '/ws');
 ws.onopen = () => out.textContent = 'connected...';
 ws.onmessage = (e) => { out.textContent = e.data; };
 ws.onclose = () => out.textContent += '\nclosed';
</script>`)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store, max-age=0")
		w.Header().Set("Pragma", "no-cache")

		fmt.Fprintf(w, `
<!doctype html>
<meta charset="utf-8" />
<title>Race Progress</title>
<style>
  body { font-family: system-ui, sans-serif; padding: 16px; }
  .progress-container { width: 500px; height: 24px; border: 1px solid #000; background:#eee; margin: 8px 0; position: relative; }
  .progress-bar { height: 100%%; background:#4caf50; width:0%%; transition: width 80ms linear; }
  .label { position:absolute; left:8px; top:0; height:100%%; display:flex; align-items:center; font-weight:600; }
  #toolbar { margin-bottom: 12px; }
</style>

<h1>Race Progress</h1>
<div id="toolbar">
  <button id="start" disabled>Start Race</button>
  <button id="trackBtn" style="margin-left:8px;">Track</button>
  <button id="track2Btn" style="margin-left:8px;">Track 2</button>
  <button id="resultsBtn" style="margin-left:8px;">Results</button>
  <span id="status" style="margin-left:8px; color:#555;">connecting…</span>
</div>
<div id="bars"></div>

<script>
document.getElementById("trackBtn").onclick = () => {
    window.location.href = "/track";
};

document.getElementById("resultsBtn").onclick = () => {
    window.location.href = "/results";
};

document.getElementById("track2Btn").onclick = () => {
    window.location.href = "/track2";
};
</script>


<script>
  const container = document.getElementById("bars");
  const btnStart  = document.getElementById("start");
  const statusEl  = document.getElementById("status");

  function renderBars(n){
    container.innerHTML = "";
    for (let i = 0; i < n; i++) {
      const div = document.createElement("div");
      div.className = "progress-container";
      div.innerHTML = '<div id="bar' + i + '" class="progress-bar"></div><span class="label"> ' + (i+1) + '</span>';
      container.appendChild(div);
    }
  }

  const proto = location.protocol === 'https:' ? 'wss://' : 'ws://';
  const ws = new WebSocket(proto + location.host + '/ws');

  ws.onopen = () => {
    statusEl.textContent = 'connected';
    btnStart.disabled = false;
  };

  ws.onclose = () => {
    statusEl.textContent = 'disconnected';
    btnStart.disabled = true;
  };

  function send(type, payload={}) {
    if (ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({ type, ...payload }));
    } else {
      console.warn('WS not open');
    }
  }

  btnStart.onclick = () => {
    send('START');
  };

  // Single onmessage handler
  ws.onmessage = (e) => {
    try {
      const msg = JSON.parse(e.data);
      if (msg.type !== 'progress') return;
      const data = msg.data || [];
      if (container.children.length !== data.length) renderBars(data.length);
      for (let i = 0; i < data.length; i++) {
        const bar = document.getElementById("bar" + i);
        if (bar) bar.style.width = data[i] + "%%";
      }
    } catch (err) {
      console.error('parse error', err);
    }
  };
</script>
`)

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
			for i := 0; i < uiNumRacers && i < len(progress); i++ {
				cp[i] = progress[i]
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
    	w.Write([]byte(generateTrackPageHTML(uiNumRacers)))
	})

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			mu.RLock()
			cp := make([]int, uiNumRacers)
			for i := 0; i < uiNumRacers && i < len(progress); i++ {
				cp[i] = progress[i]
			}
			mu.RUnlock()
			hub.broadcast <- mustJSON(WSMessage{Type: "progress", Data: cp})
		}
	}()

	http.HandleFunc("/track2", func(w http.ResponseWriter, r *http.Request) {
    	w.Header().Set("Content-Type", "text/html; charset=utf-8")
    	w.Write([]byte(generateTrack2PageHTML(uiNumRacers)))
	})


	http.Handle("/racer_pictures/", http.StripPrefix("/racer_pictures/", http.FileServer(http.Dir("./racer_pictures"))))

	resultsHandler() 

	go func() {
		fmt.Println("Server started at http://localhost:8080")
		_ = http.ListenAndServe("0.0.0.0:8080", nil)
	}()
}
