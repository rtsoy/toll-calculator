# go build -o bin/obu.exe obu/main.go; ./bin/obu.exe

# go build -o bin/receiver.exe ./data_receiver; ./bin/receiver.exe

# go build -o bin/distance_calculator.exe ./distance_calculator; ./bin/distance_calculator.exe

# go build -o bin/aggregator.exe ./aggregator; .\bin\aggregator.exe

obu:
	@go build -o bin/obu obu/main.go
	@./bin/obu

receiver:
	@go build -o bin/receiver ./data_receiver
	@./bin/receiver

calculator:
	@go build -o bin/distance_calculator ./distance_calculator
	@./bin/distance_calculator

calculator:
	@go build -o bin/aggregator ./aggregator
	@./bin/aggregator
