FROM golang:1.23-bookworm

WORKDIR /app

# Install Python 3.11+ and pip
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    python3 \
    python3-pip \
    python3-venv \
    python3-dev \
    build-essential && \
    rm -rf /var/lib/apt/lists/*

# Create virtual environment and install dependencies
COPY model/requirements.txt ./model/requirements.txt
RUN python3 -m venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"
RUN pip install --upgrade pip && \
    pip install --no-cache-dir -r model/requirements.txt

# Copy application code
COPY model ./model
COPY server ./server

# Build Go application
WORKDIR /app/server
RUN go build -o app .

WORKDIR /app

# Ensure virtual environment is used at runtime
ENV PATH="/opt/venv/bin:$PATH"

EXPOSE 8080

CMD ["./server/app"]


