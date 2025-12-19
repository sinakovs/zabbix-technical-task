FROM golang:1.24-bookworm

EXPOSE 8080


WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go mod tidy

COPY . ./

RUN go build -o app ./cmd/server/main.go

CMD ["./app"]