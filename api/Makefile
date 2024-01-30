mac:
	# go mod tidy
	# touch ./pdfs/keys/cert.pem
	# touch ./pdfs/keys/ciphertext
	# go run main.go keys --pw '${PK_PASSWORD}'
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -mod vendor -ldflags="-w -s -X '${PW_PATH}=${ADMIN_PASSWORD}'" -o ${BINARY_NAME}-mac
 
windows:
	go mod tidy
	touch ./pdfs/keys/cert.pem
	touch ./pdfs/keys/ciphertext
	go run main.go keys --pw '${PK_PASSWORD}'
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -mod vendor -ldflags="-w -s -X '${PW_PATH}=${ADMIN_PASSWORD}'" -o ${BINARY_NAME}-windows.exe

chromebook:
	go mod tidy
	touch ./pdfs/keys/cert.pem
	touch ./pdfs/keys/ciphertext
	go run main.go keys --pw '${PK_PASSWORD}'
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -mod vendor -ldflags="-w -s -X '${PW_PATH}=${ADMIN_PASSWORD}'" -o ${BINARY_NAME}-chromebook

all:
	go mod tidy
	touch ./pdfs/keys/cert.pem
	touch ./pdfs/keys/ciphertext
	go run main.go keys --pw '${PK_PASSWORD}'
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -mod vendor -ldflags="-w -s -X '${PW_PATH}=${ADMIN_PASSWORD}'" -o ${BINARY_NAME}-chromebook
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -mod vendor -ldflags="-w -s -X '${PW_PATH}=${ADMIN_PASSWORD}'" -o ${BINARY_NAME}-mac
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -mod vendor -ldflags="-w -s -X '${PW_PATH}=${ADMIN_PASSWORD}'" -o ${BINARY_NAME}-windows.exe

clean:
	go clean
	rm -f ${BINARY_NAME}-chromebook
	rm -f ${BINARY_NAME}-mac
	rm -f ${BINARY_NAME}-windows.exe