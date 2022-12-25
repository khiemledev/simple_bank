FROM golang:1.19-alpine3.16 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Running stage
FROM alpine:3.16
WORKDIR /app
COPY ./db/migration /app/db/migration
COPY --from=builder /app/main .
COPY ./app.env .
COPY ./start.sh .
COPY ./wait-for.sh .

EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT [ "/app/start.sh" ]
