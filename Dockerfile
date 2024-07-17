FROM golang:1.22.5

WORKDIR /usr/src/app

COPY go.mod go.sum ./

COPY . .

RUN go build -o /todo-service cmd/main/main.go

EXPOSE 7540

ENV TODO_PORT=7540
ENV TODO_DBFILE=./scheduler.db
ENV TODO_PASSWORD=123456
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

CMD ["./todo-service"]