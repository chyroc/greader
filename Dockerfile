FROM golang:1.19 AS build

ENV GOPATH /go
WORKDIR /app

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/greader main.go

RUN strip /go/bin/greader
RUN test -e /go/bin/greader

FROM alpine:latest

LABEL org.opencontainers.image.source=https://github.com/chyroc/greader
LABEL org.opencontainers.image.description="RSS service, providing api similar to google reader."
LABEL org.opencontainers.image.licenses="Apache-2.0"

COPY --from=build /go/bin/greader /bin/greader

CMD /bin/greader start