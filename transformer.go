package typewriter

import (
	tw "github.com/client9/typewriter"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type transformer struct {
	r *tw.Replacer
}

func (t *transformer) Transform(doc *ast.Document, reader text.Reader, pc parser.Context) {
	source := reader.Source()
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch n.Kind() {
		case ast.KindCodeBlock, ast.KindFencedCodeBlock, ast.KindCodeSpan,
			ast.KindHTMLBlock, ast.KindRawHTML:
			// The HTML renderer reads code content directly from source[]
			// via segment.Value(source) and expects ast.Text children, not
			// ast.String. Replacing nodes here would panic the renderer.
			// Use tw.ReplaceBytes on raw source to normalise code content.
			return ast.WalkSkipChildren, nil
		case ast.KindText:
			replaceText(n.(*ast.Text), source, t.r)
		}
		return ast.WalkContinue, nil
	})
}

func replaceText(node *ast.Text, source []byte, r *tw.Replacer) {
	src := string(node.Segment.Value(source))
	result := r.Replace(src)
	if result == src {
		return
	}
	node.Parent().ReplaceChild(node.Parent(), node, ast.NewString([]byte(result)))
}
