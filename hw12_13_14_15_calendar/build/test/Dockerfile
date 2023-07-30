FROM golang:1.19

ENV CODE_DIR /go/src/
WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY .. ${CODE_DIR}

CMD go test -v -race -tags=integration ./integration-tests/...
