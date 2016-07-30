export GOPATH=`pwd`
go build -o bin/echo-server ./src/echo-server
go build -o bin/bolt-server ./src/bolt-server
go build -o bin/storm-server ./src/storm-server
