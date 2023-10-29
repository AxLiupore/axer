package container

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// RunContainerInitProcess into container to process this func
func RunContainerInitProcess() error {
	//  get data from pipe
	cmdArray := readUserCommand()
	if cmdArray == nil || len(cmdArray) == 0 {
		return fmt.Errorf("run container get user command error, cmdArray is nil")
	}

	// mount file system
	err := setMount()
	if err != nil {
		return err
	}

	// find cmd absolute path
	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		logrus.Errorf("Exec loop path error %v", err)
		return err
	}

	if err := syscall.Exec(path, cmdArray, os.Environ()); err != nil {
		logrus.Errorf(err.Error())
	}
	return nil
}

// Read parameters passed from parent process
func readUserCommand() []string {
	pipe := os.NewFile(uintptr(3), "pipe")
	defer func(pipe *os.File) {
		_ = pipe.Close()
	}(pipe)
	msg, err := io.ReadAll(pipe)
	if err != nil {
		logrus.Errorf("init read pipe error %v", err)
		return nil
	}
	msgStr := string(msg)
	return strings.Split(msgStr, " ")
}

// Remount the file system in the container process
func setMount() error {
	err := syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	if err != nil {
		return err
	}

	// Get the current working directory, which is set by the parent process itself
	pwd, err := os.Getwd()
	if err != nil {
		logrus.Errorf("Get current location error %v", err)
		return err
	}

	err = pivotRoot(pwd)
	if err != nil {
		return err
	}

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	err = syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	if err != nil {
		return err
	}
	err = syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
	if err != nil {
		return nil
	}
	return nil
}

// pivotRoot function switches the root file system of the current process to the specified path
func pivotRoot(root string) error {
	// Remount the specified file system onto itself to ensure 'root' and the new root are not on the same file system
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		logrus.Errorf("mount rootfs to itself error: %v", err)
		return fmt.Errorf("mount rootfs to itself error: %v", err)
	}

	// Create a temporary directory '.pivot_root' in the new root directory to store the old root
	pivotDir := filepath.Join(root, ".pivot_root")
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		return err
	}

	// Use PivotRoot to switch the current process's root file system to the new root file system
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("pivot_root %v", err)
	}

	// Change the current working directory to the root directory
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir / %v", err)
	}

	// Set the new mount point for unmounting the old root file system later
	pivotDir = filepath.Join("/", ".pivot_root")

	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount pivot_root dir %v", err)
	}

	// remove temp dir
	return os.Remove(pivotDir)
}
