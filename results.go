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
body { font-family: system-ui, sans-serif; padding: 20px; background:#f6f8fb; }
h1 { text-align: center; }
pre { background: #fff; padding: 12px; border: 1px solid #ccc; overflow-x: auto; }
button { margin-bottom: 12px; }
</style>
</head>
<body>
<h1>Race Results</h1>
<button onclick="window.location.href='/'">Back to Progress</button>
<pre>%s</pre>

</body>
</html>`, string(data))
	})
}
