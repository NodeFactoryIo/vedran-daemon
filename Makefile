COMMIT?=${BUILDCOMMIT}
VERSION?=${BUILDTAG}

# enable cgo because it's required by OSX keychain library
CGO_ENABLED=0

# enable go modules
GO111MODULE=on

export CGO_ENABLED GO111MODULE

dep:
	go get ./...

test:
	go test ./...

lint:
	golangci-lint run

clean:
	rm vedran-daemon 2> /dev/null || exit 0

build:
	go build

install:
	make clean
	make build
	cp vedran-daemon /usr/local/bin

PLATFORMS := linux/amd64 linux/arm windows/amd64 darwin/amd64

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))


$(PLATFORMS):
	@if [ "$(os)" = "windows" ]; then \
			GOOS=$(os) GOARCH=$(arch) go build ${version_flag} -o 'build/windows/vedran-daemon.exe'; \
	else \
			GOOS=$(os) GOARCH=$(arch) go build ${version_flag} -o 'build/${os}-${arch}/vedran-daemon'; \
	fi

buildAll: $(PLATFORMS)
