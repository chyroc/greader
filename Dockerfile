FROM golang:1.19 AS build

ENV GOPATH /go
WORKDIR /app

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/greader main.go

RUN strip /go/bin/greader
RUN test -e /go/bin/greader

FROM alpine:latest

COPY --from=build /go/bin/greader /bin/greader

CMD ["/bin/greader"]