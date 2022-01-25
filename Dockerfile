FROM golang:1.17.6-alpine3.15 as builder
ARG VERSION
WORKDIR /go/src/github.com/xwjdsh/fy
COPY . .
RUN go build ./cmd/fy

FROM alpine:3.15.0
LABEL maintainer="iwendellsun@gmail.com"
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=builder /go/src/github.com/xwjdsh/fy/fy .
ENTRYPOINT ["/fy"]
