# go build -o bin/obu.exe obu/main.go; ./bin/obu.exe

# go build -o bin/receiver.exe ./data_receiver; ./bin/receiver.exe

# go build -o bin/distance_calculator.exe ./distance_calculator; ./bin/distance_calculator.exe

# go build -o bin/aggregator.exe ./aggregator; .\bin\aggregator.exe

# go build -o bin/gate.exe gateway/main.go; ./bin/gate.exe

obu:
	@go build -o bin/obu obu/main.go
	@./bin/obu

receiver:
	@go build -o bin/receiver ./data_receiver
	@./bin/receiver

calculator:
	@go build -o bin/distance_calculator ./distance_calculator
	@./bin/distance_calculator

agg:
	@go build -o bin/aggregator ./aggregator
	@./bin/aggregator

gate:
	@go build -o bin/gate gateway/main.go
	@./bin/gate

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative types/ptypes.proto
