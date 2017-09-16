BINARY=translate-api

.DEFAULT_GOAL: $(BINARY)

$(BINARY):
	govendor sync
	go test -v ./
	go build -o ${BINARY} *.go

test:
	go test -v ./

format:
	go fmt $$(go list ./... | grep -v /vendor/) ; \
	cd - >/dev/null

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
