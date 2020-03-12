# builder image
FROM golang:1.14-alpine3.11 as builder
RUN mkdir /build
ADD *.go go.mod go.sum /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -o rates .

# generate clean, final image for end users
FROM alpine:3.11
COPY --from=builder /build/rates .
COPY cronjobs /etc/crontabs/root

# executable
ENTRYPOINT [ "./rates" ]
# arguments that can be overridden
ENV RATES_DATABASE rates:rates@db:3306/rates?parseTime=true
CMD [ "serve" ]