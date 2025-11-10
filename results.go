package main

import (
	"fmt"
	"net/http"
	"os"
)

// Serve the /results page
func resultsHandler() {
	http.HandleFunc("/results", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store, max-age=0")
		w.Header().Set("Pragma", "no-cache")

		// Read the race log file
		data, err := os.ReadFile("race_log.txt")
		if err != nil {
			fmt.Fprintf(w, "<h1>Race Results</h1><p>Error reading race_log.txt: %v</p>", err)
			return
		}

		// Display the results in <pre> to preserve formatting
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
  background: linear-gradient(135deg, #74ABE2, #3d4de3ff); 
  color: white; /* <— sets default text color */
}

h1 { 
  text-align: center; 
  color: #ffffffff; 
}

pre { 
  background: #516f8dff; 
  padding: 12px; 
  border: 1px solid #000000ff; 
  overflow-x: auto; 
  color: #ffffffff; 
}

button { 
  margin-bottom: 12px; 
  background-color: #ffffffff; 
  color: black; 
  border: none; 
  padding: 8px 14px; 
  border-radius: 8px; 
  cursor: pointer; 
}

button:hover {
  background-color: #ffffffff;
}
</style>
</head>
<body>
<h1>Race Results</h1>
<button onclick="window.location.href='/'">Home</button>
<button onclick="location.href='/track'">Track 1</button>
<button onclick="location.href='/track2'">Track 2</button>
<pre>%s</pre>

</body>
</html>`, string(data))
	})
}
