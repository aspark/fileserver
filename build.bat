go build -o ./bin/fileserver.exe ./main.go
go build -o ./bin/fileserver-cli.exe ./fileserver-cli/main.go
go build -o ./bin/fileserver-svc.exe ./fileserver-svc/main.go