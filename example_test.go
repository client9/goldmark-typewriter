package typewriter_test

import (
	"os"

	typewriter "github.com/client9/goldmark-typewriter"
	tw "github.com/client9/typewriter"
	gm "github.com/yuin/goldmark"
)

// ExampleNew demonstrates the typical usage: register the extension with a
// goldmark instance and convert markdown containing typographic Unicode
// characters to their ASCII equivalents.
func ExampleNew() {
	md := gm.New(gm.WithExtensions(typewriter.New()))
	_ = md.Convert([]byte("mix ½ cup, cost © 2024, wait…"), os.Stdout)
	// Output:
	// <p>mix 1/2 cup, cost (c) 2024, wait...</p>
}

// ExampleNew_withoutCategory shows how to disable a category. Here the Math
// category is removed so the multiplication sign passes through unchanged.
func ExampleNew_withoutCategory() {
	ext := typewriter.New(typewriter.WithoutCategory(typewriter.Math))
	md := gm.New(gm.WithExtensions(ext))
	_ = md.Convert([]byte("10×"), os.Stdout)
	// Output:
	// <p>10×</p>
}

// ExampleWithMapping shows how to override a single conversion. Mapping "…"
// to itself prevents the ellipsis from being expanded to three dots.
func ExampleWithMapping() {
	ext := typewriter.New(typewriter.WithMapping("…", "…"))
	md := gm.New(gm.WithExtensions(ext))
	_ = md.Convert([]byte("wait…"), os.Stdout)
	// Output:
	// <p>wait…</p>
}

// ExampleNew_replaceBytes shows the two-pass approach for normalising content
// inside code spans, which the AST transformer cannot reach. Call
// tw.ReplaceBytes on the raw source first, then parse with goldmark.
func ExampleNew_replaceBytes() {
	src := []byte("outside … but `inside … code`")
	md := gm.New(gm.WithExtensions(typewriter.New()))

	// Extension form: prose converted, code span preserved.
	_ = md.Convert(src, os.Stdout)

	// ReplaceBytes: everything converted, including inside code spans.
	_ = md.Convert(tw.ReplaceBytes(src), os.Stdout)
	// Output:
	// <p>outside ... but <code>inside … code</code></p>
	// <p>outside ... but <code>inside ... code</code></p>
}
