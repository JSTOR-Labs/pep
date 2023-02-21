BINARY_NAME=api
PW_PATH=github.com/JSTOR-Labs/pep/api/cmd.auth_password
PASSWORD=password

mac:
	go mod tidy
	touch ./pdfs/keys/cert.pem
	touch ./pdfs/keys/ciphertext
	go run main.go keys
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -mod vendor -ldflags="-w -s -X '${PW_PATH}=${PASSWORD}'" -o ${BINARY_NAME}
 
windows:
	go mod tidy
	touch ./pdfs/keys/cert.pem
	touch ./pdfs/keys/ciphertext
	go run main.go keys
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -mod vendor -ldflags="-w -s -X '${PW_PATH}=${PASSWORD}'" -o ${BINARY_NAME}.exe

chromebook:
	go mod tidy
	touch ./pdfs/keys/cert.pem
	touch ./pdfs/keys/ciphertext
	go run main.go keys
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -mod vendor -ldflags="-w -s -X '${PW_PATH}=${PASSWORD}'" -o ${BINARY_NAME}

both:
	go mod tidy
	touch ./pdfs/keys/cert.pem
	touch ./pdfs/keys/ciphertext
	go run main.go keys
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -mod vendor -ldflags="-w -s -X '${PW_PATH}=${PASSWORD}'" -o ${BINARY_NAME}
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -mod vendor -ldflags="-w -s -X '${PW_PATH}=${PASSWORD}'" -o ${BINARY_NAME}.exe

clean:
	go clean
	rm ${BINARY_NAME}
	rm ${BINARY_NAME}.exe