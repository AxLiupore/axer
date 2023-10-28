package container

import (
	"errors"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func createPath(path string) error {
	if err := os.MkdirAll(path, 0777); err != nil {
		logrus.Errorf("MkdirAll %s error %v", path, err)
	}
	return nil
}

// NewWorkSpace Create an Overlay2 filesystem as container root workspace
func NewWorkSpace(rootPath, volume string) error {
	workerPath := filepath.Join(rootPath, "worker")
	if err := createPath(workerPath); err != nil {
		return err
	}
	err := createLower(rootPath, workerPath)
	if err != nil {
		return err
	}
	err = createDirs(workerPath)
	if err != nil {
		return err
	}
	err = mountOverlayFS(workerPath)
	if err != nil {
		return err
	}
	if volume != "" {
		volumes := volumePathExtract(volume)
		if len(volumes) == 2 && volumes[0] != "" && volumes[1] != "" {
			if err := mountVolume(workerPath, volumes); err != nil {
				return err
			}
		} else {
			logrus.Infof("volume parameter input is not correct.")
			return errors.New("volume parameter input is not correct")
		}
	}
	return nil
}

// createLower use busybox at the lower layer of the overlay filesystem
func createLower(rootPath, workerPath string) error {
	logrus.Infof("create Lowerlay")
	imagePath := filepath.Join(rootPath, "image", "image.tar")
	lowerPath := filepath.Join(workerPath, "lower")
	if err := createPath(lowerPath); err != nil {
		return err
	}
	if _, err := exec.Command("tar", "-xvf", imagePath, "-C", lowerPath+"/").CombinedOutput(); err != nil {
		logrus.Errorf("Untar dir %s error %v", imagePath, err)
		return err
	}
	return nil
}

// createDirs create overlayfs need dirs
func createDirs(workerPath string) error {
	logrus.Infof("create Upper and Work dir")
	upperPath := filepath.Join(workerPath, "upper")
	workPath := filepath.Join(workerPath, "work")
	if err := createPath(upperPath); err != nil {
		return err
	}
	if err := createPath(workPath); err != nil {
		return err
	}
	return nil
}

// mount overlay file system
func mountOverlayFS(workerPath string) error {
	// Create the corresponding mount directory
	mountPath := filepath.Join(workerPath, "container")
	if err := createPath(mountPath); err != nil {
		return nil
	}
	// e.g. lowerdir=/worker/lower,upperdir=/worker/upper,workdir=/worker/work
	lowerDir := filepath.Join(workerPath, "lower")
	upperDir := filepath.Join(workerPath, "upper")
	workDir := filepath.Join(workerPath, "work")
	dirs := "lowerdir=" + lowerDir + ",upperdir=" + upperDir + ",workdir=" + workDir
	// Full command: mount -t overlay overlay -o lowerdir=/worker/lower,upperdir=/worker/upper,workdir=/worker/work /worker/container
	cmd := exec.Command("mount", "-t", "overlay", "overlay", "-o", dirs, mountPath)
	logrus.Infof("mountOverlayFS cmd:%s", cmd.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// Execute the command
	if err := cmd.Run(); err != nil {
		logrus.Errorf("%v", err)
	}
	return nil
}

// DeleteWorkSpace Delete the overlay filesystem while container exit
func DeleteWorkSpace(rootPath, volume string) {
	workerPath := filepath.Join(rootPath, "worker")
	if volume != "" {
		volumes := volumePathExtract(volume)
		length := len(volumes)
		if length == 2 && volumes[0] != "" && volumes[1] != "" {
			umountVolume(workerPath, volumes)
		}
	}
	umountOverlayFS(workerPath)
	deleteDirs(workerPath)
}

func deleteDirs(workerPath string) {
	if err := os.RemoveAll(workerPath); err != nil {
		logrus.Errorf("RemoveAll dir %s error %v", workerPath, err)
	}
}

func umountOverlayFS(workerPath string) {
	mountPath := filepath.Join(workerPath, "container")
	cmd := exec.Command("umount", mountPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Errorf("%v", err)
	}
	if err := os.RemoveAll(mountPath); err != nil {
		logrus.Errorf("Remove dir %s error %v", mountPath, err)
	}
}

// Parsing the volume directory by splitting it with a colon, for example, -v /tmp:/tmp
func volumePathExtract(volume string) []string {
	var volumes []string
	volumes = strings.Split(volume, ":")
	return volumes
}

func mountVolume(workerPath string, volumes []string) error {
	mountPath := filepath.Join(workerPath, "container")
	// The 0th element represents the host machine directory
	parentPath := volumes[0]
	if err := createPath(parentPath); err != nil {
		return nil
	}
	// The 1st element represents the container directory
	containerPath := volumes[1]
	// Concatenate and create the corresponding container directory
	containerVolumePath := filepath.Join(mountPath, containerPath)
	if err := createPath(containerVolumePath); err != nil {
		return err
	}
	cmd := exec.Command("mount", "-o", "bind", parentPath, containerVolumePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Errorf("mount volume failed. %v", err)
		return err
	}
	return nil
}

func umountVolume(workerPath string, volumeURLs []string) {
	mountPath := filepath.Join(workerPath, "container")
	containerPath := filepath.Join(mountPath, volumeURLs[1])
	cmd := exec.Command("umount", containerPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Errorf("Umount volume failed. %v", err)
	}
}
