FROM golang:1.20

WORKDIR /go/src/app

COPY . /go/src/app

RUN go get -d -v ./...
RUN go install -v ./...

# Install Ginkgo and Gomega
RUN go install github.com/onsi/ginkgo/v2/ginkgo
RUN go get github.com/onsi/gomega/...

# Ensure that $GOPATH/bin is in $PATH
ENV PATH $PATH:/go/bin


WORKDIR /go/src/app/e2e