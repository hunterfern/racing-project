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

.slip-header { cursor: pointer; }

input[type="number"]::-webkit-outer-spin-button,
input[type="number"]::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}
input[type="number"] {
  -moz-appearance: textfield;
}

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
  <span id="balance" style="margin-left:12px;color:#ffd54d;font-weight:700">$100.00</span>
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
// Reset balance on fresh page load (simulating server start)
var RESET_ON_SERVER_START = true; 
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

var linesState = [];           // latest odds lines from server
var balance = 0;               // dollars
var activeBet = null;          // {id, amer, frac, amount}
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
  if (RESET_ON_SERVER_START) {
    balance = 100;      // reset every time the page loads
    saveBalance();      // save to localStorage
    RESET_ON_SERVER_START = false; // only do it once per page load
    renderBalance();
    return;
  }

  var v = window.localStorage.getItem('balance_v1');
  if (!v) { 
      balance = 100; 
      saveBalance(); 
  } else { 
      balance = Math.max(0, Number(v) || 0); 
  }
  renderBalance();
}
function checkGameOver() {
  // only reset if player actually has no money *and* no active bet pending
  if (balance < 1 && !activeBet) {
      alert("Game over! Your money has been reset to $100.");
      balance = 100;
      saveBalance();
      renderBalance();
      // activeBet = null;  // not necessary — we only reset when there's no active bet
      slipMsg.textContent = '';
  }
}

function saveBalance(){
  window.localStorage.setItem('balance_v1', String(balance));
}
function renderBalance(){
  balanceEl.textContent = '$' + balance.toFixed(2);
}
loadBalance();

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
    img.alt = 'Horse ' + (i + 1);
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
  '<div class="lb-name">Horse ' + (i + 1) + '</div>';

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
    winnerEl.textContent = 'Winner: Horse ' + (w.id+1) + ' — ' + secs + 's';
  } else {
    winnerEl.textContent = 'Winner: Horse ' + (Number(w)+1);
  }
  btnStart.textContent = 'Race Again';
}

// WebSocket
const proto = location.protocol === 'https:' ? 'wss://' : 'ws://';
const ws = new WebSocket(proto + location.host + '/ws');

ws.addEventListener('open', () => { statusEl.textContent='connected'; 
  btnStart.disabled=false;
  try { ws.send(JSON.stringify({ type: 'LINES', racers: NUM_RACERS })); } catch (_) {}
});
ws.addEventListener('close', () => { statusEl.textContent='disconnected'; btnStart.disabled=true; });

btnStart.addEventListener('click', () => {
  if (ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify({ type: 'START', racers: NUM_RACERS }));
    raceInProgress = true;
    updateSlipPreview();

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
  try { ws.send(JSON.stringify({ type: 'LINES', racers: NUM_RACERS })); } catch (_) {}
});

function renderOdds(lines){
  if (!Array.isArray(lines) || !oddsEl) return;

  linesState = lines.slice();

  // sort by best price (lowest absolute American odds first)
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
    // Request odds
    try {
        ws.send(JSON.stringify({ type: 'LINES', racers: NUM_RACERS }));
    } catch {}

    // Populate with placeholder horses so user can select one
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


  populateHorseSelect();
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

horseSel.addEventListener('change', updateSlipPreview);
amtInput.addEventListener('input', updateSlipPreview);

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

  // lock bet
  activeBet = { id: id, amer: amer, frac: frac, amount: amount };
  balance -= amount; saveBalance(); renderBalance();
  saveBalance(); 
  renderBalance();
 


  slipMsg.textContent = 'Bet placed: #' + (id+1) + ' for $' + amount +
                        ' at ' + (amer>=0?('+'+amer):amer) + (frac?(' ('+frac+')'):'') + '.';
  placeBtn.disabled = true;
  amtInput.value = '';
  payoutPreview.textContent = 'Payout: —';
});

function settleBet(win){
  if (!activeBet) return; // no bet

  // store a copy because we'll clear activeBet before calling checkGameOver
  const bet = activeBet;

  if (!win || typeof win.id !== 'number'){
    slipMsg.textContent = '';
    activeBet = null;
    // no money changes; still call checkGameOver in case balance is 0 and no active bet
    checkGameOver();
    return;
  }

  if (win.id === bet.id){
    var ret = amerPayout(bet.amer, bet.amount);
    var profit = ret - bet.amount;
    balance += ret; // give stake + profit back
    saveBalance();
    renderBalance();
    slipMsg.textContent = 'WIN! #' + (win.id+1) + ' — returned $' + ret.toFixed(2) + ' (profit $' + profit.toFixed(2) + ').';
  } else {
    // stake was already deducted when bet was placed
    slipMsg.textContent = 'Lost. Winner: #' + (win.id+1) + '. -$' + bet.amount.toFixed(2) + '.';
  }

  // clear the active bet so checkGameOver can act correctly
  activeBet = null;

  // now check whether balance has dropped below 1 and reset if needed
  checkGameOver();
}


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

// navigation
homeBtn.addEventListener('click', () => { btnStart.textContent='Start Race'; window.location.href='/'; });
document.getElementById('resultsBtn').addEventListener('click', ()=>window.location.href='/results');

window.addEventListener('resize', ()=>updateTrack([]));
window.addEventListener('load', () => { 
  createRacers(); 
  selectRacers.value=NUM_RACERS; 
  buildLeaderboardRows(); 
  updateTrack([]); 
  updateLeaderboard(new Array(NUM_RACERS).fill(0));
  try { ws.send(JSON.stringify({ type: 'LINES', racers: NUM_RACERS })); } catch {}
});
</script>
</body>
</html>`

	return html
}
