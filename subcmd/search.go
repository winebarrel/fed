package subcmd

import (
	"fmt"
	"sort"

	"github.com/winebarrel/kasa"
)

type SearchCmd struct {
	Query string `arg:"" help:"Search query."`
	Page  int    `short:"p" default:"1" help:"Page number."`
}

func (cmd *SearchCmd) Run(ctx *kasa.Context) error {
	posts, hasMore, err := ctx.Driver.Search(cmd.Query, cmd.Page)

	if err != nil {
		return err
	}

	sort.Slice(posts, func(i, j int) bool { return posts[i].FullName < posts[j].FullName })

	for _, v := range posts {
		fmt.Println(v.ListString())
	}

	if hasMore {
		fmt.Println("(has more page. Try increasing `-p NUM`)")
	}

	return nil
}
