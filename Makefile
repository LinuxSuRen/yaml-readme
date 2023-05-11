build:
	mkdir -p bin
	go build -o bin/yaml-readme .
copy: build
	cp bin/yaml-readme /usr/local/bin
test:
	go test ./...
benchmark:
	go test -benchmem -run=^$$ -bench ^BenchmarkGetTopN$$ github.com/linuxsuren/yaml-readme/function