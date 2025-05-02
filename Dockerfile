FROM golang:1.24.2

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o real-db ./cmd

EXPOSE 6369/tcp

CMD ["./real-db"]
