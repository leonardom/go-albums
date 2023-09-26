FROM golang:1.19 AS builder

WORKDIR /app

COPY . .

RUN go build -o server

RUN GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o server

FROM scratch

COPY --from=builder /app/server /server

EXPOSE 8000

ENTRYPOINT ["/server"]