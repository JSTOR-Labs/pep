# JSTOR Labs PEP API

This API is designed to run on both a NUC running a custom Arch Linux image, and Windows PC's off of a flashdrive.

## Prerequisites
1. [Go](https://golang.org/)
2. [Elasticsearch](https://www.elastic.co/start)

## Building

### Binary

Open a terminal and cd into the project root, to set configuration options follow the [Configuration](##Configuration) section after building.

#### Native

`go build -ldflags="-w -s" -o api`

#### Linux

`GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o api`

#### Windows

`GOOS=windows GOARCH=amd64 go build -ldflags="-w -s" -o api.exe`

### Docker

```
docker build -e "admin_password=password123" \ 
    -e "signing_key=asdfjkl12345 \
    -e "elastic_addr=http://localhost:9200" -t pep-api .
```

## Configuration

Skip this section if you built using docker.  The configuration option are set at build time during a Docker build.

`./api generate -p admin_password -k signing_key -e elastic_addr -s`

Replace `admin_password`, `signing_key`, and `elastic_addr` with your desired settings.

The resulting config file will be located at `$HOME/.labs.toml`

## Running

Once the API is configured, you can run it using `./api serve`