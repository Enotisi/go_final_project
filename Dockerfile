FROM golang:latest

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download

RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o todo-service cmd/main/main.go

EXPOSE 7540

ENV TODO_PORT=7540
ENV TODO_DBFILE=./scheduler.db
ENV TODO_PASSWORD=123456

CMD ["./todo-service"]