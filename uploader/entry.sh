#!/bin/sh
set -e

shutdown() {
    echo "Shutting down..."
    exit 1
}

# Start the Node.js app in the background
cd /app/
node index.js &
# Start the Go app in the foreground
./main || shutdown