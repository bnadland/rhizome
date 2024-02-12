package markdown

import (
	"bytes"
	"fmt"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/hashtag"
	"go.abhg.dev/goldmark/wikilink"
)

type resolver struct{}

func (*resolver) ResolveHashtag(node *hashtag.Node) (destination []byte, err error) {
	return []byte(fmt.Sprintf("/c/%s", string(node.Tag))), nil
}

func (*resolver) ResolveWikilink(node *wikilink.Node) (destination []byte, err error) {
	return []byte(fmt.Sprintf("/p/%s", string(node.Target))), nil
}

func renderer() goldmark.Markdown {
	return goldmark.New(
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithExtensions(
			&hashtag.Extender{
				Resolver: &resolver{},
			},
			&wikilink.Extender{
				Resolver: &resolver{},
			},
		),
	)
}

func Markdown(content string) (string, error) {
	var buf bytes.Buffer
	md := renderer()
	if err := md.Convert([]byte(content), &buf); err != nil {
		return "", err
	}
	return bluemonday.UGCPolicy().Sanitize(buf.String()), nil
}

func Hashtags(content string) ([]string, error) {
	md := renderer()
	doc := md.Parser().Parse(text.NewReader([]byte(content)))
	var hashtags []string
	ast.Walk(doc, func(node ast.Node, enter bool) (ast.WalkStatus, error) {
		if n, ok := node.(*hashtag.Node); ok && enter {
			hashtags = append(hashtags, string(n.Tag))
		}
		return ast.WalkContinue, nil
	})
	return hashtags, nil
}
