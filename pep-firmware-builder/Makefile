TARGET := %.tar.gz
SIG := %.sig
WORK_DIR := ./work
DIST_DIR := ./dist
OBJS := $(WORK_DIR)/boot/initramfs-linux.img $(WORK_DIR)/boot/vmlinuz-linux $(WORK_DIR)/mnt/newroot.img
VERSION=v1.0.0

all: root package sign

root: $(OBJS)

manifest.json:
	updated/updated prepare $(VERSION)

prepare: manifest.json

$(OBJS): prepare
	./build_root.sh

$(DIST_DIR)/$(TARGET): $(OBJS)
	updated/updated package $(WORK_DIR)

package: $(DIST_DIR)/$(TARGET)

$(DIST_DIR)/$(SIG): $(DIST_DIR)/$(TARGET)
	updated/updated sign $< sigkey.priv $(DIST_DIR)

sign: $(DIST_DIR)/$(SIG)

clean:
	rm -rf $(WORK_DIR) $(DIST_DIR)

deploy:
	scp -r dist snuc@192.168.1.37:/home/snuc/dist

test: %.sh
	echo $^