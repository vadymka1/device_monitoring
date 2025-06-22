FROM golang:1.23-alpine

WORKDIR /app
COPY . .

RUN apk add --no-cache protobuf
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest

ENV PATH="/go/bin:$PATH"

RUN go mod tidy
RUN go build -o app .

CMD ["./app"]
