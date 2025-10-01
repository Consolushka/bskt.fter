FROM golang:1.23.5-alpine3.21

RUN apk update && apk add --no-cache git && apk add --no-cache bash && apk add build-base &&  apk add cronie

# Change WORKDIR to match your module name
WORKDIR /imp

# Copy the app contents preserving the module structure
COPY ../../app /imp/app
COPY ../../go.mod /imp/
COPY ../../go.sum /imp/

RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN go install github.com/golang/mock/mockgen@v1.6.0

RUN go mod download
RUN go get -d -v ./...
RUN go install -v ./...
RUN cd /imp/app/cmd/scheduler && \
    go build -o /build

EXPOSE 8080

CMD [ "/build" ]
