#!/bin/bash

# Navigate to the temporary directory
cd /tmp

# Clone the repository
git clone https://github.com/PranaySiddharthBhange/letsgo.git

# Navigate to the cloned directory
cd letsgo 

# Build the Go program
go build -o letsgo main.go

# Move the executable to /bin (requires sudo)
sudo mv letsgo /bin

# Return to the /tmp directory
cd /tmp 

# Remove the cloned directory
rm -rf letsgo
