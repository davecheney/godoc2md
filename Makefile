all: examples readme

readme:
	godoc2md github.com/davecheney/godoc2md > README.md

examples:
	godoc2md -ex github.com/kr/fs > examples/fs/README.md
	godoc2md -ex github.com/codegangsta/martini > examples/martini/README.md
	godoc2md -ex github.com/gorilla/sessions > examples/sessions/README.md
	godoc2md -ex go/build > examples/build/README.md
	godoc2md -ex github.com/pkg/errors > examples/errors/README.md

.PHONY: examples readme all
