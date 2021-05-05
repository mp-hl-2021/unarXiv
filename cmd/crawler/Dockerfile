FROM golang:1.16.2-alpine3.13 as builder
RUN mkdir /build
WORKDIR /build
ADD go.mod /build/
RUN CGO_ENABLED=0 GOOS=linux go mod download
ADD . /build/
RUN CGO_ENABLED=0 GOOS=linux go build -a -o crawler cmd/crawler/main.go

FROM alpine:3.13
COPY --from=builder /build/crawler .

# executable
ENTRYPOINT [ "./crawler" ]
