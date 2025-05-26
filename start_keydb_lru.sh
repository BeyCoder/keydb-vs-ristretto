#!/bin/bash

set -e

KEYDB_PORT=6379

if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
else
  echo "[ERROR] .env file not found!"
  exit 1
fi

if [ -f "$KEYDB_PID" ]; then
  OLD_PID=$(cat "$KEYDB_PID")
  if ps -p $OLD_PID > /dev/null; then
    echo "[INFO] Killing old KeyDB process (PID $OLD_PID)..."
    kill $OLD_PID
    sleep 1
  fi
  rm -f "$KEYDB_PID"
fi

# === Start KeyDB ===
echo "[INFO] Starting KeyDB on port $KEYDB_PORT..."
keydb-server --port "$KEYDB_PORT" \
             --unixsocket /tmp/keydb.sock \
             --save "" \
             --appendonly no \
             --maxmemory "$KEYDB_MEMORY" \
             --maxmemory-policy "$KEYDB_POLICY" \
             --server-threads $KEYDB_THREADS \
             > "$KEYDB_LOG" 2>&1 &

# === Save PID ===
echo $! > "$KEYDB_PID"
echo "[INFO] KeyDB started with PID $(cat $KEYDB_PID), logs: $KEYDB_LOG"