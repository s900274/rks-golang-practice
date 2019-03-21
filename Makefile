GOPATH:=$(CURDIR)/../../../../
export GOPATH

export CGO_ENABLED=0

all: bld

bld: magneto

magneto:
	govendor sync
	go build -o bin/magneto gitlab.kingbay-tech.com/engine-lottery/magneto/cmd/magneto

clean:
	@rm -f init/magneto
	@rm -rf status
	@rm -f  log/*log*
	@rm -rf ./output

cleanlog:
	@rm -f log/*log*
