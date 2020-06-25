prog=snippetbox

build:
	go build -o $(prog) cmd/web/*


test:
	go test -v cmd/web/*

clean:
	rm -f $(prog)