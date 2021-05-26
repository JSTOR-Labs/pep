#!/bin/bash

### JSTOR NUC Recovery script

JSTOR_ROOT="/usr/local/share/jstor-dist"

# Partition the primary disk
parted /dev/sda mklabel gpt
parted /dev/sda mkpart primary fat32 0% 512MB
parted /dev/sda mkpart primary ext4 512MB 100%

# Format the partitions
mkfs.vfat -F32 /dev/sda1
mkfs.ext4 /dev/sda2

# Mount the boot device
mount /dev/sda1 /mnt

# Copy the kernel and initial ramdisk
cp ${JSTOR_ROOT}/initramfs-linux.img ${JSTOR_ROOT}/initramfs-linux-fallback.img /mnt/
cp ${JSTOR_ROOT}/vmlinuz-linux /mnt/

# Unmount the boot device
umount /mnt

# Mount the auxillary drive
mount /dev/sdb1 /mnt

# Copy root filesystem over
cp ${JSTOR_ROOT}/root.img /mnt/newroot.img

# Unmount the auxillary drive
umount /mnt

# Configure EFI
efibootmgr -d /dev/sda -p 1 -c -L "JSTOR Linux" -l /vmlinuz-linux -u 'initrd=\initramfs-linux.img'

# Poweroff the system
poweroff