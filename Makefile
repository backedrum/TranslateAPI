BINARY=translate-api

.DEFAULT_GOAL: $(BINARY)

$(BINARY):
	govendor sync
	go build -o ${BINARY} *.go

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
