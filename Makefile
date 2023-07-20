# go build -o bin/obu.exe obu/main.go; ./bin/obu.exe

# go build -o bin/receiver.exe data_receiver/main.go; ./bin/receiver.exe

obu:
	@go build -o bin/obu obu/main.go
	@./bin/obu

receiver:
	@go build -o bin/receiver ./data_receiver
	@./bin/receiver