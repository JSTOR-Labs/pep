# PEP Firmware Builder

Collection of tools for building a firmware image for a NUC.

## Requirements

1. A computer or VM running Linux (all testing done on Debian)
2. [Taskfile](https://taskfile.dev/#/installation) installed on the Linux system
3. [Go](https://golang.org/dl/) installed on the Linux system

## Setup

Start by opening a terminal and browse to the directory containing this project. All commands must be run within this project directory.

### Set firmware passwords

To set the passwords for both the standard user account and the root account, edit the variable at the top of `build_in_chroot.sh`.

### Generate SSH keys

1. `ssh-keygen -f pep-tunnel.pem`
2. `ssh-keygen -f jstor-nuc.pem`

The first SSH key generated is used for calling home to a home server.  It makes the most sense to generate a key, then use it for all subsequent builds.

The second SSH key is used for logging into NUCs running this firmware image.

### Configure home server connection

If you are running a server where the NUC can call home, you'll want to edit the `autossh.service` file and `config` file to replace the IP with the correct IP for your call home server.

### Adjust the base image

You'll want to make sure you are building from the most up-to-date base image possible.

1. Open a web browser and check the version for the latest base image [here](https://mirrors.lug.mtu.edu/archlinux/iso/latest/)
2. Record the version number contained in the archlinux-bootstrap-*.tar.gz file you see there.  It should look something like `2021.05.01`
3. Edit build_root.sh to update the ROOT_VERSION variable to match that version you've recorded.

### Provide AWS configuration

The firmware builder and firmware image rely on having AWS credentials available to them for interacting with an S3 bucket to facilitate updates, we need to provide credentials for this.  If you do not wish to use the update service, these files must still exist for the build to complete, but can simply contain dummy data.

1. Within the project root, we'll need to create a .aws directory and populate it with a valid AWS configuration.
```
mkdir .aws
```
2. Create an S3 bucket to be used for updates, along with **read-only** credentials for that bucket.
3. Add those credentials to a credentials file in .aws
```bash
cat <<EOT >> .aws/credentials
[default]
aws_access_key_id = YOURKEYHERE
aws_secret_access_key = YOURSECRETKEYHERE
EOT
```
4. Add a config file with a default region
```bash
cat <<EOT >> .aws/config
[default]
region=YOURREGIONHERE
output=json
```
5. Optionally, update the Taskfile.yml line 45 with your bucket name

### Prepare PEP software

We'll need to build the PEP tools and make them available to the image builder prior to building the image, this section will walk you through each tool and help you build them.

#### Startup Daemon

The instructions are in the pep-linux-startup daemon folder at the root of this repository.  After you build the daemon, copy it to the root of this project.

```bash
cp /path/to/startupdaemon/pep-linux-startup ./
```

#### UpdateD

UpdateD is contained within this repository, and acts as both a tool for assisting in the build process and is embedded into the image to provide update functionality.

First we'll want to browse to the updated project root
```sh
cd updated
```

Next we need to build the binary
```sh
go build -o updated
```

Copy the updated binary to the root of the firmware builder
```sh
cp updated ../
```

Return to the firmware builder project directory
```sh
cd ..
```

#### PEP API

The API is the core of the JSTOR software running within the firmware image.

You can find the instructions for building the PEP API in the project root for that project.  You will need to build both a Linux and a Windows binary.

1. Ensure you cd back to the firmware builder project before continuing

2. Create a folder to contain the binaries
```bash
mkdir pep
```

3. Copy the binaries to the new folder
```bash
cp /path/to/pep-api/api ./pep/
cp /path/to/pep-api/api.exe ./pep/
```

#### PEP App

The app is the frontend app that is loaded when you browse to the base URL of the NUC using a web browser.

Instructions for building the web app can be found in the README in the app project root.  The two steps needed for this are `installing the dependencies` and `generating a static project`

1. Copy the static project to the pep folder
```bash
cp -R /path/to/pep-app/dist ./pep/app
```

### Generate update package signing key
```bash
./updated generate sigkey.priv
```

### Build the Image
```
VERSION=v0.0.1 task 
```

### (OPTIONAL) Deploy package to S3
```
task deploy
```

## Deploy to a NUC

To deploy to a NUC, you'll need a live Linux system to run on the NUC.  

Using the ISO image located [here](https://mirrors.lug.mtu.edu/archlinux/iso/latest/) contains all of the commands needed preinstalled.  You'll want to use a tool like [Rufus](https://rufus.ie/en_US/) (Windows) or [balenaEtcher](https://www.balena.io/etcher/) (MacOS, Windows, Linux) to copy the image to a flash drive.

The following commands should be run in the live image running on the NUC.

1. Format the drive
```bash
parted /dev/sda mklabel gpt
```

2. Create boot partition
```bash
parted /dev/sda mkpart "EFI system partition" fat32 1MiB 261MiB
```

3. Enable ESP on boot partition
```bash
parted /dev/sda set 1 esp on
```

4. Create root partition
```bash
parted /dev/sda "root" ext4 261MiB 100%
```

5. Format boot partition
```bash
mkfs.fat -F32 /dev/sda1
```

6. Format root partition
```bash
mkfs.ext4 /dev/sda2
```

7. Create partition table on data drive
```bash
parted /dev/sdb mklabel gpt
```

8. Create data partition
```bash
parted /dev/sdb mkpart "data" ext4 0% 100%
```

9. Format data partition
```bash
mkfs.ext4 /dev/sdb1
```

10. Set temporary root password on the live image

Feel free to choose your own password here.
```bash
passwd
```

11. Enable the SSH service
```bash
systemctl start sshd
```

12. Find the IP of the NUC
```bash
ip addr
```

13. Mount the destination drive
```bash
mount /dev/sdb1 /mnt
mkdir -p /mnt/boot
mount /dev/sda1 /mnt/boot
```

14. Copy the built firmware to the device (run these commands on the system you built the firmware on)
```bash
scp -r dist/root.img root@IPFROMSTEP12:/mnt/newroot.img
scp -r dist/vmlinuz-linux root@IPFROMSTEP12:/mnt/boot/
scp -r dist/initramfs-linux.img root@IPFROMSTEP12:/mnt/boot/
```

15. Back on the live system, finish the install
```bash
umount -R /mnt
poweroff
```

You now have a fully installed firmware from the one you built.  You can test it out by unplugging your flashdrive used to boot the live system, and power the NUC back on.