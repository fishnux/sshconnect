#!/bin/bash
set -e
GOOS=windows GOARCH=amd64 go build -o sshconnect.exe
GOOS=linux GOARCH=amd64 go build -o sshconnect_linux_x64
GOOS=darwin GOARCH=386 go build -o sshconnect_osx
echo Done