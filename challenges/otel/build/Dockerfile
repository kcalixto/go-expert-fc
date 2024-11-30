FROM golang:1.23 AS build
WORKDIR /app

ARG CMD_PATH=./cmd/a

COPY ${CMD_PATH}/go.mod ${CMD_PATH}/go.sum ./
RUN go mod download

COPY ${CMD_PATH} .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM scratch
WORKDIR /app
COPY --from=build /app/main .
# install required for https requests with ca-certificates
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["./main"]