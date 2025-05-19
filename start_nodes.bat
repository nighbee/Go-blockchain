@echo off
echo Starting blockchain nodes...

set MINER_HOST=localhost
set PORT=5001
start cmd /k "set MINER_HOST=localhost && set PORT=5001 && go run main.go"
echo Started node on port 5001
timeout /t 2 /nobreak > nul

set PORT=5002
start cmd /k "set MINER_HOST=localhost && set PORT=5002 && go run main.go"
echo Started node on port 5002
timeout /t 2 /nobreak > nul

set PORT=5003
start cmd /k "set MINER_HOST=localhost && set PORT=5003 && go run main.go"
echo Started node on port 5003

echo All nodes started. Close the command windows to stop the nodes. 