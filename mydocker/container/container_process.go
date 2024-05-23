package container

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"mydocker/constant"
	"mydocker/utils"
	"os"
	"os/exec"
	"path"
	"syscall"
)

/**
@file:
@author: levi.Tang
@time: 2024/5/17 21:09
@description:
**/

func NewParentProcess(tty bool, volume string, id string, imageName string) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		log.Errorf("New Pipe err: %s ", err)
		return nil, nil
	}

	cmd := exec.Command("/proc/self/exe", "init")

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		// 将输出重定向
		dirPath := fmt.Sprintf(InfoLocFormat, id)
		if err := os.MkdirAll(dirPath, constant.Perm0622); err != nil {
			log.Errorf("NewParrentProcess mkdir %s, err %s", dirPath, err)
			return nil, nil
		}

		stdLogFilePath := dirPath + GetLogfile(id)
		stdLogFile, err := os.Create(stdLogFilePath)
		if err != nil {
			log.Errorf("NewParentProcess create file %s error %v", stdLogFilePath, err)
			return nil, nil
		}
		cmd.Stdout = stdLogFile
		cmd.Stderr = stdLogFile
	}
	// 将 readPipe 作为 ExtraFiles，这样 cmd 执行时就会外带着这个文件句柄去创建子进程。
	cmd.ExtraFiles = []*os.File{readPipe}
	NewWorkSpace(id, volume, imageName)
	cmd.Dir = utils.GetMerged(id)
	return cmd, writePipe
}

func NewWorkSpace(id string, volume, imageName string) {

	createLower(id, imageName)
	createDirs(id)
	mountOverlayFS(id)

	if volume != "" {
		mntPath := path.Join(id, "merged")
		hostPath, containerPath, err := volumeExtract(volume)
		if err != nil {
			log.Errorf("extract volume failed, maybe volume parameter input is not correct, defailt: %s", err)
			return
		}
		mountVolume(mntPath, hostPath, containerPath)
	}
}

func mountOverlayFS(id string) {
	// 拼接参数
	// e.g. lowerdir=/root/busybox,upperdir=/root/upper,workdir=/root/work
	dirs := utils.GetOverlayFSDirs(utils.GetLower(id), utils.GetUpper(id), utils.GetWorker(id))
	mergePath := utils.GetMerged(id)
	// 完整命令：mount -t overlay overlay -o lowerdir=/root/busybox,upperdir=/root/upper,workdir=/root/work /root/merged
	cmd := exec.Command("mount", "-t", "overlay", "overlay", "-o", dirs, mergePath)
	log.Infof("mount overlayfs: [%s]", cmd.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("%v", err)
	}
}

func createDirs(id string) {
	dirs := []string{
		utils.GetMerged(id),
		utils.GetUpper(id),
		utils.GetWorker(id),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0777); err != nil {
			log.Errorf("mkdirall dir: %s err: %s", dir, err)
		}
	}
}

func createLower(id string, imageName string) {
	// 把busybox作为overlayfs中的lower层
	lowerPath := utils.GetLower(id)
	imagePath := utils.GetImage(imageName)
	log.Infof("lower:%s image.tar:%s", lowerPath, imagePath)
	// 检查是否已经存在busybox文件夹
	exist, err := utils.PathExists(lowerPath)
	if err != nil {
		log.Infof("Fail to judge whether dir %s exists. %v", lowerPath, err)
	}
	// 不存在则创建目录并将busybox.tar解压到busybox文件夹中
	if !exist {
		if err = os.MkdirAll(lowerPath, 0777); err != nil {
			log.Errorf("Mkdir dir %s error. %v", lowerPath, err)
		}
		if _, err = exec.Command("tar", "-xvf", imagePath, "-C", lowerPath).CombinedOutput(); err != nil {
			log.Errorf("Untar dir %s error %v", lowerPath, err)
		}
	}
}

// DeleteWorkSpace Delete the AUFS filesystem while container exit
func DeleteWorkSpace(id string, volume string) {
	log.Infof("delete work space")

	if volume != "" {
		_, containerPath, err := volumeExtract(volume)
		if err != nil {
			log.Errorf("extract volume failed，maybe volume parameter input is not correct，detail:%v", err)
			return
		}
		umountVolume(utils.GetMerged(id), containerPath)
	}
	umountOverlayFS(id)
	deleteDirs(id)
}

func deleteDirs(id string) {
	dirs := []string{
		utils.GetMerged(id),
		utils.GetUpper(id),
		utils.GetWorker(id),
	}

	for _, dir := range dirs {
		if err := os.RemoveAll(dir); err != nil {
			log.Errorf("Remove dir %s error %v", dir, err)
		}
	}
}

func umountOverlayFS(id string) {
	cmd := exec.Command("umount", utils.GetMerged(id))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("%v", err)
	}
}
