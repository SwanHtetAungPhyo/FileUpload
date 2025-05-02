FROM golang:1.24-alpine
WORKDIR /app
COPY . /app
RUN go mod tidy
RUN go build -o app /app
EXPOSE 8084
CMD ["./app"]