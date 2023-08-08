FROM golang:1.20.5
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /lemmygo-users
CMD ["/lemmygo-users"]