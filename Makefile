BINARY_NAME=api
PW_PATH=github.com/JSTOR-Labs/pep/api/cmd.auth_password
ADMIN_PASSWORD=password
PK_PASSWORD=password
API_DIR=./api
FRONTEND_DIR=./app
S3_SRC="source-bucket-name"
S3_DEST="destination-bucket-name"
mac:
	make -C ${API_DIR} mac PK_PASSWORD="${PK_PASSWORD}" ADMIN_PASSWORD="${ADMIN_PASSWORD}" PW_PATH="${PW_PATH}" BINARY_NAME="${BINARY_NAME}"
	mv ${API_DIR}/${BINARY_NAME}-mac ./

windows:
	make -C ${API_DIR} windows PK_PASSWORD="${PK_PASSWORD}" ADMIN_PASSWORD="${ADMIN_PASSWORD}" PW_PATH="${PW_PATH}" BINARY_NAME="${BINARY_NAME}"
	mv ${API_DIR}/${BINARY_NAME}-windows.exe ./

chromebook:
	make -C ${API_DIR} chromebook PK_PASSWORD="${PK_PASSWORD}" ADMIN_PASSWORD="${ADMIN_PASSWORD}" PW_PATH="${PW_PATH}" BINARY_NAME="${BINARY_NAME}"
	mv ${API_DIR}/${BINARY_NAME}-chromebook ./

all-apis:
	make -C ${API_DIR} all PK_PASSWORD="${PK_PASSWORD}" ADMIN_PASSWORD="${ADMIN_PASSWORD}" PW_PATH="${PW_PATH}" BINARY_NAME="${BINARY_NAME}"
	mv ${API_DIR}/${BINARY_NAME}-windows.exe ./
	mv ${API_DIR}/${BINARY_NAME}-chromebook ./
	mv ${API_DIR}/${BINARY_NAME}-mac ./

frontend:
	npm --prefix ${FRONTEND_DIR} run build
	npm --prefix ${FRONTEND_DIR} run generate

all: all-apis frontend

clean-packages:
	rm -rf ./JSTOR-Chromebook
	rm -rf ./JSTOR-Mac
	rm -rf ./JSTOR-Windows

clean-apis:
	make -C ${API_DIR} clean PK_PASSWORD="${PK_PASSWORD}" ADMIN_PASSWORD="${ADMIN_PASSWORD}" PW_PATH="${PW_PATH}" BINARY_NAME="${BINARY_NAME}"
	rm -f ${BINARY_NAME}-chromebook
	rm -f ${BINARY_NAME}-mac
	rm -f ${BINARY_NAME}-windows.exe

clean-frontend:
	rm -rf ${FRONTEND_DIR}/dist

clean-windows:
	rm -rf ./JSTOR-windows
	rm -f ${BINARY_NAME}-windows.exe
	rm -f ${API_DIR}/${BINARY_NAME}-windows.exe
	rm -f ./JSTOR-Windows.zip

clean-chromebook:
	rm -rf ./JSTOR-chromebook
	rm -f ${BINARY_NAME}-chromebook.exe
	rm -f ${API_DIR}/${BINARY_NAME}-chromebook.exe
	rm -f ./JSTOR-Chromebook.zip

clean-mac:
	rm -rf ./JSTOR-mac
	rm -f ${BINARY_NAME}-mac.exe
	rm -f ${API_DIR}/${BINARY_NAME}-mac.exe
	rm -f ./JSTOR-Mac.zip

clean: clean-packages clean-apis clean-frontend clean-mac clean-chromebook clean-windows

download-mac:
	./${BINARY_NAME}-mac download --b ${S3_SRC}

download-chromebook:
	./${BINARY_NAME}-chromebook download --b ${S3_SRC}

assemble-chromebook:
	mkdir JSTOR-chromebook
	mkdir JSTOR-chromebook/JSTOR
	mkdir JSTOR-chromebook/JSTOR/content
	mkdir JSTOR-chromebook/JSTOR/es_config
	mkdir JSTOR-chromebook/JSTOR/es_data
	mkdir JSTOR-chromebook/JSTOR/pdfs

	cp -R ${FRONTEND_DIR}/dist ./JSTOR-chromebook/JSTOR/dist
	cp ./downloads/elasticsearch/chromebook/elasticsearch-7.15.2-amd64.deb ./JSTOR-chromebook/elasticsearch-7.15.2-amd64.deb
	cp ./downloads/es_config/chromebook/elasticsearch.yml ./JSTOR-chromebook/JSTOR/es_config/elasticsearch.yml
	cp -R ./downloads/index/ ./JSTOR-chromebook/JSTOR/es_data
	cp -R ./downloads/pdfs ./JSTOR-chromebook/JSTOR/pdfs
	cp ./shell/start.sh ./JSTOR-chromebook/start.sh
	cp ./install_guides/chromebook/README.pdf ./JSTOR-chromebook/README.pdf
	cp ./${BINARY_NAME}-chromebook ./JSTOR-chromebook/JSTOR/${BINARY_NAME}-chromebook

zip-chromebook:
	zip -r JSTOR-Chromebook.zip ./JSTOR-chromebook/  -x '**/.DS_Store' -x '**/._*' -x '**/__MACOSX'
	aws-vault exec labs --duration=12h -- aws s3 cp ./JSTOR-Chromebook.zip s3://${S3_DEST}/JSTOR-Chromebook.zip

assemble-windows:
	mkdir JSTOR-windows
	mkdir JSTOR-windows/JSTOR
	mkdir JSTOR-windows/JSTOR/content
	mkdir JSTOR-windows/JSTOR/pdfs
	mkdir JSTOR-windows/JSTOR/data

	cp -R ${FRONTEND_DIR}/dist ./JSTOR-windows/JSTOR/dist
	cp -R ./downloads/pdfs ./JSTOR-windows/JSTOR/pdfs

	cp -R ./downloads/elasticsearch/windows/ ./JSTOR-windows/JSTOR
	cp ./downloads/es_config/windows/elasticsearch.yml ./JSTOR-windows/JSTOR/elasticsearch-7.15.2/config/elasticsearch.yml
	cp ./downloads/es_config/windows/myOptions.options ./JSTOR-windows/JSTOR/elasticsearch-7.15.2/config/jvm.options.d/myOptions.options
	rm -rf ./JSTOR-windows/JSTOR/elasticsearch-7.15.2/data
	cp -R ./downloads/index/ ./JSTOR-windows/JSTOR/elasticsearch-7.15.2/data

	mkdir ./JSTOR-windows/JSTOR/elasticsearch-7.15.2/plugins
	
	cp ./shell/start.bat ./JSTOR-windows/start.bat
	cp ./install_guides/windows/README.pdf ./JSTOR-windows/README.pdf
	cp ./${BINARY_NAME}-windows.exe ./JSTOR-windows/JSTOR/${BINARY_NAME}-windows.exe

zip-windows:
	zip -r JSTOR-Windows.zip ./JSTOR-windows/ -x '**/.DS_Store' -x '**/._*' -x '**/__MACOSX'
	aws-vault exec labs --duration=12h -- aws s3 cp ./JSTOR-Windows.zip s3://${S3_DEST}/JSTOR-Windows.zip

assemble-mac:
	mkdir JSTOR-mac
	mkdir JSTOR-mac/JSTOR
	mkdir JSTOR-mac/JSTOR/content
	mkdir JSTOR-mac/JSTOR/pdfs
	mkdir JSTOR-mac/JSTOR/elasticsearch

	cp -R ${FRONTEND_DIR}/dist ./JSTOR-mac/JSTOR/dist
	cp -R ./downloads/pdfs/ ./JSTOR-mac/JSTOR/pdfs

	tar -xvf ./downloads/elasticsearch/mac/elasticsearch-7.15.2-darwin-x86_64.tar -C ./JSTOR-mac/JSTOR/elasticsearch/
	mv ./JSTOR-mac/JSTOR/elasticsearch/elasticsearch-7.15.2/* ./JSTOR-mac/JSTOR/elasticsearch
	rmdir ./JSTOR-mac/JSTOR/elasticsearch/elasticsearch-7.15.2
	mkdir JSTOR-mac/JSTOR/elasticsearch/data
	cp -R ./downloads/index/ ./JSTOR-mac/JSTOR/elasticsearch/data
	cp ./downloads/es_config/mac/elasticsearch.yml ./JSTOR-mac/JSTOR/elasticsearch/config/elasticsearch.yml

	cp ./shell/start.command ./JSTOR-mac/start.command
	cp ./${BINARY_NAME}-mac ./JSTOR-mac/JSTOR/${BINARY_NAME}-mac

zip-mac:
	zip -r JSTOR-Mac.zip ./JSTOR-mac/  -x '**/.DS_Store' -x '**/._*' -x '**/__MACOSX'
	aws-vault exec labs --duration=12h -- aws s3 cp ./JSTOR-Mac.zip s3://${S3_DEST}/JSTOR-Mac.zip

encrypt-on-mac:
	cp ./${BINARY_NAME}-mac ./downloads/${BINARY_NAME}-mac
	./downloads/${BINARY_NAME}-mac encrypt --pw "${PK_PASSWORD}"
	rm -f ./downloads/${BINARY_NAME}-mac

encrypt-on-linux:
	cp ./${BINARY_NAME}-chromebook ./downloads/${BINARY_NAME}-chromebook
	./downloads/${BINARY_NAME}-chromebook encrypt --pw "${PK_PASSWORD}"
	rm -f ./downloads/${BINARY_NAME}-chromebook

timestamp: 
	date +"%FT%T%z"

build-windows: assemble-windows zip-windows clean-windows
build-chromebook: assemble-chromebook zip-chromebook clean-chromebook
build-mac: assemble-mac zip-mac clean-mac

build-windows-on-mac: all encrypt-on-mac assemble-windows zip-windows clean
build-chromebook-on-mac: all encrypt-on-mac assemble-chromebook zip-chromebook clean
build-mac-on-mac: mac frontend encrypt-on-mac assemble-mac zip-mac clean
build-all-on-mac: timestamp all download-mac encrypt-on-mac build-windows build-chromebook build-mac timestamp

build-windows-on-linux: all encrypt-on-linux assemble-windows zip-windows clean
build-chromebook-on-linux: all encrypt-on-linux assemble-chromebook zip-chromebook clean
build-mac-on-linux: all encrypt-on-linux assemble-mac zip-mac clean
build-all-on-linux: timestamp all download-chromebook encrypt-on-linux build-windows build-chromebook build-mac timestamp