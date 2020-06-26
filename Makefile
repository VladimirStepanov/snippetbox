prog=snippetbox

build:
	go build -o $(prog) cmd/web/*


test:
	go test -v ./...

clean:
	rm -f $(prog)