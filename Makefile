prog=snippetbox

build:
	go build -o $(prog) cmd/web/*


test:
	go test -count=1 -v ./... 

clean:
	rm -f $(prog)