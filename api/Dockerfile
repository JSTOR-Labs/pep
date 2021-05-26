
############################
# STEP 1 build executable binary
############################

FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates
RUN adduser -D -g '' -h /opt/app appuser
WORKDIR /opt/build
COPY . .
RUN go mod download
RUN go mod verify
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o api

############################
# STEP 2 build a small image
############################

FROM alpine
ARG admin_password=password123
ARG signing_key=asdfjkl12345
ARG elastic_addr=http://localhost:9200
RUN apk add --update libcurl libstdc++ libgcc zlib
WORKDIR /opt/app/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /opt/build/api api
RUN chown -R appuser /opt/app
USER appuser
RUN /opt/app/api generate -p ${admin_password} -k ${signing_key} -e ${elastic_addr} -s
ENTRYPOINT ["/opt/app/api", "serve"]