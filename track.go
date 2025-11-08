package main

import "fmt"

func generateTrackPageHTML(numRacers int) string {
	if numRacers < 1 {
		numRacers = 1
	}

	html := `<!doctype html>
<html>
<head>
<meta charset="utf-8">
<title>Race Track</title>
<style>
body { font-family: system-ui, sans-serif; text-align: center; padding: 20px; background: #f6f8fb; }
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
</style>
</head>
<body>
<h1>Race Track</h1>
<div id="toolbar">
  <button id="start" disabled>Start Race</button>
  <button id="restartBtn" style="margin-left:8px;">🔁 Start New Race</button>
  <button id="backBtn" style="margin-left:8px;">Home</button>
  <button id="resultsBtn" style="margin-left:8px;">Results</button>
  <span id="status" style="margin-left:8px; color:#555;">connecting…</span>
</div>
<div id="track-container"></div>
<div id="winner"></div>

<script>
const NUM_RACERS = ` + fmt.Sprint(numRacers) + `;
const track = document.getElementById('track-container');
const dots = [];
const initialPositions=[];
(function(){const startPct=0;const gap=3.5;for(let i=0;i<NUM_RACERS;i++){initialPositions.push((startPct-i*gap+100)%100);}})();

// Create image elements
for (let i=0; i<NUM_RACERS; i++){
    const img = document.createElement('img');
    img.className='racer-dot';
    const idx=(i%12)+1;
    img.src='/racer_pictures/racer'+idx+'.png';
    img.alt='Racer '+(i+1);
    track.appendChild(img);
    dots.push(img);
}

const rx=180, ry=100, cx=1, cy=125; // oval
function getAngle(pct){return pct/100*2*Math.PI;}
function updateTrack(progress){
    for(let i=0;i<dots.length;i++){
        const pct=(progress&&progress[i]!==undefined)?progress[i]:(initialPositions[i]||0);
        const angle=getAngle(pct);
        const x=cx+rx*Math.cos(angle-Math.PI/2);
        const y=cy+ry*Math.sin(angle-Math.PI/2);
        dots[i].style.transform="translate("+x+"px, "+y+"px) translate(-50%, -50%)";
    }
}
function showWinner(winnerId){document.getElementById('winner').textContent="Winner: Racer "+(winnerId+1);}
function clearWinner(){document.getElementById('winner').textContent="";}

const proto = location.protocol==='https:'?'wss://':'ws://';
let ws;
let btnStart = document.getElementById('start');
let restartBtn = document.getElementById('restartBtn');
let statusEl = document.getElementById('status');

function connectWS(){
    ws = new WebSocket(proto+location.host+'/ws');

    ws.onopen = () => {
        statusEl.textContent = 'connected';
        btnStart.disabled = false;
    };
    ws.onclose = () => {
        statusEl.textContent = 'disconnected';
        btnStart.disabled = true;
    };

    ws.onmessage = (e) => {
        try{
            const msg=JSON.parse(e.data);
            if(msg.type==='progress') updateTrack(msg.data||[]);
            else if(msg.type==='winner') showWinner(msg.data);
            else if(msg.type==='reset'){ // reset UI if server tells us race reset
                updateTrack([]);
                clearWinner();
                btnStart.disabled = false;
            }
        }catch(err){console.error(err);}
    };
}

connectWS(); // start connection

btnStart.onclick = () => {
    if(ws && ws.readyState===WebSocket.OPEN){
        ws.send(JSON.stringify({type:'START'}));
    }
};

restartBtn.onclick = async () => {
    try{
        restartBtn.disabled = true;
        // clear current race display
        clearWinner();
        updateTrack([]);
        // tell server to reset race
        await fetch('/reset-race', {method:'POST'});
        // give the WS some time to process reset
        setTimeout(()=>{
            if(ws.readyState===WebSocket.CLOSED || ws.readyState===WebSocket.CLOSING){
                connectWS(); // reconnect if disconnected
            }
            restartBtn.disabled = false;
        }, 300);
    }catch(e){
        console.error('Restart failed:', e);
        restartBtn.disabled = false;
    }
};

window.addEventListener('load',()=>updateTrack());

document.getElementById("backBtn").onclick = () => { window.location.href = "/landing"; };
document.getElementById("resultsBtn").onclick = () => { window.location.href = "/results"; };
</script>
</body>
</html>`

	return html
}
