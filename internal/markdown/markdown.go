package markdown

import (
	"bytes"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
)

func Markdown(content string) (string, error) {
	var buf bytes.Buffer
	md := goldmark.New()
	if err := md.Convert([]byte(content), &buf); err != nil {
		return "", err
	}
	return bluemonday.UGCPolicy().Sanitize(buf.String()), nil
}
