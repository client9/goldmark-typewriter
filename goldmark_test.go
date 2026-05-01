package typewriter_test

import (
	"bytes"
	"testing"

	tw "github.com/client9/typewriter"
	typewriter "github.com/client9/goldmark-typewriter"
	gm "github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

func render(t *testing.T, ext *typewriter.Extension, src string) string {
	t.Helper()
	md := gm.New(gm.WithExtensions(ext))
	var buf bytes.Buffer
	if err := md.Convert([]byte(src), &buf); err != nil {
		t.Fatalf("Convert: %v", err)
	}
	return buf.String()
}

func renderWith(t *testing.T, src string, exts ...gm.Extender) string {
	t.Helper()
	md := gm.New(gm.WithExtensions(exts...))
	var buf bytes.Buffer
	if err := md.Convert([]byte(src), &buf); err != nil {
		t.Fatalf("Convert: %v", err)
	}
	return buf.String()
}

func TestExtensionBasic(t *testing.T) {
	// Verify the goldmark extension applies conversions to prose.
	ext := typewriter.New()
	got := render(t, ext, "wait…")
	if got != "<p>wait...</p>\n" {
		t.Errorf("got %q", got)
	}
}

func TestExtensionOptions(t *testing.T) {
	// Options pass through to the underlying Replacer.
	ext := typewriter.New(typewriter.WithoutCategory(typewriter.Math))
	got := render(t, ext, "10×")
	if got != "<p>10×</p>\n" {
		t.Errorf("× should pass through: got %q", got)
	}
}

func TestCodeContent(t *testing.T) {
	// The goldmark extension form cannot modify code span or fenced code block
	// content: the HTML renderer reads that content directly from the original
	// source bytes (segment.Value(source)) and expects ast.Text children, not
	// ast.String. Replacing those nodes would panic the renderer.
	//
	// ReplaceBytes operates before parsing and normalises everything, including
	// code content.
	src := "outside … but `inside … code`"

	// Extension form: prose converted, code preserved (architecture constraint).
	got := render(t, typewriter.New(), src)
	want := "<p>outside ... but <code>inside … code</code></p>\n"
	if got != want {
		t.Errorf("extension: got %q\nwant %q", got, want)
	}

	// ReplaceBytes: everything converted, including inside code spans.
	stripped := string(tw.ReplaceBytes([]byte(src)))
	got = render(t, typewriter.New(), stripped)
	want = "<p>outside ... but <code>inside ... code</code></p>\n"
	if got != want {
		t.Errorf("ReplaceBytes: got %q\nwant %q", got, want)
	}
}

func TestChaining(t *testing.T) {
	// Genuinely mixed input: first pair is Unicode curly quotes (U+201C/U+201D),
	// second pair is plain ASCII quotes.
	mixed := "\u201Calready curly\u201D and \"straight\""

	t.Logf("input: %s", mixed)

	// Case 1: typographer only — inconsistent output.
	t.Logf("typographer only: %s", renderWith(t, mixed, extension.Typographer))

	// Case 2: single goldmark instance with both extensions — still inconsistent.
	// Typographer (inline parser) fires first; typewriter (AST transformer) cannot
	// undo the ast.String nodes the typographer produced.
	t.Logf("single pass (typewriter+typographer): %s",
		renderWith(t, mixed, typewriter.New(), extension.Typographer))

	// Case 3: ReplaceBytes first, then typographer — consistent output.
	stripped := tw.ReplaceBytes([]byte(mixed))
	t.Logf("after ReplaceBytes: %s", stripped)
	t.Logf("two pass (strip→typographer): %s",
		renderWith(t, string(stripped), extension.Typographer))
}
