package main

import "fmt"

func generateTrackPageHTML(numRacers int) string {
	if numRacers < 1 {
		numRacers = 1
	}

	trackHeight := 55 * numRacers
	laneHeight := 50

	html := `<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8" />
<meta name="viewport" content="width=device-width,initial-scale=1" />
<title>Oval Race Track</title>
<style>
:root{
  --track-bg: #a2bf63;
  --panel-bg: #111;
  --lane-number-color: #ffea00;
  --lane-line: rgba(255,255,255,0.10);
}
body {
  margin: 0;
  padding: 20px;
  background: var(--panel-bg);
  color: #fff;
  font-family: system-ui, -apple-system, "Segoe UI", Roboto, "Helvetica Neue", Arial;
  text-align: center;
}
h1 { margin: 8px 0 12px; }

.track-wrapper {
  display: flex;
  align-items: flex-start;
  justify-content: center;
  gap: 12px;
  width: 95%;
  max-width: 1100px;
  margin: 0 auto 8px;
}

.lane-numbers {
  width: 60px;
  text-align: right;
  user-select: none;
}
.lane-number {
  height: ` + fmt.Sprint(laneHeight) + `px;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  padding-right: 8px;
  font-weight: 700;
  color: var(--lane-number-color);
  text-shadow: 0 1px 2px rgba(0,0,0,0.6);
  box-sizing: border-box;
}

/* Oval track container */
#track-container {
  position: relative;
  width: 820px;
  height: ` + fmt.Sprint(trackHeight) + `px;
  background: var(--track-bg);
  border-radius: 50% / 30%; /* oval shape */
  border: 2px solid rgba(0,0,0,0.25);
  box-shadow: 0 6px 18px rgba(0,0,0,0.45);
  overflow: hidden;
}

/* faint horizontal lane lines (for readability) */
.lane-line {
  position: absolute;
  left: 0;
  right: 0;
  height: 1px;
  background: var(--lane-line);
}

/* racer images */
.racer-dot {
  position: absolute;
  width: 40px;
  height: 40px;
  object-fit: contain;
  transform: translate(-50%, -50%); /* center on left/top */
  transition: left 0.08s linear, top 0.08s linear;
  pointer-events: none;
}

/* finish marker (vertical) */
.finish-line {
  position: absolute;
  top: 0;
  left: 48%;
  width: 6px;
  height: 50px;
  background: repeating-linear-gradient(
    45deg,
    #fff,
    #fff 6px,
    #000 6px,
    #000 12px
  );
  transform: translateX(-50%);
  opacity: 0.9;
}

#toolbar { margin-bottom: 10px; }
button { cursor: pointer; }
#status { margin-left: 8px; color: #bbb; }

#winner {
  margin-top: 14px;
  font-size: 20px;
  font-weight: 700;
  color: #ffd54d;
}

/* responsive */
@media (max-width: 900px) {
  #track-container { width: 92%; }
  .lane-numbers { display: none; }
}
</style>
</head>
<body>
<h1>🏁 Oval Race Track</h1>

<div id="toolbar">
  <button id="start" disabled>Start Race</button>
  <button id="resultsBtn">Results</button>
  <span id="status" style="margin-left:8px; color:#bbb;">connecting…</span>
</div>

<div class="track-wrapper">
  <div class="lane-numbers" id="laneNumbers"></div>
  <div id="track-container"></div>
</div>

<div id="winner">Winner: —</div>

<script>
const NUM_RACERS = ` + fmt.Sprint(numRacers) + `;
const trackEl = document.getElementById('track-container');
const winnerEl = document.getElementById('winner');

const laneHeight = ` + fmt.Sprint(laneHeight) + `;
const rx = Math.max(140, trackEl ? trackEl.clientWidth * 0.42 : 180); // horizontal radius
const ry =  (laneHeight * NUM_RACERS) * 0.45; // vertical radius proportional to lanes
const cx = trackEl ? trackEl.clientWidth / 2 : 410; // center x (approx)
const cy = ` + fmt.Sprint(trackHeight/2) + `; // center y inside container (middle)


const dots = [];
for (let i = 0; i < NUM_RACERS; i++) {

  // create racer img (use your racer images; fallback will display nothing if missing)
  const img = document.createElement('img');
  img.className = 'racer-dot';
  const idx = (i % 12) + 1;
  img.src = '/racer_pictures/racer' + idx + '.png';
  img.alt = 'Racer ' + (i + 1);
  // start at initial offset along oval (evenly spaced)
  const initPct = (100 - (i * (100 / NUM_RACERS))) % 100;
  const ang = (initPct / 100) * Math.PI * 2;
  const x0 = (trackEl.clientWidth / 2) + rx * Math.cos(ang - Math.PI/2);
  const y0 = cy + ry * Math.sin(ang - Math.PI/2);
  img.style.left = x0 + 'px';
  img.style.top = y0 + 'px';
  trackEl.appendChild(img);
  dots.push(img);
}

// finish line (vertical marker near right side)
const finish = document.createElement('div');
finish.className = 'finish-line';
trackEl.appendChild(finish);

// helper: convert percent (0..100) to coordinates on oval (left/top)
function pctToXY(pct) {
  const w = trackEl.clientWidth;
  // recompute radii & centers based on current width/height
  const localRx = Math.max(140, w * 0.42);
  const localRy = (laneHeight * NUM_RACERS) * 0.45;
  const localCx = w / 2;
  const localCy = trackEl.clientHeight / 2;
  const angle = (pct / 100) * 2 * Math.PI;
  // offset the angle so pct=0 starts near top center (matching previous behavior)
  const x = localCx + localRx * Math.cos(angle - Math.PI/2);
  const y = localCy + localRy * Math.sin(angle - Math.PI/2);
  return { x, y };
}

// update positions from a progress array (0..100)
function updateTrack(progress) {
  // ensure layout has sized trackEl (sometimes clientWidth is zero until displayed)
  for (let i = 0; i < dots.length; i++) {
    const pct = (progress && progress[i] !== undefined) ? progress[i] : 0;
    const pos = pctToXY(pct);
    dots[i].style.left = pos.x + 'px';
    dots[i].style.top = pos.y + 'px';
  }
}

// Show winner - accepts either {id, finishMs} or numeric id
function showWinner(w) {
  if (w == null) {
    winnerEl.textContent = 'Winner: —';
    return;
  }
  if (typeof w === 'object' && w.id !== undefined) {
    const secs = (Number(w.finishMs) / 1000).toFixed(3);
    winnerEl.textContent = 'Winner: Racer ' + (w.id + 1) + ' — ' + secs ;
  } else {
    winnerEl.textContent = 'Winner: Racer ' + (Number(w) + 1) ;
  }
}
function clearWinner() { winnerEl.textContent = ''; }

// WebSocket connection & UI wiring
const proto = location.protocol === 'https:' ? 'wss://' : 'ws://';
let ws;
let btnStart = document.getElementById('start');
let restartBtn = document.getElementById('restartBtn');
let statusEl = document.getElementById('status');

function connectWS(){
  ws = new WebSocket(proto + location.host + '/ws');

  ws.onopen = () => {
    statusEl.textContent = 'connected';
    if (btnStart) btnStart.disabled = false;
  };
  ws.onclose = () => {
    statusEl.textContent = 'disconnected';
    if (btnStart) btnStart.disabled = true;
  };

  ws.onmessage = (e) => {
    try {
      const msg = JSON.parse(e.data);
      if (msg.type === 'progress') updateTrack(msg.data || []);
      else if (msg.type === 'winner') showWinner(msg.data);
      else if (msg.type === 'reset') {
        updateTrack([]);
        clearWinner();
        if (btnStart) btnStart.disabled = false;
      }
    } catch (err) {
      console.error('ws parse error', err);
    }
  };
}
connectWS();

btnStart.onclick = () => {
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify({ type: 'START' }));
    // clear winner when new race starts
    showWinner(null);
  }
};

// Keep partner's existing restart behavior (POST /reset-race) AND also send RESTART over WS.
// That way existing server-side reset endpoint is preserved while hub-based restart works too.
restartBtn.onclick = async () => {
  try {
    restartBtn.disabled = true;
    clearWinner();
    updateTrack([]);
    // try POSTing to /reset-race (preserve existing partner workflow)
    try {
      await fetch('/reset-race', { method: 'POST' });
    } catch (e) {
      // ignore fetch error (server may not implement endpoint)
      console.warn('reset endpoint POST failed (ok if not present):', e);
    }
    // also send RESTART on websocket if available
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({ type: 'RESTART' }));
    }
    // re-establish connection if ws closed (keeps previous logic)
    setTimeout(() => {
      if (!ws || ws.readyState === WebSocket.CLOSED || ws.readyState === WebSocket.CLOSING) {
        connectWS();
      }
      restartBtn.disabled = false;
    }, 300);
  } catch (e) {
    console.error('Restart failed:', e);
    restartBtn.disabled = false;
  }
};

document.getElementById("backBtn").onclick = () => { window.location.href = "/landing"; };
document.getElementById('resultsBtn').addEventListener('click', () => { window.location.href = '/results'; });

// keep UI layout up to date if the window resizes
window.addEventListener('resize', () => updateTrack([]));
window.addEventListener('load', () => updateTrack([]));
</script>
</body>
</html>`

	return html
}
