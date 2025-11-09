package main

import "fmt"

func generateTrackPageHTML(defaultNumRacers int) string {
	if defaultNumRacers < 4 {
		defaultNumRacers = 4
	} else if defaultNumRacers > 12 {
		defaultNumRacers = 12
	}

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
	width: 400px; 
	height: 250px; 
	margin: 40px auto; 
	background: #a2bf63; 
	border-radius: 50% / 30%; 
	border: 2px solid #ccc; }
.racer-dot { 
	position: absolute; 
	width: 40px; 
	height: 40px; 
	object-fit: contain; 
	transform: translate(-50%, -50%); 
	transition: transform 0.1s linear; 
	pointer-events: none; }
#winner { margin-top: 20px; font-size: 1.2em; font-weight: bold; }
#toolbar { margin-bottom: 12px; }
button, select { margin: 4px; padding: 8px 14px; font-size: 1em; border-radius: 6px; border: 1px solid #aaa; background: #fff; cursor: pointer; }
button:hover, select:hover { background: #f0f0f0; }
</style>
</head>
<body>
<h1>🏁 Oval Race Track</h1>

<div id="toolbar">
  <label for="numRacers"># of Horses:</label>
  <select id="numRacers">
    <option value="4">4</option>
    <option value="5">5</option>
    <option value="6">6</option>
    <option value="7">7</option>
    <option value="8">8</option>
    <option value="9">9</option>
    <option value="10">10</option>
    <option value="11">11</option>
    <option value="12">12</option>
  </select>
  <button id="start" disabled>Start Race</button>
  <button id="homeBtn" style="margin-left:8px;">Home</button>
  <button id="resultsBtn" style="margin-left:8px;">Results</button>
  <span id="status" style="margin-left:8px; color:#555;">connecting…</span>
</div>

<div class="track-wrapper">
  <div class="lane-numbers" id="laneNumbers"></div>
  <div id="track-container"></div>
</div>

<div id="winner">Winner: —</div>

<script>
let NUM_RACERS = ` + fmt.Sprint(defaultNumRacers) + `;

const track = document.getElementById('track-container');
const dots = [];
const selectRacers = document.getElementById('numRacers');

// create racers
function createRacers() {
	track.innerHTML = '';
	dots.length = 0;
	for (let i = 0; i < NUM_RACERS; i++) {
		const img = document.createElement('img');
		img.className = 'racer-dot';
		const idx = (i % 12) + 1;
		img.src = '/racer_pictures/racer' + idx + '.png';
		img.alt = 'Horse ' + (i + 1);
		track.appendChild(img);
		dots.push(img);
	}
}

const rx = 180, ry = 100, cx = 1, cy = 125;
function getAngle(pct) { return pct / 100 * 2 * Math.PI; }
function updateTrack(progress) {
	for (let i = 0; i < dots.length; i++) {
		const pct = (progress && progress[i] !== undefined) ? progress[i] : (initialPositions[i] || 0);
		const angle = getAngle(pct);
		const x = cx + rx * Math.cos(angle - Math.PI / 2);
		const y = cy + ry * Math.sin(angle - Math.PI / 2);
		dots[i].style.transform = "translate(" + x + "px, " + y + "px) translate(-50%, -50%)";
	}
}

// handle winner objects or numeric IDs and show time if available
function showWinner(w) {
	const winnerEl = document.getElementById('winner');
	if (!w) {
		winnerEl.textContent = "Winner: —";
		return;
	}

	let id, timeText = "";
	if (typeof w === 'object' && w.id !== undefined) {
		id = Number(w.id);
		if (w.finishMs !== undefined) {
			const secs = (Number(w.finishMs) / 1000).toFixed(3);
			timeText = " — " + secs + "s";
		}
	} else {
		id = Number(w);
	}
	if (isNaN(id)) {
		winnerEl.textContent = "Winner: —";
		return;
	}
	winnerEl.textContent = "Winner: Horse " + (id + 1) + timeText;
	document.getElementById('start').textContent = "Race Again";
}

const proto = location.protocol === 'https:' ? 'wss://' : 'ws://';
const ws = new WebSocket(proto + location.host + '/ws');

const btnStart = document.getElementById('start');
const statusEl = document.getElementById('status');

ws.onopen = () => {
	statusEl.textContent = 'connected';
	btnStart.disabled = false;
};
ws.onclose = () => {
	statusEl.textContent = 'disconnected';
	btnStart.disabled = true;
};

// Start or restart race
btnStart.onclick = async () => {
	if (ws.readyState === WebSocket.OPEN) {
		const n = parseInt(selectRacers.value);
		NUM_RACERS = n;
		createRacers();
		ws.send(JSON.stringify({ type: 'START', racers: n }));
		document.getElementById('winner').textContent = "";
		document.getElementById('start').textContent = "Start Race";
	}
};

// Track updates
ws.onmessage = (e) => {
	try {
		const msg = JSON.parse(e.data);
		if (msg.type === 'progress') updateTrack(msg.data || []);
		else if (msg.type === 'winner') showWinner(msg.data); // ✅ works for both old/new server formats
	} catch (err) {
		console.error(err);
	}
};

const initialPositions = [];
(function() {
	const startPct = 0;
	const gap = 3.5;
	for (let i = 0; i < 12; i++) {
		initialPositions.push((startPct - i * gap + 100) % 100);
	}
})();

window.addEventListener('load', () => {
	createRacers();
	updateTrack();
	selectRacers.value = NUM_RACERS;
});

// Navigation
document.getElementById("homeBtn").onclick = () => { window.location.href = "/landing"; };
document.getElementById("resultsBtn").onclick = () => { window.location.href = "/results"; };
</script>
</body>
</html>`
	return html
}
