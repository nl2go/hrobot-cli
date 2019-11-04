PROJECT := hrobot-cli
GO := $(shell which go)

UNIX_EXECUTABLES := \
	linux/386/$(PROJECT)_linux_386 \
	linux/amd64/$(PROJECT)_linux_amd64

COMPRESSED_EXECUTABLES=$(UNIX_EXECUTABLES:%=%.tar.bz2)
COMPRESSED_EXECUTABLE_TARGETS=$(COMPRESSED_EXECUTABLES:%=bin/%)

all: $(EXECUTABLE)

# 386
bin/linux/386/$(PROJECT)_linux_386:
	GOARCH=386 GOOS=linux ${GO} build -ldflags "-extldflags '-static'" -o "$@"

# amd64
bin/linux/amd64/$(PROJECT)_linux_amd64:
	GOARCH=amd64 GOOS=linux ${GO} build -race -ldflags "-extldflags '-static'" -o "$@"

# compress artifacts
%.tar.bz2: %
	tar -jcvf "$<.tar.bz2" "$<"
%.zip: %.exe
	zip "$@" "$<"

release: clean
	$(MAKE) $(UNIX_EXECUTABLES:%=bin/%)

test:
	${GO} test -race ./...

cover:
	${GO} test -race -coverprofile=coverage.out ./...
	${GO} tool cover -func coverage.out

$(PROJECT):
	${GO} build -race -ldflags "-extldflags '-static'" -o "$@"

fmt:
	${GO} fmt ./...

install:
	${GO} install

clean:
	rm -rf bin/

.PHONY: clean release install test