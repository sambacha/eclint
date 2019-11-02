package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/editorconfig/editorconfig-core-go/v2"
	"github.com/go-logr/logr"
	"k8s.io/klog/v2"
	"k8s.io/klog/v2/klogr"
)

var (
	version = "dev"
	log     logr.Logger
)

func walk(paths ...string) ([]string, error) {
	files := make([]string, 0)
	for _, path := range paths {
		err := filepath.Walk(path, func(p string, i os.FileInfo, e error) error {
			if e != nil {
				return e
			}
			mode := i.Mode()
			if mode.IsRegular() && !mode.IsDir() {
				log.V(4).Info("index %s\n", p)
				abs, err := filepath.Abs(p)
				if err != nil {
					return err
				}
				files = append(files, abs)
			}
			return nil
		})
		if err != nil {
			return files, err
		}
	}
	return files, nil
}

func lint(filename string) error {
	// XXX editorconfig should be able to treat a flux of
	// filenames with caching capabilities.
	def, err := editorconfig.GetDefinitionForFilename(filename)
	if err != nil {
		return fmt.Errorf("Cannot open file %s. %w", filename, err)
	}
	log.V(1).Info("lint", "filename", filename)

	fp, err := os.Open(filename)
	if err != nil {
		return err
	}

	err = readLines(fp, func(index int, data []byte) error {
		var err error

		if def.EndOfLine != "" {
			err = endOfLine(def.EndOfLine, data)
		}

		if err == nil && def.IndentStyle != "" {
			err = indentStyle(def.IndentStyle, data)
		}

		if err == nil && def.Charset != "" {
			err = charset(def.Charset, data)
		}

		if err == nil && def.TrimTrailingWhitespace != nil && *def.TrimTrailingWhitespace {
			err = trimTrailingWhitespace(data)
		}

		if err != nil {
			return fmt.Errorf("line %d: %s", index, err)
		}
		return nil
	})

	return err
}

func main() {
	var flagVersion bool

	klog.InitFlags(nil)

	flag.BoolVar(&flagVersion, "version", false, "print the version number")
	flag.Parse()

	if flagVersion {
		fmt.Printf("eclint %s\n", version)
		return
	}

	log = klogr.New()

	args := flag.Args()
	if len(args) == 0 {
		args = append(args, ".")
	}

	files, err := walk(args...)
	if err != nil {
		log.Error(err, "error while handling the arguments")
		flag.Usage()
		os.Exit(1)
		return
	}
	log.V(1).Info("files", "count", len(files))

	c := 0
	for _, filename := range files {
		err := lint(filename)
		if err != nil {
			log.V(4).Info("lint error", "filename", filename, "error", err)
			fmt.Printf("filename %s: %s", filename, err)
			c++
		}
	}
	if c > 0 {
		os.Exit(1)
	}
}
