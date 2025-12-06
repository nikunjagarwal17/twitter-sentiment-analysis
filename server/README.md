## Server: Twitter Sentiment API

This folder contains a small Go HTTP server that exposes the sentiment model through a simple web page. The server calls the Python predictor in the `model` folder and renders the result.

### Requirements

- Go 1.18+ (tested with recent versions).
- Python 3 with the dependencies from `model/requirements.txt` installed.
- Trained model artifacts present in `model/artifacts/`:
  - `sentiment_lstm.h5`
  - `tokenizer.pkl`
  - `meta.npy`

### Running the server locally

From the project root:

```bash
cd model
pip install -r requirements.txt

# make sure artifacts exist under model/artifacts/
# (train them via the notebook if needed)

cd ../server
go run ./...
```

The server listens on `http://localhost:8080`.

Open that URL in a browser, type a tweet-like text into the textarea, and submit. 