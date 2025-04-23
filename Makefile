run-server:
	go	run	./cmd/server/main.go

run-worker:
	go run ./worker/main.go

build-server:
	go	build	-o	bin/server ./cmd/server

build-worker:
	go	build	-o	bin/worker ./worker