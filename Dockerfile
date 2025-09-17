FROM golang:1.25-alpine
WORKDIR /app
COPY . .
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
RUN go mod tidy
RUN go build -o node ./cmd/main.go
ENTRYPOINT ["/entrypoint.sh"]
CMD ["./node"]