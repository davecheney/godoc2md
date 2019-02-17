all: examples readme

bin/godoc2md: $(shell find ./ -name \*.go)
	go build -o bin/godoc2md ./cmd/godoc2md/main.go

readme: bin/godoc2md
	bin/godoc2md github.com/WillAbides/godoc2md > README.md

examples: bin/godoc2md
	bin/godoc2md github.com/kr/fs > examples/fs/README.md
	bin/godoc2md github.com/codegangsta/martini > examples/martini/README.md
	bin/godoc2md github.com/gorilla/sessions > examples/sessions/README.md
	bin/godoc2md go/build > examples/build/README.md

.PHONY: examples readme all
