## Model: Twitter Sentiment Analysis

This folder contains the end-to-end sentiment analysis model: data preparation, preprocessing, LSTM training, artifact saving, and a small prediction script used by the Go backend.

### Contents

- `Dataset analysis.ipynb`: basic exploration of the raw dataset.
- `Twitter_Sentiment_End_to_End.ipynb`: unified notebook that loads data, preprocesses tweets, trains the LSTM, evaluates it, saves the model, and demonstrates a sample prediction.
- `train-test-split.py`, `preprocessing.py`, `lstm.py`, `ml_tfidf.py`: original script-based workflow (no changes required if you prefer this path).
- `predict.py`: small CLI script that loads the saved model and tokenizer and returns a JSON prediction for a single input text.
- `data/dataset.csv`: raw Sentiment140 data (not tracked in Git).
- `artifacts/`: saved model and tokenizer (created after training, not tracked in Git).

### Setup

From the project root:

```bash
cd model
pip install -r requirements.txt
```

Place the Sentiment140 CSV at `model/data/dataset.csv`. The CSV should match the original Kaggle format.

### Training and saving artifacts

The recommended way to train and export artifacts is via the notebook:

1. Start Jupyter:
   ```bash
   jupyter notebook
   ```
2. Open `Twitter_Sentiment_End_to_End.ipynb`.
3. Run the cells in order. At the end, the notebook will save:
   - `artifacts/sentiment_lstm.h5`
   - `artifacts/tokenizer.pkl`
   - `artifacts/meta.npy`

These files are required by both `predict.py` and the Go backend.

If you prefer scripts, you can still run:

```bash
python train-test-split.py
python preprocessing.py
python lstm.py
```

but the Go backend is wired to use the artifacts produced by the notebook or any equivalent process that writes the same files.

### Using the predictor from the command line

After training and saving artifacts, you can get a sentiment prediction for an arbitrary text:

```bash
cd model
python predict.py --text "I love this project"
```

The script prints a JSON object:

```json
{"label": 1, "confidence": 0.93}
```

`label` follows the original Sentiment140 convention (`0` negative, `1` positive). `confidence` is the model probability for the predicted class.

