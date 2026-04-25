#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cp "$SCRIPT_DIR/pls" /usr/local/bin/pls
chmod +x /usr/local/bin/pls

echo "pls installed successfully"

nohup pls serve >> pls.log 2>&1 &
