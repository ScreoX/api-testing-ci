FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY tests ./tests

CMD ["go", "test", "-count=1", "-v", "./..."]
