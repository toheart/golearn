package container

import (
	"errors"
	errors2 "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

/*
*
@file:
@author: levi.Tang
@time: 2024/5/17 21:09
@description:
*
*/
const fdIndex = 3

func readCommand() []string {

	pipe := os.NewFile(uintptr(fdIndex), "pipe")
	defer pipe.Close()

	msg, err := io.ReadAll(pipe)
	if err != nil {
		log.Errorf("init read pipe error %s", err)
		return nil
	}
	msgStr := string(msg)
	return strings.Split(msgStr, " ")
}

func RunContainerInitProcess() error {
	cmdArray := readCommand()
	if len(cmdArray) == 0 {
		return errors.New("run container get user command error, cmdArray is nil")
	}
	// 挂载文件系统
	setUpMount()
	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		log.Errorf("exec loop path err:%s", err)
		return err
	}
	log.Infof("find path: %s", path)
	if err := syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
		log.Errorf(err.Error())
	}
	return nil
}

/*
*
Init 挂载点
*/

func setUpMount() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Errorf("get current location err: %s", err)
		return
	}
	log.Infof("current location : %s", pwd)
	// systemd 加入linux之后, mount namespace 就变成 shared by default, 所以你必须显示
	// 声明你要这个新的mount namespace独立。
	// 如果不先做 private mount，会导致挂载事件外泄，后续执行 pivotRoot 会出现 invalid argument 错误
	err = syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	err = pivoRoot(pwd)
	if err != nil {
		log.Errorf("pivotRoot failed,detail: %v", err)
		return
	}
	// mount /proc
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	_ = syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	// 由于前面 pivotRoot 切换了 rootfs，因此这里重新 mount 一下 /dev 目录
	// tmpfs 是基于 件系 使用 RAM、swap 分区来存储。
	// 不挂载 /dev，会导致容器内部无法访问和使用许多设备，这可能导致系统无法正常工作
	syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
}

func pivoRoot(root string) error {
	/**
	  NOTE：PivotRoot调用有限制，newRoot和oldRoot不能在同一个文件系统下。
	  因此，为了使当前root的老root和新root不在同一个文件系统下，这里把root重新mount了一次。
	  bind mount是把相同的内容换了一个挂载点的挂载方法
	*/
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return errors2.Wrap(err, "mount rootfs to itself")
	}
	// 创建 rootfs/.pivot_root 目录用于存储 old_root
	pivotDir := filepath.Join(root, ".pivot_root")
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		return err
	}
	// 执行pivot_root调用,将系统rootfs切换到新的rootfs,
	// PivotRoot调用会把 old_root挂载到pivotDir,也就是rootfs/.pivot_root,挂载点现在依然可以在mount命令中看到
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return errors2.WithMessagef(err, "pivotRoot failed,new_root:%v old_put:%v", root, pivotDir)
	}
	// 修改当前的工作目录到根目录
	if err := syscall.Chdir("/"); err != nil {
		return errors2.WithMessage(err, "chdir to / failed")
	}

	// 最后再把old_root umount了，即 umount rootfs/.pivot_root
	// 由于当前已经是在 rootfs 下了，就不能再用上面的rootfs/.pivot_root这个路径了,现在直接用/.pivot_root这个路径即可
	pivotDir = filepath.Join("/", ".pivot_root")
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return errors2.WithMessage(err, "unmount pivot_root dir")
	}
	// 删除临时文件夹
	return os.Remove(pivotDir)
}
