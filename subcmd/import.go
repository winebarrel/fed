package subcmd

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/kanmu/kasa"
	"github.com/kanmu/kasa/esa/model"
	"github.com/kanmu/kasa/postname"
)

type ImportCmd struct {
	File   string `arg:"" help:"Source file (stdin:'-')."`
	Path   string `arg:"" help:"Post name or Post URL('https://<TEAM>.esa.io/posts/<NUM>' or '//<NUM>')."`
	Notice bool   `negatable:"" default:"true" help:"Post with notify."`
	Wip    bool   `negatable:"" help:"Post as WIP."`
}

func (cmd *ImportCmd) Run(ctx *kasa.Context) error {
	if cmd.File == "-" {
		return cmd.importFile(ctx, os.Stdin, cmd.Path)
	} else {
		fi, err := os.Stat(cmd.File)

		if err != nil {
			return err
		}

		if fi.IsDir() {
			return cmd.importDir(ctx)
		} else {
			f, err := os.OpenFile(cmd.File, os.O_RDONLY, 0)

			if err != nil {
				return err
			}

			defer f.Close()
			return cmd.importFile(ctx, f, cmd.Path)
		}
	}
}

func (cmd *ImportCmd) importFile(ctx *kasa.Context, file io.Reader, path string) error {
	cat, name := postname.Split(path)

	if name == "" {
		name = filepath.Base(cmd.File)
	}

	newPost := &model.NewPostBody{
		Name:     name,
		Category: cat,
		Wip:      &cmd.Wip,
	}

	bodyMd, err := io.ReadAll(file)

	if err != nil {
		return err
	}

	newPost.BodyMd = string(bodyMd)
	url, err := ctx.Driver.Post(newPost, 0, cmd.Notice)

	if err != nil {
		return err
	}

	ctx.Fmt.Println(url)

	return nil
}

func (cmd *ImportCmd) importDir(ctx *kasa.Context) error {
	return filepath.Walk(cmd.File, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		path, err = filepath.Abs(path)

		if err != nil {
			return err
		}

		root, err := filepath.Abs(cmd.File)

		if err != nil {
			return err
		}

		f, err := os.OpenFile(path, os.O_RDONLY, 0)

		if err != nil {
			return err
		}

		path = strings.TrimPrefix(path, root)
		path = filepath.Join(cmd.Path, path)
		path = strings.TrimPrefix(path, "/")

		return cmd.importFile(ctx, f, path)
	})

}
