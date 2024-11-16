FROM golang:1.23 AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o cloudrun

FROM scratch
WORKDIR /app
COPY --from=build /app/cloudrun .
# install required for https requests with ca-certificates
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["./cloudrun"]