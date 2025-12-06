package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
)

type prediction struct {
	Label      int     `json:"label"`
	Confidence float64 `json:"confidence"`
}

var pageTmpl = template.Must(template.New("index").Parse(`
<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>Twitter Sentiment Analysis</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        h1 { font-size: 24px; margin-bottom: 16px; }
        form { margin-bottom: 24px; }
        textarea { width: 100%; height: 120px; padding: 8px; box-sizing: border-box; font-family: inherit; }
        button { padding: 8px 16px; font-size: 14px; cursor: pointer; }
        button:disabled { opacity: 0.6; cursor: not-allowed; }
        .result { margin-top: 16px; padding: 12px; border: 1px solid #ccc; }
        .error { color: #b00000; margin-top: 8px; }
        .loading { color: #666; margin-top: 8px; }
        label { display: block; margin-bottom: 8px; }
    </style>
</head>
<body>
    <h1>Twitter Sentiment Analysis</h1>
    <form method="POST" action="/" id="sentimentForm">
        <label for="text">Tweet text</label>
        <textarea id="text" name="text" required>{{.Input}}</textarea>
        <br/>
        <button type="submit" id="submitBtn">Analyze</button>
    </form>

    <div id="loading" class="loading" style="display: none;">Analyzing...</div>

    {{if .Error}}
    <div class="error" id="errorMsg">{{.Error}}</div>
    {{end}}

    {{if .HasResult}}
    <div class="result" id="result">
        <div><strong>Predicted sentiment:</strong> {{.SentimentText}}</div>
        <div><strong>Label:</strong> {{.Label}}</div>
        <div><strong>Confidence:</strong> {{printf "%.4f" .Confidence}}</div>
    </div>
    {{end}}

    <script>
        const form = document.getElementById('sentimentForm');
        const submitBtn = document.getElementById('submitBtn');
        const loading = document.getElementById('loading');
        const errorMsg = document.getElementById('errorMsg');
        const result = document.getElementById('result');

        form.addEventListener('submit', function(e) {
            submitBtn.disabled = true;
            loading.style.display = 'block';
            if (errorMsg) errorMsg.style.display = 'none';
            if (result) result.style.display = 'none';
        });
    </script>
</body>
</html>
`))

type pageData struct {
	Input         string
	HasResult     bool
	Label         int
	Confidence    float64
	SentimentText string
	Error         string
}

func runPredict(text string) (prediction, error) {
	var outBuf, errBuf bytes.Buffer

	cmd := exec.Command("python", "../model/predict.py", "--text", text)
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()
	if err != nil {
		log.Printf("predict error: %v, stderr: %s", err, errBuf.String())
		return prediction{}, err
	}

	output := outBuf.String()
	
	jsonRegex := regexp.MustCompile(`\{[^{}]*"label"[^{}]*"confidence"[^{}]*\}`)
	jsonMatch := jsonRegex.FindString(output)
	if jsonMatch == "" {
		log.Printf("no JSON found in output: %s", output)
		return prediction{}, errors.New("no valid JSON found in output")
	}

	var p prediction
	if err := json.Unmarshal([]byte(jsonMatch), &p); err != nil {
		log.Printf("json decode error: %v, extracted: %s", err, jsonMatch)
		return prediction{}, err
	}

	return p, nil
}

func sentimentText(label int) string {
	switch label {
	case 0:
		return "Negative"
	case 1:
		return "Positive"
	default:
		return "Unknown"
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		pageTmpl.Execute(w, pageData{})
	case http.MethodPost:
		text := strings.TrimSpace(r.FormValue("text"))
		data := pageData{Input: text}

		if text == "" {
			data.Error = "Please enter some text."
			pageTmpl.Execute(w, data)
			return
		}

		p, err := runPredict(text)
		if err != nil {
			data.Error = "Failed to run prediction. Check server logs."
			pageTmpl.Execute(w, data)
			return
		}

		data.HasResult = true
		data.Label = p.Label
		data.Confidence = p.Confidence
		data.SentimentText = sentimentText(p.Label)
		pageTmpl.Execute(w, data)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Println("Server listening on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}


