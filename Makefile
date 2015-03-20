ROOT = /tmp/golang
PROJECT = springmonitor

# First target is the default target
$(PROJECT): clean prepare get test build

clean:
	rm -rf $(ROOT)

prepare:
	mkdir -p $(ROOT)/src/github.com/lukaf/$(PROJECT)
	rsync --delete -auh ./ $(ROOT)/src/github.com/lukaf/$(PROJECT)/

get: prepare
	cd $(ROOT)/src/github.com/lukaf/$(PROJECT)/ && \
		GOPATH=$(ROOT) go get -v

build: get
	cd $(ROOT)/src/github.com/lukaf/$(PROJECT)/ && \
		GOPATH=$(ROOT) go build

test: get
	cd $(ROOT)/src/github.com/lukaf/$(PROJECT)/ && \
		GOPATH=$(ROOT) go test -v

install: build
	cd $(ROOT)/src/github.com/lukaf/$(PROJECT)/ && \
		GOPATH=$(ROOT) go install

.PHONY: clean prepare get build test $(PROJECT)
