build:
	go build -o bin src/main/main.go
	bin/main.exe

compile:
	go build -o bin src/main/main.go

run:
	bin/main.exe