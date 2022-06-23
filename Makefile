build:
	mkdir -p bin
	go build -o bin/yaml-readme .
copy: build
	cp bin/yaml-readme /usr/local/bin

