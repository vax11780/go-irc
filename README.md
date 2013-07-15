/client contains a simple text based network client that connects to a hardcoded IP and Port
/server contains a simple text based network multi-client server that receives connections from an arbitrary number of clients and re-broadcasts messages from one client to all other connected clients.  Binds to a hardcoded network TCP Port

go build client.go
go build server.go

On the server:
./server -d


On each client:
./client -d


type /quit to end the client session in the client



TODO:
* Add IP:PORT as arguments
* Add ability to execute LUA or other scripts with a "/command <string>" similiar to /quit


go run client/client.go -s localhost -p 9999
