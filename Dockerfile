FROM golang:1.13.5-alpine
WORKDIR /app
RUN apk add -U --no-cache ca-certificates
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /tmp/http-echo ./

FROM alpine:3.11
COPY --from=0 /tmp/http-echo /bin/
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 5678
ENTRYPOINT ["/bin/http-echo"]
