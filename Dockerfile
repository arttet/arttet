FROM golang:1.15-alpine AS builder

RUN apk add git

COPY /src /src
WORKDIR /src
RUN go get -d -v ./...
RUN go build -o /out/app .

FROM alpine:latest

COPY --from=builder /out/app /app
CMD ["/app"]
