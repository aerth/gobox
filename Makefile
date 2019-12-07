buildflags += -v
buildflags += --ldflags='-w -s'
buildflags += -tags 'netgo osusergo'

gobox: *.go applets/*/*.go pkg/common/*.go
	CGO_ENABLED=0 go build ${buildflags} -o $@

clean:
	rm -f gobox
