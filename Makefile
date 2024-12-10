build_app:
	go build -o enricher_app cmd/enricher/main.go

test:
	go test ./...

clean:
	rm enricher

