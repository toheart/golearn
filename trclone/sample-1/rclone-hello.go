package main

/**
@file:
@author: levi.Tang
@time: 2024/6/14 10:42
@description:
**/

import (
	"context"
	"io"
	"os"

	"github.com/rclone/rclone/cmd"
	"github.com/rclone/rclone/fs"
	"github.com/rclone/rclone/fs/operations"
	mutex "sync" // renamed as "sync" already in use
)

var opt = operations.LoggerOpt{}

func main() {
	// Load the configuration file
	//args := []string{"D:\\code", "D:\\code_bak"}
	os.Args = []string{
		"rclone",
		"sync",
		"D:\\code",
		"D:\\code_bak",
		"--dry-run",
		"--combined",
		"sync.log",
	}
	cmd.Main()
}

func GetSyncLoggerOpt(ctx context.Context, fdst fs.Fs, combinedFile string) (operations.LoggerOpt, func(), error) {
	closers := []io.Closer{}

	opt.LoggerFn = syncLoggerFn
	if opt.TimeFormat == "max" {
		opt.TimeFormat = operations.FormatForLSFPrecision(fdst.Precision())
	}
	opt.NewListJSON(ctx, fdst, "")

	open := func(name string, pout *io.Writer) error {
		if name == "" {
			return nil
		}
		if name == "-" {
			*pout = os.Stdout
			return nil
		}
		out, err := os.Create(name)
		if err != nil {
			return err
		}
		*pout = out
		closers = append(closers, out)
		return nil
	}

	if err := open(combinedFile, &opt.Combined); err != nil {
		return opt, nil, err
	}

	close := func() {
		for _, closer := range closers {
			err := closer.Close()
			if err != nil {
				fs.Errorf(nil, "Failed to close report output: %v", err)
			}
		}
	}

	return opt, close, nil
}

var lock mutex.Mutex

func syncLoggerFn(ctx context.Context, sigil operations.Sigil, src, dst fs.DirEntry, err error) {
	lock.Lock()
	defer lock.Unlock()

	if err == fs.ErrorIsDir && !opt.FilesOnly && opt.DestAfter != nil {
		opt.PrintDestAfter(ctx, sigil, src, dst, err)
		return
	}

	_, srcOk := src.(fs.Object)
	_, dstOk := dst.(fs.Object)
	var filename string
	if !srcOk && !dstOk {
		return
	} else if srcOk && !dstOk {
		filename = src.String()
	} else {
		filename = dst.String()
	}

	if sigil.Writer(opt) != nil {
		operations.SyncFprintf(sigil.Writer(opt), "%s\n", filename)
	}
	if opt.Combined != nil {
		operations.SyncFprintf(opt.Combined, "%c %s\n", sigil, filename)
		fs.Debugf(nil, "Sync Logger: %s: %c %s\n", sigil.String(), sigil, filename)
	}
	if opt.DestAfter != nil {
		opt.PrintDestAfter(ctx, sigil, src, dst, err)
	}
}

func anyNotBlank(s ...string) bool {
	for _, x := range s {
		if x != "" {
			return true
		}
	}
	return false
}
