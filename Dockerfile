# builder image
FROM golang:latest as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=1 GOOS=linux go build -a -tags netgo  -o app


FROM frolvlad/alpine-glibc:alpine-3.13_glibc-2.32
RUN apk add --no-cache ca-certificates tzdata
RUN ln -snf /usr/share/zoneinfo/Europe/Moscow /etc/localtime
RUN mkdir /app
COPY --from=builder /build/app /app/
RUN chmod +x /app/app
CMD [ "/app/app" ]