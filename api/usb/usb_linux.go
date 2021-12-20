package usb

import (
	"bufio"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/JSTOR-Labs/pep/api/elasticsearch"
	"github.com/JSTOR-Labs/pep/api/files"
	"github.com/JSTOR-Labs/pep/api/globals"
	"github.com/JSTOR-Labs/pep/api/utils"
	"github.com/JSTOR-Labs/pep/api/which"
	"github.com/jaypipes/ghw"
	"github.com/vmihailenco/msgpack/v5"
	"gopkg.in/yaml.v2"
)

func FindUSBDrives() []USBDrive {
	drives := make([]USBDrive, 0)
	block, err := ghw.Block()
	if err != nil {
		fmt.Printf("Error gettin block storage info: %v", err)
	}
	for _, disk := range block.Disks {
		log.Debug().Interface("disk", disk).Msg("Found disk")
		if strings.Contains(disk.BusPath, "usb") {
			// we found a usb device
			partitions := make([]Partition, 0)
			for _, part := range disk.Partitions {
				partPath := fmt.Sprintf("/dev/%s", part.Name)

				// Get FS type
				blkid, err := BlkID(partPath)
				if err != nil {
					continue
				}

				partitions = append(partitions, Partition{
					Name:       part.Name,
					Label:      blkid.Label,
					Size:       utils.ByteCountSI(part.SizeBytes),
					Filesystem: blkid.Type,
					SizeKB:     int32(part.SizeBytes / 1024),
				})
			}

			drives = append(drives, USBDrive{
				Name:       disk.Name,
				Size:       utils.ByteCountSI(disk.SizeBytes),
				Model:      disk.Model,
				Partitions: partitions,
			})
		}
	}
	log.Debug().Msgf("Found %d USB drives", len(drives))
	return drives
}

func FormatDrive(name string) error {
	mounted := checkMounted(name)
	if len(mounted) > 0 {
		for _, drive := range mounted {
			log.Debug().Msgf("Unmounting %s", drive)
			err := syscall.Unmount(drive, 0)
			if err != nil {
				log.Error().Err(err).Msg("Unmounting drive")
				return err
			}
		}
	}
	devicePath := fmt.Sprintf("/dev/%s", name)
	cmd := exec.Command(which.LookupExecutable("parted"), devicePath, "mklabel", "gpt", "-s")
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Error().Err(err).Msg("Failed to create partition table")
		return err
	}
	cmd = exec.Command(which.LookupExecutable("parted"), devicePath, "mkpart", "primary", "ntfs", "0%", "100%", "-s")
	err = cmd.Run()
	if err != nil {
		log.Error().Err(err).Msg("Failed to create partition")
		return err
	}
	time.Sleep(time.Second * 5)
	cmd = exec.Command(which.LookupExecutable("mkfs.exfat"), fmt.Sprintf("%s%d", devicePath, 1))
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Error().Err(err).Msg("Failed to create filesystem")
		return err
	}
	cmd.Run()
	return nil
}

// checkMounted checks a drive and its partition to see if it is mounted
func checkMounted(name string) (out []string) {
	out = make([]string, 0)
	drives, err := filepath.Glob(fmt.Sprintf("/dev/%s*", name))
	if err != nil {
		return
	}

	file, err := os.Open("/proc/mounts")
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		for _, drive := range drives {
			if strings.HasPrefix(scanner.Text(), drive+" ") {
				params := strings.Split(scanner.Text(), " ")
				unquotedPath, err := strconv.Unquote("\"" + params[1] + "\"")
				if err != nil {
					log.Error().Err(err).Msg("Failed to unquote path")
					out = append(out, params[1])
				} else {
					out = append(out, unquotedPath)
				}
			}
		}
	}
	return
}

func writeDriveInfo(mountPath string, info driveInfo) error {
	f, err := os.Create(fmt.Sprintf("%s/.driveinfo", mountPath))
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := msgpack.Marshal(info)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	return err
}

func writeFlashdriveConfig(esRoot string, config elasticConfig) error {
	f, err := os.Create(fmt.Sprintf("%sconfig/elasticsearch.yml", esRoot))
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	return err
}

func BuildFlashDrive(name string, snapshotName string, pdfs []string) error {
	globals.BuildJobs[snapshotName] = 0
	drivePath := "/dev/" + name
	mountPoint := "/mnt/" + utils.RandString(16)
	err := os.Mkdir(mountPoint, 0775)
	if err != nil {
		return err
	}
	err = os.Chown(mountPoint, 1000, 1000)
	drive := NewDrive(drivePath, mountPoint)
	err = drive.Mount(false)
	if err != nil {
		log.Error().Err(err).Msg("Failed to mount drive")
		return err
	}

	javaPath, err := files.DownloadFile(files.JavaURL, "OpenJDK11U-jre_x64_windows_hotspot_11.0.10_9.zip")
	if err != nil {
		log.Error().Err(err).Msg("Failed to download java")
		return err
	}

	elasticPath, err := files.DownloadFile(files.ElasticURL, "elasticsearch-7.10.2-no-jdk-windows-x86_64.zip")
	if err != nil {
		log.Error().Err(err).Msg("Failed to download elasticsearch")
		return err
	}
	var javaFiles []string
	if f, err := os.Open(javaPath); err == nil {
		defer f.Close()
		javaFiles, err = files.Untar(mountPoint, f)
		if err != nil {
			log.Error().Err(err).Msg("Failed to untar java")
			return err
		}
	}

	esFiles, err := files.Unzip(elasticPath, mountPoint)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unzip elasticsearch")
		return err
	}

	err = files.CopyFile("/usr/share/pep/pepapi.exe", fmt.Sprintf("%s/api.exe", mountPoint))
	if err != nil {
		return err
	}

	err = files.CopyRecursive(viper.GetString("web.root"), mountPoint)
	if err != nil {
		return err
	}

	pdfLoc := fmt.Sprintf("%s/pdfs", mountPoint)

	// Create pdf directory
	if err := os.Mkdir(pdfLoc, 0775); err != nil && !os.IsExist(err) {
		return err
	}

	// Copy PDFs
	for _, sub := range pdfs {
		if err := files.CopyRecursive(fmt.Sprintf("/mnt/pdfs/%s", sub), pdfLoc); err != nil {
			return err
		}
	}

	// Template out the batch file
	esRoot := strings.ReplaceAll(esFiles[0], "/", "\\")
	javaRoot := strings.ReplaceAll(javaFiles[0], "/", "\\")

	cfg := viper.New()
	cfg.Set("auth.password", viper.GetString("auth.password"))
	cfg.Set("auth.signing_key", viper.GetString("auth.signing_key"))
	cfg.Set("elasticsearch.address", "http://localhost:9201")
	cfg.Set("runtime.flash_drive_mode", true)
	cfg.Set("web.root", "./app")

	cfg.WriteConfigAs(filepath.Join(mountPoint, "api.toml"))

	f, err := os.Open("/usr/share/pep/start.bat")
	if err != nil {
		return err
	}
	defer f.Close()

	t, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	tmpl, err := template.New("start.bat").Parse(string(t))
	if err != nil {
		return err
	}

	output, err := os.Create(fmt.Sprintf("%s/start.bat", mountPoint))
	if err != nil {
		return err
	}
	defer output.Close()

	tmpl.Execute(output, struct {
		ESRoot   string
		JavaHome string
	}{
		ESRoot:   esRoot,
		JavaHome: javaRoot,
	})

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "Unknown JSTOR Appliance"
	}

	err = writeDriveInfo(mountPoint, driveInfo{
		APIPath:  "/api.exe",
		JavaHome: javaRoot,
		ESRoot:   esRoot,
		Created:  time.Now(),
		Hostname: hostname,
	})

	err = writeFlashdriveConfig(fmt.Sprintf("%s/%s", mountPoint, esFiles[0]), elasticConfig{
		Port:                 9201,
		Repo:                 []string{"/mnt/es_backup"},
		XPackMachineLearning: false,
		DiscoveryType:        "single-node",
	})

	if err != nil {
		log.Error().Err(err).Msg("Failed to write Elasticsearch config")
		return err
	}

	go func() {
		defer func() {
			err = drive.Unmount()
			if err != nil {
				log.Error().Err(err).Msg("Failed to unmount drive")
			}
			os.Remove(mountPoint)
		}()
		log.Debug().Msg("Waiting for snapshot to complete")
		status := "IN PROGRESS"
		for status != "SUCCESS" {
			var err error
			status, err = elasticsearch.GetSnapshotStatus(snapshotName)
			if err != nil {
				log.Error().Err(err).Msg("Failed to get snapshot status")
				globals.BuildJobs[snapshotName] = 2
				return
				// cmd.Process.Signal(syscall.SIGTERM)
			}
			time.Sleep(5 * time.Second)
		}
		log.Debug().Msg("Starting flashdrive Elasticsearch")

		u, err := user.Lookup("elasticsearch")
		if err != nil {
			panic("elasticsearch user not found")
		}

		g, err := user.LookupGroup("elasticsearch")
		if err != nil {
			panic("elasticsearch group not found")
		}

		uid, err := strconv.Atoi(u.Uid)
		if err != nil {
			panic("user ID is not a number")
		}

		gid, err := strconv.Atoi(g.Gid)
		if err != nil {
			panic("group ID is not a number")
		}

		log.Printf("Using UID %d and GID %d to start Elasticsearch\n", uid, gid)

		// Pre-create elasticsearch keystore
		cmd := exec.Command(fmt.Sprintf("%s/%sbin/elasticsearch-keystore", mountPoint, esFiles[0]), "create")
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		cmd.SysProcAttr.Credential = &syscall.Credential{
			Uid: uint32(uid),
			Gid: uint32(gid),
		}
		cmd.Dir = fmt.Sprintf("%s/%sbin", mountPoint, esFiles[0])
		cmd.Env = append(os.Environ(), "JAVA_HOME=/usr")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			log.Error().Err(err).Msg("Failed to create elasticsearch keystore")
		}

		// Start elasticsearch
		cmd = exec.Command(fmt.Sprintf("%s/%sbin/elasticsearch", mountPoint, esFiles[0]))
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		cmd.SysProcAttr.Credential = &syscall.Credential{
			Uid: uint32(uid),
			Gid: uint32(gid),
		}
		cmd.Dir = fmt.Sprintf("%s/%sbin", mountPoint, esFiles[0])
		cmd.Env = append(os.Environ(), "JAVA_HOME=/usr")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Start() // Run es in the background
		if err != nil {
			log.Error().Err(err).Msg("Failed to start elasticsearch")
			globals.BuildJobs[snapshotName] = 2
			return
		}

		// Wait for ES to start
		respChan := make(chan int, 1)
		go func() {
			client := http.Client{Timeout: time.Second * 5}
			for {
				resp, err := client.Get("http://localhost:9201/")
				if err == nil && resp.StatusCode == http.StatusOK {
					respChan <- resp.StatusCode
				}
				time.Sleep(time.Second * 1)
			}
		}()
		log.Debug().Msg("Waiting for Elasticsearch to start")
		select {
		case <-respChan:
			break
		case <-time.After(120 * time.Second):
			globals.BuildJobs[snapshotName] = 2
			cmd.Process.Signal(syscall.SIGTERM)
			err := cmd.Wait()
			if err != nil {
				log.Error().Err(err).Msg("Failed to wait for Elasticsearch to exit")
				return
			}
			if cmd.ProcessState.ExitCode() != 0 {
				cmdOut, err := cmd.CombinedOutput()
				if err != nil {
					log.Error().Err(err).Msg("Failed to get output from Elasticsearch")
				} else {
					log.Error().Int("exit_code", cmd.ProcessState.ExitCode()).Msgf("Elasticsearch failed to start: %s", string(cmdOut))
				}
			} else {
				log.Error().Msg("Timed out waiting for Elasticsearch to start")
			}
			return
		}
		// time.Sleep(time.Second * 30)
		log.Debug().Msg("Elasticsearch started, restoring snapshot")
		err = elasticsearch.LoadSnapshot("http://localhost:9201", snapshotName)
		if err != nil {
			log.Error().Err(err).Msg("Failed to load snapshot")
			globals.BuildJobs[snapshotName] = 2
			cmd.Process.Signal(syscall.SIGTERM)
			return
		}

		log.Debug().Msg("Snapshot restore finished, killing elasticsearch")

		cmd.Process.Signal(syscall.SIGTERM)

		log.Debug().Msg("All done.")

		syscall.Sync()

		globals.BuildJobs[snapshotName] = 1
	}()

	return nil
}
