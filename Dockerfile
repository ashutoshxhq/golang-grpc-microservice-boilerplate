FROM golang:alpine as builder
RUN mkdir /build 
WORKDIR /build
COPY go.mod ./
COPY . .
RUN go mod tidy
RUN go build -o main .
FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/main /app/
WORKDIR /app
EXPOSE 8000
EXPOSE 9000
CMD ["./main"]