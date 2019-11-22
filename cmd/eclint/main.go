package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/editorconfig/editorconfig-core-go/v2"
	"github.com/mattn/go-colorable"
	"gitlab.com/greut/eclint"
	"k8s.io/klog/v2"
	"k8s.io/klog/v2/klogr"
)

var (
	version = "dev"
)

func main() { //nolint:funlen
	flagVersion := false
	log := klogr.New()
	opt := eclint.Option{
		Stdout:            os.Stdout,
		ShowErrorQuantity: 10,
		IsTerminal:        terminal.IsTerminal(syscall.Stdout),
		Log:               log,
	}

	if runtime.GOOS == "windows" {
		opt.Stdout = colorable.NewColorableStdout()
	}

	// Flags
	klog.InitFlags(nil)
	flag.BoolVar(&flagVersion, "version", false, "print the version number")
	flag.BoolVar(&opt.NoColors, "no_colors", false, "enable or disable colors")
	flag.BoolVar(&opt.Summary, "summary", false, "enable the summary view")
	flag.BoolVar(
		&opt.ShowAllErrors,
		"show_all_errors",
		false,
		fmt.Sprintf("display all errors for each file (otherwise %d are kept)", opt.ShowErrorQuantity),
	)
	flag.StringVar(&opt.Exclude, "exclude", "", "paths to exclude")
	flag.Parse()

	if flagVersion {
		fmt.Fprintf(opt.Stdout, "eclint %s\n", version)
		return
	}

	args := flag.Args()
	files, err := eclint.ListFiles(log, args...)
	if err != nil {
		log.Error(err, "error while handling the arguments")
		flag.Usage()
		os.Exit(1)
		return
	}

	log.V(1).Info("files", "count", len(files), "exclude", opt.Exclude)

	if opt.Summary {
		opt.ShowAllErrors = true
		opt.ShowErrorQuantity = int(^uint(0) >> 1)
	}

	c := 0
	for _, filename := range files {
		// Skip excluded files
		if opt.Exclude != "" {
			ok, err := editorconfig.FnmatchCase(opt.Exclude, filename)
			if err != nil {
				log.Error(err, "exclude pattern failure", "exclude", opt.Exclude)
				c++
				break
			}
			if ok {
				continue
			}
		}

		errs := eclint.Lint(filename, opt.Log)
		c += len(errs)
		err := eclint.PrintErrors(opt, filename, errs)
		if err != nil {
			log.Error(err, "print errors failure", "filename", filename)
		}
	}

	if c > 0 {
		opt.Log.V(1).Info("Some errors were found.", "count", c)
		os.Exit(1)
	}
}