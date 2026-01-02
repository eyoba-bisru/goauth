FROM golang:1.25.4-alpine

WORKDIR /app
RUN apk add --no-cache git

# Add this line to disable VCS globally in the container
ENV GOFLAGS="-buildvcs=false"

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go install github.com/air-verse/air@latest

CMD ["air"]