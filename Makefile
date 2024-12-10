build_app:
	go build -o enricher cmd/enricher/main.go

test:
	go test ./...

clean:
	rm enricher

