#!/bin/bash

### CONFIGURATION ###

ROOTPW=""
USERPW=""


### DO NOT EDIT BELOW THIS LINE ###

### DO NOT RUN THIS SCRIPT DIRECTLY ###

# Load logging helpers
source ./logging.sh

# Set the timezone
ln -sf /usr/share/zoneinfo/America/Detroit /etc/localtime

# Set the hostname
echo "jstor" > /etc/hostname

cat > /etc/hosts << END
127.0.0.1   localhost
::1         localhost
127.0.1.1   jstor.localdomain jstor
END

# Set the root password
echo -e "${ROOTPW}\n${ROOTPW}" | (passwd) &> /dev/null

# Add simple network config
cat > /etc/systemd/network/20-ethernet.network << END
[Match]
Name=en*
Name=eth*

[Network]
DHCP=yes
IPv6PrivacyExtensions=yes

[DHCP]
RouteMetric=512
END

# Enable Networking
systemctl enable systemd-networkd &> /dev/null
systemctl enable systemd-resolved &> /dev/null

# Configure pacman
echo "Server = https://mirrors.lug.mtu.edu/archlinux/\$repo/os/\$arch" >> /etc/pacman.d/mirrorlist

# Initialize the pacman keyring
pacman-key --init &> /dev/null
pacman-key --populate archlinux &> /dev/null

# Update the system
pacman --noconfirm -Syu sed &> /dev/null
inf "Performing full system upgrade..."
pacman --noconfirm -Syu base linux linux-firmware &> /dev/null

# Rebuild custom initramfs
# mv /init_functions /usr/lib/initcpio/init_functions
cp /initcpio/install/jstor /usr/lib/initcpio/install/jstor
cp /initcpio/hooks/jstor /usr/lib/initcpio/hooks/jstor
mv /mkinitcpio.conf /etc/mkinitcpio.conf
inf "Generating custom initramfs..."
mkinitcpio -P # &> /dev/null

# Configure the system locale
echo "en_US.UTF-8 UTF-8" >> /etc/locale.gen
locale-gen &> /dev/null
echo "LANG=en_US.UTF-8" > /etc/locale.conf

# Install Elasticsearch
inf "Installing required packages..."
pacman --noconfirm -Syu elasticsearch pacman-contrib openssh sudo autossh nfs-utils parted exfatprogs &> /dev/null

# Initialize Elasticsearch keystore
inf "Creating elasticsearch keystore..."
elasticsearch-keystore create &> /dev/null

# Enable Elasticsearch
systemctl enable elasticsearch &> /dev/null

# Enable the SSH server
systemctl enable sshd &> /dev/null

# Enable updateD
systemctl enable updated &> /dev/null

# Generate the SSH host keys
inf "Generating host SSH keys"
ssh-keygen -A &> /dev/null

# Configure sshd
sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin no/g' /etc/ssh/sshd_config
sed -i 's/#PasswordAuthentication yes/PasswordAuthentication no/g' /etc/ssh/sshd_config

# Add jstor user
useradd -m -G elasticsearch -s /bin/bash jstor

# Set jstor password
echo -e "${USERPW}\n${USERPW}" | (passwd jstor) &> /dev/null

# Create user SSH directory
sudo -u jstor mkdir -p /home/jstor/.ssh
chmod 700 /home/jstor/.ssh

# Add SSH key to authorized keys
cat /jstor-nuc.pem.pub > /home/jstor/.ssh/authorized_keys
chown jstor:jstor /home/jstor/.ssh/authorized_keys
chmod 600 /home/jstor/.ssh/authorized_keys
rm /jstor-nuc.pem.pub

# Copy JSTOR files into place
mv /pep /home/jstor/pep

# Fix up perms
chown jstor:jstor /home/jstor/pep

# Move systemd configuration
mv /pep-api.service /lib/systemd/system/pep-api.service

# Copy elasticsearch.yml into place
mv /elasticsearch.yml /etc/elasticsearch/elasticsearch.yml

# Set permissions
chown -R jstor: /home/jstor/pep
chmod 755 /home/jstor/pep/app
chown -R root:elasticsearch /etc/elasticsearch

# Enable startup actions
systemctl enable pep-linux-startup &> /dev/null

# Enable the JSTOR API
systemctl enable pep-api &> /dev/null

# Setup reverse SSH tunnel
inf "Configuring reverse SSH tunnel"
useradd -m -s /sbin/nologin autotunnel
mv /autossh.service /lib/systemd/system/autossh.service
mkdir -p /home/autotunnel/.ssh
mv /pep-tunnel.pem /home/autotunnel/.ssh/pep-tunnel.pem
mv /config /home/autotunnel/.ssh/config
chown -R autotunnel:autotunnel /home/autotunnel/.ssh
chmod 700 /home/autotunnel/.ssh
chmod 400 /home/autotunnel/.ssh/*
systemctl enable autossh &> /dev/null


inf "Cleaning up the pacman cache..."
# Clean up pacman cache
paccache -rk0 &> /dev/null
