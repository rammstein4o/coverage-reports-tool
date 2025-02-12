ARG IMAGE_GO=golang:1.14-alpine
FROM ${IMAGE_GO}

WORKDIR ${GOPATH}/src/coverage-reports-tool

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -a -installsuffix cgo -o /coverage-reports-tool
