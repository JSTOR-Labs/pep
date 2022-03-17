GOCMD=go
GOBUILD=$(GOCMD) build
OUTDIR=output
TMPDIR=$(OUTDIR)/initramfs
CPIOOUT=initramfs.img
BINDIR=$(OUTDIR)/bin
BINARY=jinitrd

LDFLAGS=-ldflags="-s -w -X 'main.PersistentVolumeDevPath=$(PERSISTENT_VOLUME)' -X 'main.BootVolumeDevPath=$(BOOT_VOLUME)' -X 'main.InitialFSPath=$(ROOT_VOLUME)'"

all: cpio

$(BINDIR)/$(BINARY):
	mkdir -p $(BINDIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINDIR)/$(BINARY)

clean:
	rm -rf $(OUTDIR)

$(CPIOOUT): $(BINDIR)/$(BINARY)
	mkdir -p $(TMPDIR)
	cp $(BINDIR)/$(BINARY) $(TMPDIR)/init
	cd $(TMPDIR) && find . -mindepth 1 -printf '%P\0' \
		| sort -z | LANG=C bsdtar --uid 0 --gid 0 --null -cnf - -T - \
		| LANG=C bsdtar --null -cf - --format=newc @- | zstd "-T0" > ../$(CPIOOUT)
