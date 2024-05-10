# syntax=docker/dockerfile:1

#FROM golang:1.21 alpine
FROM golang:1.21-alpine

RUN mkdir /app
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./ ./
RUN go build -o boiler-plate ./

# RUN go run ./script/migration/create_migration_script.go

EXPOSE 9004

CMD ["./boiler-plate"]