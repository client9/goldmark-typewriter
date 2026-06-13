package typewriter

import (
	tw "github.com/client9/typewriter"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// transformer walks ast.KindText nodes and applies typewriter replacements.
// Code and raw HTML nodes are skipped; see the architecture note in CLAUDE.md.
type transformer struct {
	r *tw.Replacer
}

func (t *transformer) Transform(doc *ast.Document, reader text.Reader, pc parser.Context) {
	source := reader.Source()
	type edit struct {
		node   *ast.Text
		result string
	}
	var edits []edit
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch n.Kind() {
		case ast.KindCodeBlock, ast.KindFencedCodeBlock, ast.KindCodeSpan,
			ast.KindHTMLBlock, ast.KindRawHTML:
			// The HTML renderer reads code content directly from source[]
			// via segment.Value(source) and expects ast.Text children, not
			// ast.String. Replacing those nodes would panic the renderer.
			// Use tw.ReplaceBytes on raw source to normalise code content.
			return ast.WalkSkipChildren, nil
		case ast.KindText:
			node := n.(*ast.Text)
			src := string(node.Segment.Value(source))
			if result := t.r.Replace(src); result != src {
				edits = append(edits, edit{node, result})
			}
		}
		return ast.WalkContinue, nil
	})
	// Apply replacements after the walk. Calling ReplaceChild inside the walk
	// callback causes goldmark's RemoveChild to zero the replaced node's
	// nextSibling pointer, which terminates the sibling iteration early and
	// silently skips all remaining text nodes in the same parent.
	for _, e := range edits {
		replaceText(e.node, e.result)
	}
}

func replaceText(node *ast.Text, result string) {
	// ast.String has no SoftLineBreak/HardLineBreak fields and renderString
	// emits nothing after the text. Embed the newline so the line break
	// survives in rendered output.
	if node.SoftLineBreak() || node.HardLineBreak() {
		result += "\n"
	}
	if p := node.Parent(); p != nil {
		p.ReplaceChild(p, node, ast.NewString([]byte(result)))
	}
}
