# JSTOR Labs PEP API

This API is designed to run on both a NUC running a custom Arch Linux image, and Windows PC's off of a flashdrive.

## Prerequisites
1. [Go](https://golang.org/)
2. [Elasticsearch](https://www.elastic.co/start)

## Building


### Binary

Open a terminal and cd into the api root, to set configuration options follow the [Configuration](##Configuration) section after building. The Makefile has options for various standard builds. Building the api binary with the Makefile will also generate a fresh Certificate, private key, and encrypted user password in the `pdfs/keys` directory. Note that this will overwrite any existing content.

The only requirement before building is that you update the value for `ADMIN_PASSWORD` in the Makefile to include the desired admin password. Similarly, be sure to add a value for `PK_PASSWORD` in the Makefile to include the PDF password. Alternatively, both these values can be set in the Make command, as described below. Note that these passwords should be entirely separate, but that both should be memorable.

#### Native

`go build -ldflags="-w -s" -o api`

#### Chromebook

`make chromebook ADMIN_PASSWORD=admin_password PK_PASSWORD=key_password`

#### Windows

`make window ADMIN_PASSWORD=admin_password PK_PASSWORD=key_password`

#### Mac

`make mac ADMIN_PASSWORD=admin_password PK_PASSWORD=key_password`

### Docker

```
docker build -e "admin_password=password123" \ 
    -e "signing_key=asdfjkl12345 \
    -e "elastic_addr=http://localhost:9200" -t pep-api .
```

## PDF Encryption

When the api binary is in the same directory as a directory `pdfs` with unencrypted PDF files, you can run `./api encrypt` to encrypt all pdfs in the directory using the encrypted symmetric key.


## Configuration

Skip this section if you built using docker.  The configuration option are set at build time during a Docker build.

`./api generate -p admin_password -k signing_key -e elastic_addr -s`

Replace `admin_password`, `signing_key`, and `elastic_addr` with your desired settings.

The resulting config file will be located at `$HOME/.labs.toml`

## Running

Once the API is configured, you can run it using `./api serve`