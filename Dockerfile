FROM golang:1.15-alpine AS builder
RUN apk add --no-cache git
WORKDIR /tmp/user-api
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o ./out/user-api .


FROM alpine:3.9
RUN apk add ca-certificates

COPY --from=builder /tmp/user-api/out/user-api /user-api

EXPOSE 8080

CMD ["/user-api"]