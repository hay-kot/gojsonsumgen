// Package commands contains the CLI commands for the application
package commands

import (
	"context"
	"fmt"
	"go/format"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/hay-kot/gojsonsumgen/internal/gojsonsum"
	"github.com/rs/zerolog/log"
)

type Flags struct {
	LogLevel string
}

type Controller struct {
	Flags *Flags
}

type GenJob struct {
	FilePath string
	Defs     []gojsonsum.SumTypeDef
}

func (c *Controller) Gen(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no arguments provided")
	}

	rootdir := args[0]
	rootfs := os.DirFS(rootdir)

	files, err := globfiles(rootfs)
	if err != nil {
		return err
	}

	var jobs []GenJob

	for _, file := range files {
		osfile, err := os.OpenFile(file, os.O_RDONLY, 0)
		if err != nil {
			return err
		}

		defs, err := gojsonsum.ParseDefinitionComments(osfile)
		if err != nil {
			return err
		}

		if len(defs) == 0 {
			continue
		}

    log.Info().Str("path", file).Msg("Found definitions")
		job := GenJob{FilePath: file}

		for _, def := range defs {
			v := gojsonsum.ParseSumTypeDef(def)
			job.Defs = append(job.Defs, v)
		}

		jobs = append(jobs, job)
	}

	for _, job := range jobs {
		packagename := filepath.Base(filepath.Dir(job.FilePath))

		code, err := gojsonsum.Render(packagename, job.Defs)
		if err != nil {
			return err
		}

		bits, err := format.Source([]byte(code))
		if err != nil {
			return err
		}

		// write to adjacent file
		outname := strings.TrimSuffix(job.FilePath, ".go") + ".sumtype.go"

		log.Info().Str("path", outname).Msg("Writing file")
		err = os.WriteFile(outname, bits, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func globfiles(rootfs fs.FS) ([]string, error) {
	const GoExt = ".go"

	var files []string

	err := fs.WalkDir(rootfs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		ext := filepath.Ext(d.Name())

		if ext != GoExt {
			return nil
		}

		files = append(files, path)
		return nil
	})

	return files, err
}
