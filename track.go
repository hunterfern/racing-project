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
    height: 15%;           
    top: 3%;
    left: 50%;
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
    z-index: 3;
}


.lane-numbers {
  display: none;


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
    width: 720px;      
    height: 380px;
    margin: 40px auto; 
    background: #88b65f; 
    border-radius: 50% / 45%; 
    border: 2px solid #000; 
    overflow: hidden;
}

/* Dirt oval + inner grass */
#track-container::before,
#track-container::after {
    content: "";
    position: absolute;
    left: 50%;
    top: 50%;
    transform: translate(-50%, -50%);
    border-radius: 50% / 45%;
}

/* Outer dirt lane */
#track-container::before {
    width: 90%;
    height: 96%;
    background: linear-gradient(135deg, #c89b6b, #b27c47);
    box-shadow:
      0 0 0 4px #ffffffcc,
      inset 0 0 0 6px #a36a3b,
      0 4px 12px rgba(0,0,0,0.35);
    z-index: 0;
}

/* Inner grass infield */
#track-container::after {
    width: 50%;
    height: 62%;
    background: radial-gradient(circle at 30% 30%, #7ecf5c, #5f9f3d);
    box-shadow: 
      0 0 0 3px #ffffffcc,
      inset 0 0 0 4px rgba(0,0,0,0.25);
    z-index: 0;
}

.flower-ring {
  position: absolute;
  width: 52%;
  height: 64%;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
  pointer-events: none;
  z-index: 1;
  /* limit to inner-oval-ish area */
  clip-path: ellipse(48% 42% at 50% 50%);
}

.flower-cluster {
  position: absolute;
  width: 0;
  height: 0;
}

.flower-dot {
  position: absolute;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #ff7aa2;
  opacity: 0.8;
  box-shadow: 0 0 3px rgba(0,0,0,0.25);
}


/* colors */
.flower-dot:nth-child(3n) { background: #ffd15c; }
.flower-dot:nth-child(4n) { background: #9bde7e; }

.racer-dot { 
    position: absolute; 
    width: 50px;      
    height: 50px;     
    object-fit: contain; 
    transform: translate(-50%, -50%); 
    transition: left 0.1s linear, top 0.1s linear; 
    pointer-events: none; 
    z-index: 12;
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

#odds .odds-row{
  display: flex;
  justify-content: space-between;
  padding: 6px 8px;
  border-bottom: 1px solid rgba(255,255,255,0.06);
  font-weight: 600;
}
#odds .odds-row:last-child{ border-bottom:none; }

#betSlip{
  position: fixed;
  left: 0; right: 0; bottom: 0;
  display: flex; justify-content: center;
  transform: translateY(calc(100% - 48px));
  transition: transform 300ms ease;
  z-index: 1000;
  pointer-events: auto;
}

#betSlip[aria-hidden="false"]{
  transform: translateY(0%);
}

.slip-paper{
  width: 420px; max-width: 92vw;
  background: #fffef6;
  color: #222;
  border: 1px dashed #222;
  border-radius: 10px;
  box-shadow: 0 10px 32px rgba(0,0,0,.35);
  padding: 14px 14px 16px;
  font-family: ui-sans-serif, system-ui, -apple-system, "Segoe UI", Roboto, Arial;
}

.slip-header{
  display:flex; align-items:center; justify-content:space-between;
  font-weight: 800; letter-spacing: .5px;
  border-bottom: 1px dashed #444; padding-bottom:8px; margin-bottom:10px;
  cursor: pointer;
  user-select: none;
}

.slip-header button{
  border:none; background:transparent; font-size:20px; cursor:pointer;
}

.slip-row{
  display:grid; grid-template-columns: 110px 1fr; gap:10px; align-items:center;
  margin: 10px 0;
}

.slip-row label{ font-weight:700; }

.slip-row select, .slip-row input{
  padding:8px 10px; border:1px solid #bbb; border-radius:6px; font-size:15px;
}

.slip-note{
  font-size: 14px; color:#444; margin:6px 0;
}

.slip-place{
  width:100%; margin-top:10px; padding:10px 12px;
  border:none; border-radius:8px; cursor:pointer;
  background:#111; color:#fff; font-weight:700;
}

.slip-place:disabled{ opacity:.5; cursor:not-allowed; }

.slip-msg{ margin-top:10px; min-height:18px; font-weight:700; }

input[type="number"]::-webkit-outer-spin-button,
input[type="number"]::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}
input[type="number"] {
  -moz-appearance: textfield;
}


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
		func() string {
			s := ""
			for i := 4; i <= 12; i++ {
				s += fmt.Sprintf("<option value=\"%d\">%d</option>", i, i)
			}
			return s
		}() +
		`</select>
  <button id="start" disabled>Start Race</button>
  <button id="homeBtn" style="margin-left:8px;">Home</button>
  <button id="resultsBtn" style="margin-left:8px;">Results</button>
  <span id="balance" style="margin-left:12px;color:#ffd54d;font-weight:700">$100.00</span>
  <span id="status" style="margin-left:8px; color:#555;">connecting…</span>
</div>

<div class="track-wrapper">
  <div class="lane-numbers" id="laneNumbers"></div>
  <div id="track-container">
  <div id="finish-line"></div>
  <div class="flower-ring" id="flowerRing"></div>
</div>
  <div class="leaderboard" id="leaderboard">
    <h2>Leaderboard</h2>
    <div id="lbBody"></div>
  </div>
   <div class="leaderboard" id="oddsPanel">
    <h2>Odds</h2>
    <div id="odds"></div>
  </div>
</div>

<!-- Bet Slip Panel -->
<div id="betSlip" aria-hidden="true">
  <div class="slip-paper">
    <div class="slip-header">
      <div>BET SLIP</div>
      <button id="slipClose" title="Close">&times;</button>
    </div>

    <div class="slip-row">
      <label for="betHorse">Horse</label>
      <select id="betHorse"></select>
    </div>

    <div class="slip-row">
      <label for="betAmount">Amount ($)</label>
      <input id="betAmount" type="number" min="1" step="1" placeholder="Enter stake"/>
    </div>

    <div id="oddsPreview" class="slip-note">Odds: —</div>
    <div id="payoutPreview" class="slip-note">Payout: —</div>

    <button id="placeBet" class="slip-place" disabled>Place Bet</button>

    <div id="slipMsg" class="slip-msg"></div>
  </div>
</div>

<div id="winner">Winner: —</div>

<script>
let NUM_RACERS = ` + fmt.Sprint(defaultNumRacers) + `;
const track = document.getElementById('track-container');
const laneNumbersEl = document.getElementById('laneNumbers');
const dots = [];
const selectRacers = document.getElementById('numRacers');
const lbBody = document.getElementById('lbBody');
const winnerEl = document.getElementById('winner');
const btnStart = document.getElementById('start');
const statusEl = document.getElementById('status');
const oddsEl = document.getElementById('odds');

let lastOrder = [];
const ROW_H = 44;
const rowEls = [];
let lastLBUpdate = 0;
const LB_INTERVAL = 300;


var linesState = [];
var balance = 0;
var activeBet = null;
var raceInProgress = false;


var betSlipEl = document.getElementById('betSlip');
var slipCloseBtn = document.getElementById('slipClose');
var slipHeader = document.querySelector('#betSlip .slip-header');
var horseSel = document.getElementById('betHorse');
var amtInput = document.getElementById('betAmount');
var placeBtn = document.getElementById('placeBet');
var oddsPreview = document.getElementById('oddsPreview');
var payoutPreview = document.getElementById('payoutPreview');
var slipMsg = document.getElementById('slipMsg');
var balanceEl = document.getElementById('balance');

function loadBalance(){
  var v = window.localStorage.getItem('balance_v1');
  if (!v) {
    balance = 100;
    saveBalance();
  } else {
    balance = Math.max(0, Number(v) || 0);
  }
  renderBalance();
}
function saveBalance(){
  window.localStorage.setItem('balance_v1', String(balance));
}
function renderBalance(){
  balanceEl.textContent = '$' + balance.toFixed(2);
}
loadBalance();

function checkGameOver() {
  if (balance < 1) {
    alert("Game over! Your money has been reset to $100.");
    balance = 100;
    saveBalance();
    renderBalance();
    activeBet = null;
    slipMsg.textContent = '';
  }
}

// create racers
function createRacers() {
    const oldDots = track.querySelectorAll('.racer-dot');
    oldDots.forEach(function (el) { el.remove(); });

    laneNumbersEl.innerHTML = '';




    dots.length = 0;
    lbBody.innerHTML = '';
    rowEls.length = 0;
    lastOrder = [];

    for (let i = 0; i < NUM_RACERS; i++) {
        const ln = document.createElement('div');
        ln.className = 'lane-number';
        ln.textContent = '#' + (i + 1);
        laneNumbersEl.appendChild(ln);

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

function generateFlowerClusters() {
  var flowerRing = document.getElementById('flowerRing');
  if (!flowerRing) return;

  // clear any existing clusters
  flowerRing.innerHTML = '';

  var clusterCount = 45;
  var flowersPerCluster = [3, 4, 5];

  for (var i = 0; i < clusterCount; i++) {
    var cluster = document.createElement('div');
    cluster.className = 'flower-cluster';

    var x, y;
    while (true) {
      x = Math.random() * 100;
      y = Math.random() * 100;

      var dx = (x - 50) / 48;
      var dy = (y - 50) / 42;

      if (dx * dx + dy * dy <= 1) break;
    }

    cluster.style.left = x + '%';
    cluster.style.top  = y + '%';

    // scatter flowers around the center
    var count = flowersPerCluster[Math.floor(Math.random() * flowersPerCluster.length)];
    for (var j = 0; j < count; j++) {
      var flower = document.createElement('div');
      flower.className = 'flower-dot';

      flower.style.left = (Math.random() * 18 - 9) + 'px';
      flower.style.top  = (Math.random() * 18 - 9) + 'px';

      cluster.appendChild(flower);
    }

    flowerRing.appendChild(cluster);
  }
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

function renderOdds(lines){
  if (!Array.isArray(lines) || !oddsEl) return;

  linesState = lines.slice();

  var sorted = lines.slice().sort(function(a, b){
    var aa = a.amer >= 0 ? a.amer : Math.abs(a.amer);
    var bb = b.amer >= 0 ? b.amer : Math.abs(b.amer);
    return aa - bb;
  });

  oddsEl.innerHTML = sorted.map(function(l){
    var sign = l.amer >= 0 ? ('+' + l.amer) : ('' + l.amer);
    var frac = (l.fracN && l.fracD) ? (' (' + l.fracN + '/' + l.fracD + ')') : '';
    return '<div class="odds-row">' +
             '<span>#' + (l.id + 1) + ' Horse ' + (l.id + 1) + '</span>' +
             '<span>' + sign + frac + '</span>' +
           '</div>';
  }).join('');

  if (betSlipEl.getAttribute('aria-hidden') === 'false') {
    populateHorseSelect();
    updateSlipPreview();
  }
}

function populateHorseSelect() {
  horseSel.innerHTML = linesState.map(function(l){
    var sign = l.amer >= 0 ? ('+' + l.amer) : ('' + l.amer);
    var frac = (l.fracD && l.fracN) ? (' (' + l.fracN + '/' + l.fracD + ')') : '';
    return '<option value="'+ l.id +'" data-amer="'+ l.amer +'" data-frac="'+ (l.fracD && l.fracN ? (l.fracN + '/' + l.fracD) : '') +'">#'+ (l.id+1) +' Horse '+ (l.id+1) +' — '+ sign + frac +'</option>';
  }).join('');
}

function closeSlip(){ betSlipEl.setAttribute('aria-hidden','true'); }

function openSlip(){
  slipMsg.textContent = '';
  betSlipEl.setAttribute('aria-hidden','false');

  if (!linesState.length) {
    try { ws.send(JSON.stringify({ type: 'LINES', racers: NUM_RACERS })); } catch {}
    var options = [];
    for (var i = 0; i < NUM_RACERS; i++) {
      options.push(
        '<option value="' + i + '" data-amer="0">#' + (i+1) + ' Horse ' + (i+1) + ' — (loading odds…)</option>'
      );
    }
    horseSel.innerHTML = options.join('');
    oddsPreview.textContent = 'Odds: — (loading…)';
    payoutPreview.textContent = 'Payout: —';
    placeBtn.disabled = true;
  } else {
    populateHorseSelect();
  }

  if (horseSel.options.length && horseSel.selectedIndex < 0) {
    horseSel.selectedIndex = 0;
  }
  updateSlipPreview();
}

function toggleSlip(){
  const isOpen = betSlipEl.getAttribute('aria-hidden') === 'false';
  if (isOpen) closeSlip(); else openSlip();
}

slipHeader.addEventListener('click', toggleSlip);
slipCloseBtn.addEventListener('click', closeSlip);

function amerPayout(amer, stake){
  var a = Number(amer);
  var s = Number(stake);
  if (!isFinite(a) || !isFinite(s) || s <= 0) return 0;
  if (a > 0) return s + (s * (a/100));
  return s + (s * (100/Math.abs(a)));
}

function updateSlipPreview(){
  var opt = horseSel.options[horseSel.selectedIndex];
  if (!opt){
    oddsPreview.textContent = 'Odds: —';
    payoutPreview.textContent = 'Payout: —';
    placeBtn.disabled = true;
    return;
  }
  var amer = Number(opt.getAttribute('data-amer'));
  var frac = opt.getAttribute('data-frac');
  var amount = Number(amtInput.value || 0);

  var sign = amer >= 0 ? ('+' + amer) : ('' + amer);
  oddsPreview.textContent = 'Odds: ' + sign + (frac ? (' (' + frac + ')') : '');

  var valid = isFinite(amount) && amount > 0 && amount <= balance;
  var payout = amerPayout(amer, amount);
  payoutPreview.textContent = amount > 0 ? ('Payout: $' + payout.toFixed(2)) : 'Payout: —';
  placeBtn.disabled = !valid;
}

horseSel.addEventListener('change', updateSlipPreview);
amtInput.addEventListener('input', updateSlipPreview);

placeBtn.addEventListener('click', function(){
  var opt = horseSel.options[horseSel.selectedIndex];
  if (!opt) return;
  var amer = Number(opt.getAttribute('data-amer'));
  var frac = opt.getAttribute('data-frac');
  var id = Number(opt.value);
  var amount = Math.floor(Number(amtInput.value || 0));
  if (!isFinite(amount) || amount <= 0) { slipMsg.textContent = 'Enter a valid amount.'; return; }
  if (amount > balance) { slipMsg.textContent = 'Insufficient balance.'; return; }
  if (raceInProgress) { slipMsg.textContent = 'Wait for this race to finish.'; return; }
  if (activeBet) { slipMsg.textContent = 'You already have a bet placed.'; return; }

  activeBet = { id: id, amer: amer, frac: frac, amount: amount };
  balance -= amount;
  saveBalance();
  renderBalance();

  slipMsg.textContent = 'Bet placed: #' + (id+1) + ' for $' + amount +
                        ' at ' + (amer>=0?('+'+amer):amer) + (frac?(' ('+frac+')'):'') + '.';
  placeBtn.disabled = true;
  amtInput.value = '';
  payoutPreview.textContent = 'Payout: —';
});

function settleBet(win){
  if (!activeBet) return;
  if (!win || typeof win.id !== 'number'){ slipMsg.textContent = ''; activeBet = null; return; }

  if (win.id === activeBet.id){
    var ret = amerPayout(activeBet.amer, activeBet.amount);
    var profit = ret - activeBet.amount;
    balance += ret;
    saveBalance();
    renderBalance();
    checkGameOver();
    slipMsg.textContent = 'WIN! #' + (win.id+1) + ' — returned $' + ret.toFixed(2) + ' (profit $' + profit.toFixed(2) + ').';
  } else {
    slipMsg.textContent = 'Lost. Winner: #' + (win.id+1) + '. -$' + activeBet.amount.toFixed(2) + '.';
    checkGameOver();
  }
  activeBet = null;
}

const proto = location.protocol === 'https:' ? 'wss://' : 'ws://';
const ws = new WebSocket(proto + location.host + '/ws');

ws.addEventListener('open', () => {
  statusEl.textContent = 'connected';
  btnStart.disabled = false;
  try { ws.send(JSON.stringify({ type: 'LINES', racers: NUM_RACERS })); } catch (_) {}
});
ws.addEventListener('close', () => {
  statusEl.textContent = 'disconnected';
  btnStart.disabled = true;
});

ws.addEventListener('message', (evt) => {
  let msg;
  try { msg = JSON.parse(evt.data); } catch { return; }

  switch (msg.type) {
    case 'progress':
      updateTrack(msg.data || []);
      updateSlipPreview();
      break;

    case 'winner':
      raceInProgress = false;
      showWinner(msg.data);
      settleBet(msg.data);
      updateSlipPreview();
      btnStart.textContent = 'Race Again';
      btnStart.disabled = false;
      break;

    case 'lines':
      linesState = Array.isArray(msg.data) ? msg.data.slice() : [];
      renderOdds(linesState);
      raceInProgress = false;
      if (betSlipEl.getAttribute('aria-hidden') === 'false') {
        populateHorseSelect();
        if (horseSel.options.length && horseSel.selectedIndex < 0) {
          horseSel.selectedIndex = 0;
        }
        updateSlipPreview();
      }
      break;
  }
});

document.getElementById('homeBtn').onclick = () => {
  btnStart.textContent = 'Start Race';
  window.location.href = '/';
};


document.getElementById('resultsBtn').onclick = () => {
  window.location.href = '/results';
};

btnStart.addEventListener('click', () => {
  if (ws.readyState === WebSocket.OPEN) {
    const n = parseInt(selectRacers.value);
    NUM_RACERS = n;
    createRacers();
    ws.send(JSON.stringify({ type: 'START', racers: n }));
    raceInProgress = true;
    updateSlipPreview();
    showWinner(null);
    btnStart.textContent = 'Racing...';
    btnStart.disabled = true;
    setTimeout(() => { btnStart.disabled = false; }, 1500);
  }
});

selectRacers.addEventListener('change', () => {
  NUM_RACERS = parseInt(selectRacers.value);
  createRacers();
  showWinner(null);
  updateLeaderboard(new Array(NUM_RACERS).fill(0));
  btnStart.textContent = 'Start Race';
  try { ws.send(JSON.stringify({ type: 'LINES', racers: NUM_RACERS })); } catch (_) {}
});

window.addEventListener('load', () => {
  generateFlowerClusters();
  createRacers();
  selectRacers.value = NUM_RACERS;
  updateTrack([]);
  updateLeaderboard(new Array(NUM_RACERS).fill(0));
  try { ws.send(JSON.stringify({ type: 'LINES', racers: NUM_RACERS })); } catch {}
});

window.addEventListener('resize', () => updateTrack([]));
</script>
</body>
</html>`

	return html
}
