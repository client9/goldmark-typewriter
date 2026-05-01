// Package typewriter provides a goldmark extension that applies typewriter
// conversions as an AST transformer.
//
// All Category constants and Option constructors are defined here so callers
// need only one import. For preprocessing raw markdown source before parsing,
// use github.com/client9/typewriter directly.
package typewriter

import (
	tw "github.com/client9/typewriter"
	gm "github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
)

// Category is an alias for the core Category type.
type Category = tw.Category

// Category constants — all active by default.
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

// UnicodeStyle is an alias for the core UnicodeStyle type.
type UnicodeStyle = tw.UnicodeStyle

// UnicodeStyle constants.
const (
	Bold        = tw.Bold
	Italic      = tw.Italic
	BoldItalic  = tw.BoldItalic
	Monospace   = tw.Monospace
	Superscript = tw.Superscript
	Subscript   = tw.Subscript
)

// Option configures the extension.
type Option func(*tw.Config)

// WithCategory sets the active categories to exactly c, replacing the default.
func WithCategory(c Category) Option {
	return func(cfg *tw.Config) { cfg.Categories = c }
}

// WithoutCategory removes one or more categories from the active set.
func WithoutCategory(c Category) Option {
	return func(cfg *tw.Config) { cfg.Categories &^= c }
}

// WithMapping adds or overrides a single character conversion. Set to to an
// empty string to leave the character unchanged.
func WithMapping(from, to string) Option {
	return func(cfg *tw.Config) {
		if cfg.Overrides == nil {
			cfg.Overrides = make(map[string]string)
		}
		cfg.Overrides[from] = to
	}
}

// WithBold converts runs of Unicode bold characters, wrapping with prefix and suffix.
// Empty prefix and suffix strips to plain ASCII.
func WithBold(prefix, suffix string) Option {
	return withRun(tw.Bold, prefix, suffix)
}

// WithItalic converts runs of Unicode italic characters, wrapping with prefix and suffix.
func WithItalic(prefix, suffix string) Option {
	return withRun(tw.Italic, prefix, suffix)
}

// WithBoldItalic converts runs of Unicode bold-italic characters, wrapping with prefix and suffix.
func WithBoldItalic(prefix, suffix string) Option {
	return withRun(tw.BoldItalic, prefix, suffix)
}

// WithMonospace converts runs of Unicode monospace characters, wrapping with prefix and suffix.
func WithMonospace(prefix, suffix string) Option {
	return withRun(tw.Monospace, prefix, suffix)
}

// WithSuperscript converts runs of superscript characters, wrapping with prefix and suffix.
// A common convention is prefix "^" with empty suffix.
func WithSuperscript(prefix, suffix string) Option {
	return withRun(tw.Superscript, prefix, suffix)
}

// WithSubscript converts runs of subscript characters, wrapping with prefix and suffix.
func WithSubscript(prefix, suffix string) Option {
	return withRun(tw.Subscript, prefix, suffix)
}

func withRun(style tw.UnicodeStyle, prefix, suffix string) Option {
	return func(cfg *tw.Config) {
		cfg.Runs = append(cfg.Runs, tw.RunStyle{Style: style, Prefix: prefix, Suffix: suffix})
	}
}

// Extension is the goldmark extension. Create with New.
type Extension struct {
	r *tw.Replacer
}

// New creates the goldmark extension. With no options all Default categories
// are active and no Unicode style runs are converted.
func New(opts ...Option) *Extension {
	cfg := tw.Config{Categories: tw.Default}
	for _, o := range opts {
		o(&cfg)
	}
	return &Extension{r: tw.New(cfg)}
}

// Extend implements goldmark.Extender.
func (e *Extension) Extend(m gm.Markdown) {
	m.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(&transformer{r: e.r}, 100),
		),
	)
}
