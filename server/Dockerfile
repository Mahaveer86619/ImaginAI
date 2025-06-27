FROM golang:1.24-alpine

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .

RUN go build -o main cmd/main.go

EXPOSE 5050

CMD ["./main"]