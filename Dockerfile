FROM golang:1.25-alpine
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o node ./cmd/main.go
CMD ["./node"]