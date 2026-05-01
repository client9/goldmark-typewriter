// Package typewriter provides a goldmark extension that applies typewriter
// conversions as an AST transformer. For preprocessing raw markdown source
// before parsing, use the parent typewriter package directly.
//
// Type aliases and re-exported identifiers let callers use this package as
// their sole import — no need to also import github.com/client9/typewriter.
package typewriter

import (
	tw "github.com/client9/typewriter"
	gm "github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
)

// Type aliases — identical to the core types, no conversion required.
type (
	Category = tw.Category
	Option   = tw.Option
)

// Category constants re-exported by value.
const (
	Quotes      = tw.Quotes
	Dashes      = tw.Dashes
	Ellipsis    = tw.Ellipsis
	Fractions   = tw.Fractions
	Symbols     = tw.Symbols
	Math        = tw.Math
	Ligatures   = tw.Ligatures
	Bullets     = tw.Bullets
	Spaces      = tw.Spaces
	Default     = tw.Default
	CategoryAll = tw.CategoryAll
)

// Option constructors re-exported as package-level vars.
var (
	WithCategory    = tw.WithCategory
	WithoutCategory = tw.WithoutCategory
	WithMapping     = tw.WithMapping
)

// Extension is the goldmark extension. Create with New.
type Extension struct {
	r *tw.Replacer
}

// New creates the goldmark extension. With no options all Default categories
// are active. Accepts the same options as tw.New.
func New(opts ...tw.Option) *Extension {
	return &Extension{r: tw.New(opts...)}
}

// Extend implements goldmark.Extender.
func (e *Extension) Extend(m gm.Markdown) {
	m.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(&transformer{r: e.r}, 100),
		),
	)
}
