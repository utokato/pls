#!/bin/bash

# 编译前端应用
echo "Building front end app..."
cd web
rm -rf ../offline/dist
rm -rf dist
pnpm build
cd ..
mv web/dist/ offline/dist/
echo "Build front end app success!"

# 编译 go 应用
echo "Building go app..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOFLAGS=-trimpath go build -ldflags="-s -w" -o bin/pls main.go
echo "Build go app success!"