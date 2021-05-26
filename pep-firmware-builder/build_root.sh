#!/bin/bash

# Configuration
MIRROR="https://mirrors.lug.mtu.edu/archlinux"
ROOT_VERSION="2021.01.01"
ROOT_TAR="archlinux-bootstrap-${ROOT_VERSION}-x86_64.tar.gz"
ROOT_URL="${MIRROR}/iso/${ROOT_VERSION}/${ROOT_TAR}"
OUTPUT_IMAGE=${1:-newroot.img}

# END CONFIGURATION - DO NOT EDIT BELOW THIS LINE

# Load logging helpers
source ./logging.sh

# Remove any old root
sudo rm -rf root.x86_64

# Download root base tarball
inf "Downloading bootstrap image..."
curl -O $ROOT_URL &> /dev/null || error_exit "Failed to download bootstrap image!"

# Extra base root
inf "Unzipping bootstrap image..."
sudo tar xf $ROOT_TAR &> /dev/null

# Delete base root tarball
rm $ROOT_TAR 1> /dev/null

# Copy build script into chroot
sudo cp build_in_chroot.sh root.x86_64/
sudo cp logging.sh root.x86_64/

# Copy in custom init_fuctions
# sudo cp init_fuctions root.x86_64/usr/lib/initcpio/init_functions
# sudo cp init_functions root.x86_64/init_functions
sudo cp mkinitcpio.conf root.x86_64/mkinitcpio.conf
sudo cp jstor-nuc.pem.pub root.x86_64/jstor-nuc.pem.pub
sudo cp -R pep root.x86_64/pep
sudo cp pep-api.service root.x86_64/pep-api.service
sudo cp pep-tunnel.pem root.x86_64/pep-tunnel.pem
sudo cp autossh.service root.x86_64/autossh.service
sudo cp config root.x86_64/config
sudo cp fstab root.x86_64/etc/fstab
sudo cp issue root.x86_64/etc/issue
sudo cp pep-linux-startup root.x86_64/usr/local/bin/pep-linux-startup
sudo cp pep-linux-startup.service root.x86_64/lib/systemd/system/pep-linux-startup.service
sudo cp escli root.x86_64/usr/local/bin/escli
sudo cp elasticsearch.yml root.x86_64/elasticsearch.yml
sudo cp -R archjstor/initcpio root.x86_64/initcpio
sudo cp jstorsig root.x86_64/usr/bin/jstorsig
sudo cp updated/updated root.x86_64/usr/bin/updated
sudo cp sigkey.priv.pub root.x86_64/etc/sigkey.priv.pub
sudo cp manifest.json root.x86_64/manifest.json
sudo cp updated.service root.x86_64/lib/systemd/system/updated.service
sudo cp -R .aws root.x86_64/root/.aws


# Fix pacman in chroot
sudo sed -i 's/CheckSpace/#CheckSpace/g' root.x86_64/etc/pacman.conf

# Enter chroot and run script
inf "Configuring the root image..."
sudo root.x86_64/bin/arch-chroot root.x86_64 /build_in_chroot.sh

# Remove chroot build script
sudo rm root.x86_64/build_in_chroot.sh root.x86_64/logging.sh

# Create dist directory
mkdir -p work/{boot,mnt}

# Move latest kernel out of the squash image
sudo mv root.x86_64/boot/* ./work/boot/

sync

# Squash the filesystem
inf "Creating squashed filesystem..."
sudo mksquashfs root.x86_64 work/mnt/$OUTPUT_IMAGE -b 1024k &> /dev/null

# Delete chroot
sudo rm -rf root.x86_64

inf "Created dist/mnt/${OUTPUT_IMAGE}"
