# toll-calculator

```
docker-compose up -d
```

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
```

```
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

```
docker run -p 9090:9090 -v "$(pwd)/.config:/etc/prometheus" prom/prometheus
```

```
for prometeus local usage
prometeus --config.file=<your_config_file>yml
```