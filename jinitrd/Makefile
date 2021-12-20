GOCMD=go
GOBUILD=$(GOCMD) build
OUTDIR=output
TMPDIR=$(OUTDIR)/initramfs
CPIOOUT=initramfs.cpio.gz
BINDIR=$(OUTDIR)/bin
BINARY=jinitrd
LDFLAGS=-ldflags="-s -w"

all: cpio

build:
	mkdir -p $(BINDIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINDIR)/$(BINARY)

clean:
	rm -rf $(OUTDIR)

cpio: build
	mkdir -p $(TMPDIR)
	cp $(BINDIR)/$(BINARY) $(TMPDIR)/init
	pushd $(TMPDIR) && find . -print0 | cpio --null --create --verbose --format=newc | gzip --best > ../$(CPIOOUT) && popd
