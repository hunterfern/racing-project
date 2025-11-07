package main

import "fmt"

// generateTrack2PageHTML creates the HTML for the left-to-right race track
func generateTrack2PageHTML(numRacers int) string {
	if numRacers < 1 {
		numRacers = 1
	}

	trackHeight := 55 * numRacers // total track height
	laneHeight := 55             // height per lane

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
      --lane-number-color: #ffea00;
    }
    body {
      font-family: system-ui, -apple-system, "Segoe UI", Roboto, "Helvetica Neue", Arial;
      text-align: center;
      background: #111;
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

    /* Left column for lane numbers */
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
      text-shadow: 0 1px 2px #a2bf63);
      box-sizing: border-box;
    }

    /* Track container */
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

    .lane {
      position: absolute;
      left: 0;
      right: 0;
      height: ` + fmt.Sprint(laneHeight) + `px;
      box-sizing: border-box;
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
      color: #ffd54d;
    }

    #toolbar {
      margin-bottom: 12px;
    }

    /* small responsive adjustments */
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
    <button id="start" disabled>Start Race</button>
    <button id="backBtn">Progress</button>
    <button id="resultsBtn">Results</button>
    <span id="status" style="margin-left:12px;color:#bbb">connecting…</span>
  </div>

  <div class="track-wrapper">
    <div class="lane-numbers" id="laneNumbers"></div>
    <div id="track"></div>
  </div>

  <div id="winner">Winner: —</div>

<script>
const NUM_RACERS = ` + fmt.Sprint(numRacers) + `;
const laneHeight = ` + fmt.Sprint(laneHeight) + `;
const trackEl = document.getElementById('track');
const laneNumbersEl = document.getElementById('laneNumbers');
const winnerEl = document.getElementById('winner');

// create lane numbers and lane lines and racer image placeholders
const dots = [];
for (let i = 0; i < NUM_RACERS; i++) {
  // lane number element (fixed left column)
  const ln = document.createElement('div');
  ln.className = 'lane-number';
  ln.textContent = '#' + (i + 1);
  ln.style.height = laneHeight + 'px';
  laneNumbersEl.appendChild(ln);

  // lane container positioning (we won't add child lane divs — we draw lines and dots)
  const top = i * laneHeight;
  // lane line (separator) except topmost
  if (i > 0) {
    const line = document.createElement('div');
    line.className = 'lane-line';
    line.style.top = top + 'px';
    trackEl.appendChild(line);
  }

  // create racer image element
  const img = document.createElement('img');
  img.className = 'racer-dot';
  // use your racer picture files if available, fallback to colored circle data URL if not
  const idx = (i % 12) + 1;
  img.src = '/racer_pictures/racer' + idx + '.png';
  img.alt = 'Racer ' + (i + 1);
  // vertical center of lane
  img.style.top = (top + laneHeight / 2) + 'px';
  img.style.left = '0px';
  trackEl.appendChild(img);
  dots.push(img);
}

// draw finish line
const finish = document.createElement('div');
finish.className = 'finish-line';
trackEl.appendChild(finish);

// helper to update positions (progress array of 0..100)
function updateTrack(progress) {
  const trackWidth = trackEl.clientWidth;
  const maxX = trackWidth * 0.94 - 40; // leave some margin and image width, finish-line at ~94%
  for (let i = 0; i < dots.length; i++) {
    const pct = (progress && progress[i] !== undefined) ? progress[i] : 0;
    const x = Math.max(0, Math.min(maxX, (maxX * pct / 100)));
    dots[i].style.left = x + 'px';
  }
}

// show winner payload: either {id, finishMs} or numeric id
function showWinner(w) {
  if (w == null) {
    winnerEl.textContent = 'Winner: —';
    return;
  }
  if (typeof w === 'object' && w.id !== undefined) {
    const secs = (Number(w.finishMs) / 1000).toFixed(3);
    winnerEl.textContent = 'Winner: Racer ' + (w.id + 1) + ' — ' + secs;
  } else {
    winnerEl.textContent = 'Winner: Racer ' + (Number(w) + 1) ;
  }
}

// WebSocket setup (reuse your /ws endpoint)
const proto = location.protocol === 'https:' ? 'wss://' : 'ws://';
const ws = new WebSocket(proto + location.host + '/ws');

const btnStart = document.getElementById('start');
const statusEl = document.getElementById('status');

ws.addEventListener('open', () => {
  statusEl.textContent = 'connected';
  btnStart.disabled = false;
});

ws.addEventListener('close', () => {
  statusEl.textContent = 'disconnected';
  btnStart.disabled = true;
});

btnStart.addEventListener('click', () => {
  if (ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify({ type: 'START' }));
    // clear previous winner text when starting new race
    showWinner(null);
  }
});

ws.addEventListener('message', (evt) => {
  try {
    const msg = JSON.parse(evt.data);
    if (msg.type === 'progress') {
      updateTrack(msg.data || []);
    } else if (msg.type === 'winner') {
      showWinner(msg.data);
    }
  } catch (err) {
    console.error('ws parse error', err);
  }
});

// navigation buttons
document.getElementById('backBtn').addEventListener('click', () => { window.location.href = '/'; });
document.getElementById('resultsBtn').addEventListener('click', () => { window.location.href = '/results'; });

// initial layout update on load/resize
window.addEventListener('resize', () => updateTrack([]));
window.addEventListener('load', () => updateTrack([]));
</script>
</body>
</html>`

	return html
}
