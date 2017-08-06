build:
	GOBIN=${PWD} go install -v .

install:
	install helper /usr/local/bin/
