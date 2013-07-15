#!/bin/sh

mkdir bin

echo "Building client with go"
go build -o bin/client client/client.go

if [ $? -ne 0 ]
then
    echo "Client build failed, exiting"
    exit 1   
fi

echo "Building server with go"
go build -o bin/server server/server.go

if [ $? -ne 0 ]
then
    echo "Server build failed, exiting"
    exit 1   
fi

echo "Build completed successfully"

exit 0
