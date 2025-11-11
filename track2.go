package main

import "fmt"

func generateTrack2PageHTML(defaultNumRacers int) string {
	if defaultNumRacers < 4 {
		defaultNumRacers = 4
	} else if defaultNumRacers > 12 {
		defaultNumRacers = 12
	}

	trackHeight := 55 * defaultNumRacers // total track height
	laneHeight := 55                     // height per lane

	html := `<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8" />
<meta name="viewport" content="width=device-width,initial-scale=1" />
<title>Race Track 2</title>
<style>
:root {
  --track-width: 80%;
  --track-bg: #a2bf63;
  --track-radius: 10px;
  --track-shadow: 0 0 15px rgba(0,0,0,0.5);
  --lane-line: rgba(255,255,255,0.12);
  --lane-number-color: #ffffffff;
}
body {
  font-family: system-ui, -apple-system, "Segoe UI", Roboto, "Helvetica Neue", Arial;
  text-align: center;
  background: linear-gradient(135deg, #74ABE2, #5563DE);
  color: #fff;
  margin: 0;
  padding: 20px;
}
h1 { margin-top: 8px; }

.track-wrapper {
  position: relative;
  width: 90%;
  max-width: 1100px;
  margin: 20px auto;
  display: flex;
  align-items: flex-start;
  justify-content: center;
}

.lane-numbers {
  width: 60px;
  margin-right: 8px;
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
  text-shadow: 0 1px 2px #000;
  box-sizing: border-box;
}

#track {
  position: relative;
  width: var(--track-width);
  background: var(--track-bg);
  height: ` + fmt.Sprint(trackHeight) + `px;
  border-radius: var(--track-radius);
  overflow: hidden;
  box-shadow: var(--track-shadow);
  border: 2px solid rgba(255,255,255,0.06);
}

.lane-line {
  position: absolute;
  left: 0;
  right: 0;
  height: 1px;
  background: var(--lane-line);
}

.racer-dot {
  position: absolute;
  width: 40px;
  height: 40px;
  object-fit: contain;
  transform: translateY(-50%);
  transition: left 0.1s linear;
  pointer-events: none;
}

.finish-line {
  position: absolute;
  top: 0;
  right: 10%;
  width: 6px;
  height: 100%;
  background: repeating-linear-gradient(
    45deg,
    #fff,
    #fff 6px,
    #000 6px,
    #000 12px
  );
  opacity: 0.9;
}

#winner {
  margin-top: 16px;
  font-size: 20px;
  font-weight: 700;
  color: #ffffffff;
}

#toolbar {
  margin-bottom: 12px;
}

select, button {
  margin: 4px;
  padding: 8px 14px;
  font-size: 1em;
  border-radius: 6px;
  border: 1px solid #aaa;
  background: #fff;
  cursor: pointer;
}
select:hover, button:hover { background: #f0f0f0; }

.leaderboard {
  width: 260px;
  margin-left: 12px;
  background: rgba(0,0,0,0.35);
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 10px;
  padding: 12px;
  box-shadow: 0 4px 14px rgba(0,0,0,0.3);
}

.leaderboard h2 {
  margin: 0 0 8px 0;
  font-size: 16px;
  letter-spacing: 0.5px;
  color: #ffffffff;
}

#lbBody {  position: relative; height: 0; }

.lb-row {
  position: absolute;
  left: 8px; right: 8px;
  height: 44px;
  display: grid;
  grid-template-columns: 28px 1fr;
  gap: 8px;
  padding: 6px 8px;
  border-bottom: 1px solid rgba(255,255,255,0);
  align-items: center;
  font-weight: 600;
  background: rgba(0,0,0,0);
  border-radius: 8px;
  transition: transform 600ms cubic-bezier(.25,.8,.25,1);
  will-change: transform;
}

.lb-row:last-child { border-bottom: none; }
.lb-rank { text-align: right; opacity: 0.9; }
.lb-name { text-align: left; }

#odds .odds-row{
  display: flex;
  justify-content: space-between;
  padding: 6px 8px;
  border-bottom: 1px solid rgba(255,255,255,0.06);
  font-weight: 600;
}
#odds .odds-row:last-child{ border-bottom:none; }


@media (max-width: 700px) {
  .lane-numbers { display: none; }
  .finish-line { right: 10%; width: 5px; }
  .racer-dot { width: 32px; height: 32px; }
}
</style>
</head>
<body>
<h1>Race Track</h1>
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
  <button id="homeBtn">Home</button>
  <button id="resultsBtn">Results</button>
  <span id="status" style="margin-left:12px;color:#bbb">connecting…</span>
</div>

<div class="track-wrapper">
  <div class="lane-numbers" id="laneNumbers"></div>
  <div id="track"></div>
  <div class="leaderboard" id="leaderboard">
    <h2>Leaderboard</h2>
    <div id="lbBody"></div>
  </div>
  <div class="leaderboard" id="oddsPanel">
    <h2>Odds</h2>
    <div id="odds"></div>
  </div>
</div>

<div id="winner">Winner: —</div>

<script>
let NUM_RACERS = ` + fmt.Sprint(defaultNumRacers) + `;
const laneHeight = ` + fmt.Sprint(laneHeight) + `;
const trackEl = document.getElementById('track');
const laneNumbersEl = document.getElementById('laneNumbers');
const winnerEl = document.getElementById('winner');
const btnStart = document.getElementById('start');
const selectRacers = document.getElementById('numRacers');
const statusEl = document.getElementById('status');
const homeBtn = document.getElementById('homeBtn');
const lbBody = document.getElementById('lbBody');
const oddsEl = document.getElementById('odds');
let lastOrder = [];
const ROW_H = 44;
const rowEls = [];
let lastLBUpdate = 0;
const LB_INTERVAL = 300;

const dots = [];

// draw horses and lanes
function createRacers() {
  const newTrackHeight = laneHeight * NUM_RACERS;
  trackEl.style.height = newTrackHeight + 'px';
  trackEl.innerHTML = '';
  laneNumbersEl.innerHTML = '';
  dots.length = 0;

  for (let i = 0; i < NUM_RACERS; i++) {
    const ln = document.createElement('div');
    ln.className = 'lane-number';
    ln.textContent = '#' + (i + 1);
    ln.style.height = laneHeight + 'px';
    laneNumbersEl.appendChild(ln);

    if (i > 0) {
      const line = document.createElement('div');
      line.className = 'lane-line';
      line.style.top = (i * laneHeight) + 'px';
      trackEl.appendChild(line);
    }

    const img = document.createElement('img');
    img.className = 'racer-dot';
    img.src = '/racer_pictures/racer' + ((i % 12) + 1) + '.png';
    img.alt = 'Racer ' + (i + 1);
    img.style.top = (i * laneHeight + laneHeight / 2) + 'px';
    img.style.left = '0px';
    trackEl.appendChild(img);
    dots.push(img);
  }

  const finish = document.createElement('div');
  finish.className = 'finish-line';
  trackEl.appendChild(finish);

  buildLeaderboardRows();
}

function getSuffix(n) {
  if (n % 10 === 1 && n % 100 !== 11) return 'st';
  if (n % 10 === 2 && n % 100 !== 12) return 'nd';
  if (n % 10 === 3 && n % 100 !== 13) return 'rd';
  return 'th';
}

function buildLeaderboardRows() {
  lbBody.innerHTML = '';
  rowEls.length = 0;

  lbBody.style.height = (ROW_H * NUM_RACERS) + 'px';

  for (let i = 0; i < NUM_RACERS; i++) {
    const row = document.createElement('div');
    row.className = 'lb-row';
    row.style.transform = 'translateY(' + (i * ROW_H) + 'px)'; 
    row.dataset.id = i;

    row.innerHTML =
  '<div class="lb-rank">' + (i + 1) + getSuffix(i + 1) + '</div>' +
  '<div class="lb-name">Racer ' + (i + 1) + '</div>';

    lbBody.appendChild(row);
    rowEls.push(row);
  }

  lastOrder = Array.from({length: NUM_RACERS}, (_, k) => k);
}

// update positions
function updateTrack(progress) {
  const trackWidth = trackEl.clientWidth;
  const maxX = trackWidth * 0.94 - 40;
  for (let i = 0; i < dots.length; i++) {
    const pct = progress && progress[i] !== undefined ? progress[i] : 0;
    dots[i].style.left = Math.max(0, Math.min(maxX, maxX * pct / 100)) + 'px';
  }
updateLeaderboard(progress);
}

function updateLeaderboard(progress) {
  const now = performance.now();
  if (now - lastLBUpdate < LB_INTERVAL) return;
  lastLBUpdate = now;
  const rows = [];
  for ( let i = 0; i < NUM_RACERS; i++) {
    const pct = progress && progress[i] !== undefined ? progress[i] : 0;
    rows.push([i, pct]);
  }

  const orderIndex = new Map();
  if (lastOrder.length === NUM_RACERS) {
    lastOrder.forEach((id, idx) => orderIndex.set(id, idx));
  }
  rows.sort((a, b) => {
    if (b[1] !== a[1]) return b[1] - a[1];
    const ao = orderIndex.has(a[0]) ? orderIndex.get(a[0]) : 9999;
    const bo = orderIndex.has(b[0]) ? orderIndex.get(b[0]) : 9999;
    if (ao !== bo) return ao - bo;
    return a[0] - b[0];
  });

  for (let rank = 0; rank < rows.length; rank++) {
    const id = rows[rank][0];
    const row = rowEls[id];
    if (!row) continue;

    row.querySelector('.lb-rank').textContent = (rank + 1) + getSuffix(rank + 1);

    row.style.transform = 'translateY(' + (rank * ROW_H) + 'px)';
  }

  lastOrder = rows.map(r => r[0]);
}


// display winner
function showWinner(w) {
  if (!w || (typeof w === 'object' && w.id === undefined)) {
    winnerEl.textContent = 'Winner: —';
    return;
  }
  if (typeof w === 'object') {
    const secs = (Number(w.finishMs)/1000).toFixed(3);
    winnerEl.textContent = 'Winner: Racer ' + (w.id+1) + ' — ' + secs + 's';
  } else {
    winnerEl.textContent = 'Winner: Racer ' + (Number(w)+1);
  }
  btnStart.textContent = 'Race Again';
}

// WebSocket
const proto = location.protocol === 'https:' ? 'wss://' : 'ws://';
const ws = new WebSocket(proto + location.host + '/ws');

ws.addEventListener('open', () => { statusEl.textContent='connected'; btnStart.disabled=false; });
ws.addEventListener('close', () => { statusEl.textContent='disconnected'; btnStart.disabled=true; });

btnStart.addEventListener('click', () => {
  if (ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify({ type: 'START', racers: NUM_RACERS }));
    showWinner(null);
    btnStart.textContent = 'Racing...';
    btnStart.disabled = true;
    setTimeout(() => { btnStart.disabled = false; }, 1500);
  }
});

// redraw track immediately on dropdown change
selectRacers.addEventListener('change', () => {
  NUM_RACERS = parseInt(selectRacers.value);
  createRacers();
  buildLeaderboardRows();
  showWinner(null);
  updateLeaderboard(new Array(NUM_RACERS).fill(0));
  btnStart.textContent = 'Start Race';
});

function renderOdds(lines){
  if (!Array.isArray(lines) || !oddsEl) return;

  // Sort by best price (lowest absolute American odds first)
  var sorted = lines.slice().sort(function(a, b){
    var aa = Math.abs(a.amer);
    var bb = Math.abs(b.amer);
    return aa - bb;
  });

  oddsEl.innerHTML = sorted.map(function(l){
    // american odds string
    var sign = (l.amer >= 0 ? '+' : '') + String(l.amer);

    // fractional odds from Go: fracN/fracD (note the field names)
    var frac = (l.fracN && l.fracD) ? (' (' + l.fracN + '/' + l.fracD + ')') : '';

    // build one row using plain string concatenation
    return '<div class="odds-row">' +
             '<span>#' + (l.id + 1) + ' Racer ' + (l.id + 1) + '</span>' +
             '<span>' + sign + frac + '</span>' +
           '</div>';
  }).join('');
}

// handle progress/winner from server
ws.addEventListener('message', (evt) => {
  try {
    const msg = JSON.parse(evt.data);
    if (msg.type==='progress'){
      updateTrack(msg.data||[]);
      updateLeaderboard(msg.data || []);
    }
    else if (msg.type==='winner') showWinner(msg.data);
    else if (msg.type==='lines') renderOdds(msg.data);
  } catch(err){console.error(err);}
});


// navigation
homeBtn.addEventListener('click', () => { btnStart.textContent='Start Race'; window.location.href='/'; });
document.getElementById('resultsBtn').addEventListener('click', ()=>window.location.href='/results');

window.addEventListener('resize', ()=>updateTrack([]));
window.addEventListener('load', () => { createRacers(); selectRacers.value=NUM_RACERS; buildLeaderboardRows(); updateTrack([]); updateLeaderboard(new Array(NUM_RACERS).fill(0)); });
</script>
</body>
</html>`

	return html
}
