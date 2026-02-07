FROM golang:1.22-alpine

WORKDIR /app

# Install git (needed for some go modules)
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./cmd/api

EXPOSE 8080

CMD ["./app"]
