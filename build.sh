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

echo "Build completed successfully, deploying to test box"

scp bin/server ubuntu@ec2-54-227-84-172.compute-1.amazonaws.com

if [ `whoami` = "jenkins" ]
then
    SSH_IDENTITY=~/.ssh/id_rsa
else
    SSH_IDENTITY=~/.ssh/virginia.pem
fi

#ssh -i $SSH_IDENTITY ubuntu@ec2-54-227-84-172.compute-1.amazonaws.com "server -d &"

#if [ $? -ne 0 ]
#then
#    echo "SSH of server to test box failed, exiting"
#    exit 1
#fi

echo "Deployment completed successfully"

exit 0
