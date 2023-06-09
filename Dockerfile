# syntax=docker/dockerfile:1

FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build  -o /leetcrawl cmd/leetcrawl/leetcrawl.go

CMD ["/leetcrawl"]