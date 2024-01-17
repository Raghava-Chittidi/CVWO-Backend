FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod .
RUN go mod tidy 
RUN go mod download

COPY . .
RUN go build -o bin /app/cmd/server/main.go

EXPOSE 8000

ENTRYPOINT [ "/app/bin" ]