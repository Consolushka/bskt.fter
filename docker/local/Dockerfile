FROM golang:alpine

RUN apk update && apk add --no-cache git && apk add --no-cach bash && apk add build-base

# Change WORKDIR to match your module name
WORKDIR /FTER

# Copy the app contents preserving the module structure
COPY ./app /FTER/app
COPY go.mod /FTER/
COPY go.sum /FTER/
COPY main.go /FTER/
COPY .env /FTER/

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

RUN go mod download
RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -o /build

EXPOSE 8080

CMD [ "/build" ]
