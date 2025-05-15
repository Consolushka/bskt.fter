FROM golang:alpine

RUN apk update && apk add --no-cache git && apk add --no-cach bash && apk add build-base

# Change WORKDIR to match your module name
WORKDIR /IMP

# Copy the app contents preserving the module structure
COPY ./app /IMP/app
COPY go.mod /IMP/
COPY go.sum /IMP/
COPY main.go /IMP/
COPY .env /IMP/

RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN go install github.com/golang/mock/mockgen@v1.6.0

RUN go mod download
RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -o /build

EXPOSE 8080

CMD [ "/build" ]
