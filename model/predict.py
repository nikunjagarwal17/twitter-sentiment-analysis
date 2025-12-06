import os
import sys
import argparse
import json
import pickle

os.environ['TF_CPP_MIN_LOG_LEVEL'] = '3'
os.environ['PYTHONUNBUFFERED'] = '1'

import numpy as np
import tensorflow as tf
tf.get_logger().setLevel('ERROR')

from tensorflow.keras.models import load_model
from tensorflow.keras.preprocessing.sequence import pad_sequences


BASE_DIR = os.path.dirname(os.path.abspath(__file__))
ARTIFACT_DIR = os.path.join(BASE_DIR, "artifacts")

MODEL_PATH = os.path.join(ARTIFACT_DIR, "sentiment_lstm.h5")
TOKENIZER_PATH = os.path.join(ARTIFACT_DIR, "tokenizer.pkl")
META_PATH = os.path.join(ARTIFACT_DIR, "meta.npy")


with open(TOKENIZER_PATH, "rb") as f:
    TOKENIZER = pickle.load(f)

META = np.load(META_PATH, allow_pickle=True).item()
MAX_LEN = int(META["max_len"])

MODEL = load_model(MODEL_PATH)


def predict_sentiment(text: str):
    seq = TOKENIZER.texts_to_sequences([text])
    padded = pad_sequences(seq, maxlen=MAX_LEN)
    probs = MODEL.predict(padded, verbose=0)[0]
    label = int(np.argmax(probs))
    confidence = float(probs[label])
    return label, confidence


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--text", required=True)
    args = parser.parse_args()

    try:
        label, confidence = predict_sentiment(args.text)
        result = json.dumps({"label": label, "confidence": confidence})
        sys.stdout.write(result)
        sys.stdout.flush()
    except Exception as e:
        sys.stderr.write(f"Error: {str(e)}\n")
        sys.exit(1)


if __name__ == "__main__":
    main()


