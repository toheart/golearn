package windows

import (
	"bytes"
	"fmt"
	"github.com/sourcegraph/conc/pool"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

/**
@file:
@author: levi.Tang
@time: 2024/3/29 9:26
@description:
**/

var srcdir = "D:\\Program Files"
var destdir = "D:\\Program Files_bak"

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func TestCreatePath(t *testing.T) {
	// 生产指定长度的路径
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 300)
	var j int
	for i := range b {
		j++
		if j%10 == 0 {
			b[i] = '\\'
			continue
		}
		b[i] = charset[rand.Intn(len(charset))]
	}
	t.Logf("%s \n", b)
	// 创建目录
	mkdest := path.Join(srcdir, string(b))
	if err := os.MkdirAll(mkdest, os.ModePerm); err != nil {
		t.Errorf("%s \n", err.Error())
		return
	}
	t.Logf("mkdir: %s \n", mkdest)
	// 创建文件
	cmd := exec.Command("powershell", "fsutil", "file", "createnew", path.Join(mkdest, "test.txt"), "1073741824")

	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	var sout, serr bytes.Buffer
	cmd.Stdout = &sout
	cmd.Stderr = &serr

	//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
	err := cmd.Run()
	if err != nil {
		fmt.Printf("%s", gbk_to_utf8(serr.Bytes()))
		return
	}
	t.Logf("file: %s , cmd: %s\n", path.Join(mkdest, "test.txt"), gbk_to_utf8(sout.Bytes()))
}

func CopyFile(source string, dest string) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourcefile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}

	}

	return
}

func CopyDir(source string, dest string) (err error) {

	// get properties of source dir
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	// create dest dir

	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(source)

	objects, err := directory.Readdir(-1)

	for _, obj := range objects {

		sourcefilepointer := source + "/" + obj.Name()

		destinationfilepointer := dest + "/" + obj.Name()

		if obj.IsDir() {
			// create sub-directories - recursively
			err = CopyDir(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			// perform copy
			err = CopyFile(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
	return
}

func BenchmarkCopyDir(b *testing.B) {
	for i := 0; i < b.N; i++ {
		os.RemoveAll(destdir)
		CopyDir(srcdir, destdir)
	}
}

func TestCopyDir(t *testing.T) {
	os.RemoveAll(destdir)
	CopyDir(srcdir, destdir)
}

func CopyDirByRoboCopy() {
	cmd := exec.Command("powershell", "robocopy", "\""+srcdir+"\"", "\""+destdir+"\"", "/e", "/mt:12")

	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	var sout, serr bytes.Buffer
	cmd.Stdout = &sout
	cmd.Stderr = &serr

	//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
	err := cmd.Run()
	if err != nil {
		fmt.Printf("%s", gbk_to_utf8(serr.Bytes()))
		return
	}
	fmt.Printf("%s \n ", gbk_to_utf8(sout.Bytes()))
}

func TestCopyDirByRoboCopy(t *testing.T) {
	os.RemoveAll(destdir)
	CopyDirByRoboCopy()
}

func BenchmarkCopyDirByRoboCopy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		os.RemoveAll(destdir)
		CopyDirByRoboCopy()
	}
}

func CopyDirByXcopy(t *testing.T) {
	os.MkdirAll(destdir, os.ModePerm)
	cmd := exec.Command("powershell", "xcopy", "/Y", "/E", "/H", "/R", "\""+srcdir+"\"", "\""+destdir+"\"")

	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	var sout, serr bytes.Buffer
	cmd.Stdout = &sout
	cmd.Stderr = &serr

	//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
	err := cmd.Run()
	if err != nil {
		t.Logf("%s \n", gbk_to_utf8(serr.Bytes()))
		t.Logf("%s \n", gbk_to_utf8(sout.Bytes()))
		return
	}
	t.Logf("%s \n", gbk_to_utf8(sout.Bytes()))
}

func TestCopyDirByXcopy(t *testing.T) {
	os.RemoveAll(destdir)
	CopyDirByXcopy(t)
}

func gbk_to_utf8(s []byte) []byte {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	if d, e := ioutil.ReadAll(reader); e == nil {
		return d
	}
	return s
}

func TestCopyDirWithParalle(t *testing.T) {
	os.RemoveAll(destdir)
	copyDir(srcdir, destdir)
}

func copyDir(sourceDir, destinationDir string) error {
	err := os.MkdirAll(destinationDir, 0755)
	if err != nil {
		return err
	}

	var p = pool.New().WithErrors().WithMaxGoroutines(12)

	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		dest := filepath.Join(destinationDir, relPath)
		if info.IsDir() {
			err := os.MkdirAll(dest, info.Mode())
			if err != nil {
				return err
			}
			return nil
		}

		p.Go(func() error {
			err := copyFile(path, dest)
			return err
		})

		return nil
	})

	return p.Wait()
}

func copyFile(source, destination string) error {
	srcFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return destFile.Sync()
}
