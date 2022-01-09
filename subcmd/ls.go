package subcmd

import (
	"fmt"
	"sort"

	"github.com/winebarrel/kasa"
)

type LsCmd struct {
	Path      string `arg:"" optional:"" help:"Post name or Post category or Post tag."`
	Page      int    `short:"p" default:"1" help:"Page number."`
	Recursive bool   `short:"r" default:"true" negatable:"" help:"Recursively list posts."`
}

func (cmd *LsCmd) Run(ctx *kasa.Context) error {
	posts, hasMore, err := ctx.Driver.ListOrTagSearch(cmd.Path, cmd.Page, cmd.Recursive)

	if err != nil {
		return err
	}

	sort.Slice(posts, func(i, j int) bool { return posts[i].FullName < posts[j].FullName })

	for _, v := range posts {
		fmt.Println(v.ListString())
	}

	if hasMore {
		fmt.Printf("(has more pages. current page is %d, try `-p %d`)\n", cmd.Page, cmd.Page+1)
	}

	return nil
}
