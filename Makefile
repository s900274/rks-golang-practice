export GO111MODULE=on
export CGO_ENABLED=0

all: bld

bld: rks-golang-practice

rks-golang-practice:
	go mod download
	go build -o init/rks-golang-practice ./cmd/rks-golang-practice

clean:
	@rm -f init/rks-golang-practice
	@rm -rf status
	@rm -f  log/*log*
	@rm -rf ./output

cleanlog:
	@rm -f log/*log*
