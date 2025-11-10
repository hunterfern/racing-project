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
  background: linear-gradient(135deg, #74ABE2, #5563DE);
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

/* Finish line for oval */
#finish-line {
    position: absolute;
    width: 6px;
    height: 15%;           /* keep at 90% */
    top: 0%;
    left: 48%;
    transform: translateX(-50%);
    background: repeating-linear-gradient(
        45deg,
        #fff,
        #fff 6px,
        #000 6px,
        #000 12px
    );
    border-radius: 2px;
    opacity: 0.9;
    pointer-events: none;
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
    width: 650px;      /* bigger */
    height: 420px;     /* bigger */
    margin: 40px auto; 
    background: var(--track-bg); 
    border-radius: 50% / 30%; 
    border: 2px solid #000; 
}



.racer-dot { 
    position: absolute; 
    width: 50px;      /* was 40px */
    height: 50px;     /* was 40px */
    object-fit: contain; 
    transform: translate(-50%, -50%); 
    transition: left 0.1s linear, top 0.1s linear; 
    pointer-events: none; 
}


#winner {
  margin-top: 16px;
  font-size: 20px;
  font-weight: 700;
  color: #ffffffff;
}

#toolbar { margin-bottom: 12px; }
button, select { margin: 4px; padding: 8px 14px; font-size: 1em; border-radius: 6px; border: 1px solid #aaa; background: #fff; cursor: pointer; }
button:hover, select:hover { background: #f0f0f0; }

/* Leaderboard */
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
margin:0 0 8px 0; 
font-size:16px; 
letter-spacing:0.5px; 
color:#fff; }
#lbBody { position: relative; height:0; }
.lb-row { position: absolute; 
left:8px; 
right:8px; 
height:44px; 
display:grid; 
grid-template-columns:28px 1fr; 
gap:8px; 
padding:6px 8px; 
border-bottom:1px solid rgba(255,255,255,0.00); 
align-items:center; 
font-weight:600; 
background:rgba(0,0,0,0.0); 
border-radius:8px; 
transition: transform 600ms cubic-bezier(.25,.8,.25,1); 
will-change: transform; 
}
.lb-row:last-child { border-bottom:none; }
.lb-rank { text-align:right; opacity:0.9; }
.lb-name { text-align:left; }

@media (max-width:900px) { #track-container { width: 92%; } .lane-numbers { display:none; } }
@media (max-width:700px) { 
    .racer-dot { width: 40px; height: 40px; } 
}
</style>
</head>
<body>
<h1>Oval Race Track</h1>

<div id="toolbar">
  <label for="numRacers"># of Horses:</label>
  <select id="numRacers">` +
	func() string { s := ""; for i := 4; i <= 12; i++ { s += fmt.Sprintf("<option value=\"%d\">%d</option>", i, i) }; return s }() +
	`</select>
  <button id="start" disabled>Start Race</button>
  <button id="homeBtn" style="margin-left:8px;">Home</button>
  <button id="resultsBtn" style="margin-left:8px;">Results</button>
  <span id="status" style="margin-left:8px; color:#555;">connecting…</span>
</div>

<div class="track-wrapper">
  <div class="lane-numbers" id="laneNumbers"></div>
  <div id="track-container">
  <div id="finish-line"></div>
</div>
  <div class="leaderboard" id="leaderboard">
    <h2>Leaderboard</h2>
    <div id="lbBody"></div>
  </div>
</div>

<div id="winner">Winner: —</div>

<script>
let NUM_RACERS = ` + fmt.Sprint(defaultNumRacers) + `;
const track = document.getElementById('track-container');
const dots = [];
const selectRacers = document.getElementById('numRacers');
const lbBody = document.getElementById('lbBody');
let lastOrder = [];
const ROW_H = 44;
const rowEls = [];
let lastLBUpdate = 0;
const LB_INTERVAL = 300;
const winnerEl = document.getElementById('winner');

// create racers
function createRacers() {
    track.innerHTML = '';

    // Re-add finish line
    const finishLine = document.createElement('div');
    finishLine.id = 'finish-line';
    track.appendChild(finishLine);

    dots.length = 0;
    lbBody.innerHTML = '';
    rowEls.length = 0;
    lastOrder = [];

    for (let i = 0; i < NUM_RACERS; i++) {
        const img = document.createElement('img');
        img.className = 'racer-dot';
        img.src = '/racer_pictures/racer' + ((i % 12) + 1) + '.png';
        img.alt = 'Horse ' + (i + 1);
        track.appendChild(img);
        dots.push(img);

        // leaderboard rows
        const row = document.createElement('div');
        row.className = 'lb-row';
        row.style.transform = 'translateY(' + (i * ROW_H) + 'px)';
        row.dataset.id = i;
        row.innerHTML = '<div class="lb-rank">' + (i + 1) + '</div><div class="lb-name">Horse ' + (i + 1) + '</div>';
        lbBody.appendChild(row);
        rowEls.push(row);
    }
    lbBody.style.height = (ROW_H * NUM_RACERS) + 'px';
}


// convert pct to oval coordinates
function pctToXY(pct){
    const w = track.clientWidth;
    const h = track.clientHeight;
    const rx = w * 0.42;  // radius x
    const ry = h * 0.45;  // radius y
    const cx = w / 2;
    const cy = h / 2;
    const angle = pct / 100 * 2 * Math.PI;
    const x = cx + rx * Math.cos(angle - Math.PI/2);
    const y = cy + ry * Math.sin(angle - Math.PI/2);
    return {x, y};
}


// update racer positions
function updateTrack(progress){
	for(let i=0;i<dots.length;i++){
		const pct = progress && progress[i]!==undefined?progress[i]:0;
		const pos=pctToXY(pct);
		dots[i].style.left=pos.x+'px';
		dots[i].style.top=pos.y+'px';
	}
	updateLeaderboard(progress);
}

// update leaderboard
function updateLeaderboard(progress){
	const now=performance.now();
	if(now-lastLBUpdate<LB_INTERVAL)return;
	lastLBUpdate=now;
	const rows=[];
	for(let i=0;i<NUM_RACERS;i++){
		const pct = progress && progress[i]!==undefined?progress[i]:0;
		rows.push([i,pct]);
	}
	const orderIndex=new Map();
	if(lastOrder.length===NUM_RACERS){
		lastOrder.forEach((id,idx)=>orderIndex.set(id,idx));
	}
	rows.sort((a,b)=>{if(b[1]!=a[1])return b[1]-a[1];const ao=orderIndex.has(a[0])?orderIndex.get(a[0]):9999;const bo=orderIndex.has(b[0])?orderIndex.get(b[0]):9999;if(ao!=bo)return ao-bo;return a[0]-b[0];});
	for(let rank=0;rank<rows.length;rank++){
		const id=rows[rank][0];
		const row=rowEls[id];
		if(!row)continue;
		row.querySelector('.lb-rank').textContent=(rank+1);
		row.style.transform='translateY('+(rank*ROW_H)+'px)';
	}
	lastOrder=rows.map(r=>r[0]);
}

// display winner
function showWinner(w){
	if(!w){ winnerEl.textContent='Winner: —'; return;}
	let id,timeText="";
	if(typeof w==='object'&&w.id!==undefined){
		id=Number(w.id);
		if(w.finishMs!==undefined){ timeText=" — "+(Number(w.finishMs)/1000).toFixed(3)+"s"; }
	} else id=Number(w);
	if(isNaN(id)){ winnerEl.textContent='Winner: —'; return;}
	winnerEl.textContent='Winner: Horse '+(id+1)+timeText;
}

// WebSocket connection & UI
const proto = location.protocol==='https:'?'wss://':'ws://';
let ws = new WebSocket(proto+location.host+'/ws');
const btnStart = document.getElementById('start');
const statusEl = document.getElementById('status');

ws.onopen=()=>{ statusEl.textContent='connected'; btnStart.disabled=false; };
ws.onclose=()=>{ statusEl.textContent='disconnected'; btnStart.disabled=true; };
ws.onmessage=(e)=>{
	try{
		const msg=JSON.parse(e.data);
		if(msg.type==='progress') updateTrack(msg.data||[]);
		else if(msg.type==='winner') showWinner(msg.data);
	}catch(err){console.error(err);}
};

document.getElementById('homeBtn').onclick = () => {
    window.location.href = '/'; 
};

// Results button: go to results page
document.getElementById('resultsBtn').onclick = () => {
    window.location.href = '/results'; 
};


btnStart.onclick=()=>{
	if(ws.readyState===WebSocket.OPEN){
		const n=parseInt(selectRacers.value);
		NUM_RACERS=n;
		createRacers();
		ws.send(JSON.stringify({type:'START',racers:n}));
		showWinner(null);
	}
};

selectRacers.onchange=()=>{
	NUM_RACERS=parseInt(selectRacers.value);
	createRacers();
	showWinner(null);
	updateLeaderboard(new Array(NUM_RACERS).fill(0));
};
</script>
</body>
</html>`

	return html
}
