#!/bin/bash

# Start multiple blockchain nodes on different ports
ports=(5001 5002 5003)

for port in "${ports[@]}"; do
    PORT=$port go run main.go &
    echo "Started node on port $port"
    sleep 2  # Wait between starting nodes
done

echo "All nodes started. Press Ctrl+C to stop all nodes."
wait  # Wait for all background processes 