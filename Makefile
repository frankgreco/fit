ALL_SRC := $(shell find . -name "*.go" | grep -v -e vendor \
        -e ".*/\..*" \
        -e ".*/_.*" \
        -e ".*/mocks.*")
		
GOFMT=gofmt

.PHONY: fmt
fmt:
	$(GOFMT) -e -s -l -w $(ALL_SRC)