# goldmark-typewriter
[![Go Reference](https://pkg.go.dev/badge/github.com/client9/goldmark-typewriter.svg)](https://pkg.go.dev/github.com/client9/goldmark-typewriter)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://github.com/client9/goldmark-typewriter/actions/workflows/go.yml/badge.svg)](https://github.com/client9/goldmark-typewriter/actions)

A [goldmark](https://github.com/yuin/goldmark) extension that converts typographic
("smart") Unicode characters — and Unicode style variants like bold and italic — back
to their plain ASCII equivalents.

Built on [github.com/client9/typewriter](https://github.com/client9/typewriter). 

See also [github.com/client9/demoji](https://github.com/client9/demoji) and [github.com/client9/goldmark-demoji](https://github.com/client9/goldmark-demoji) for emoji conversion and normalization.

## Install

```sh
go get github.com/client9/goldmark-typewriter
```

## Quick start

```go
import typewriter "github.com/client9/goldmark-typewriter"

md := goldmark.New(goldmark.WithExtensions(typewriter.New()))
```

With options:

```go
md := goldmark.New(goldmark.WithExtensions(
    typewriter.New(
        typewriter.WithoutCategory(typewriter.Math),
        typewriter.WithBold("**", "**"),
        typewriter.WithItalic("_", "_"),
    ),
))
```

## What it converts

### Character substitutions (all active by default)

| Category | Examples | Result |
|----------|---------|--------|
| Quotes | `"` `"` `'` `'` `«` `»` `„` | `"` `'` `<<` `>>` |
| Dashes | em dash `—`, en dash `–`, minus `−` | `---` `--` `-` |
| Ellipsis | `…` | `...` |
| Fractions | `½` `¼` `¾` `⅓` `⅛` | `1/2` `1/4` `3/4` `1/3` `1/8` |
| Symbols | `©` `®` `™` | `(c)` `(r)` `(tm)` |
| Math | `×` `÷` `≠` `≤` `≥` `→` | `x` `/` `!=` `<=` `>=` `->` |
| Ligatures | `ﬁ` `ﬂ` `ﬀ` `ﬃ` | `fi` `fl` `ff` `ffi` |
| Bullets | `•` `†` `‡` | `*` `*` `**` |
| Spaces | NBSP, thin, en, em, figure, hair, U+2028, U+2029 | plain space |

### Unicode style variants (opt-in)

Runs of styled Unicode characters — common in copy-pasted LinkedIn posts, Twitter
formatting tricks, and AI-generated content — are detected and wrapped.

| Option | Input | Output |
|--------|-------|--------|
| `WithBold("**", "**")` | `𝗛𝗲𝗹𝗹𝗼` | `**Hello**` |
| `WithItalic("_", "_")` | `𝘸𝘰𝘳𝘭𝘥` | `_world_` |
| `WithBoldItalic("***", "***")` | `𝙃𝙚𝙡𝙡𝙤` | `***Hello***` |
| `WithMonospace("` + "`" + `", "` + "`" + `")` | `𝙷𝚎𝚕𝚕𝚘` | `` `Hello` `` |
| `WithSuperscript("^", "")` | `mc²` | `mc^2` |
| `WithSubscript("", "")` | `H₂O` | `H2O` |

## Prose vs code content

| | Prose | Code spans / fenced blocks |
|---|---|---|
| goldmark extension | converted | **preserved** |
| `typewriter.ReplaceBytes` (preprocessor) | converted | converted |

The extension preserves code content because goldmark's HTML renderer reads code spans
and fenced blocks directly from the original source bytes — AST-level replacement is not
possible there. This is a constraint of goldmark's architecture, not a policy choice.

To also normalise code content (e.g., smart quotes inside a pasted shell command), use
`typewriter.ReplaceBytes` on raw source before parsing:

```go
import tw "github.com/client9/typewriter"

clean := tw.ReplaceBytes(src)
md.Convert(clean, &buf)
```

## Configuration

### Enable only specific categories

```go
// Only convert dashes and ellipses.
typewriter.New(typewriter.WithCategory(typewriter.Dashes | typewriter.Ellipsis))
```

### Disable specific categories

```go
typewriter.New(typewriter.WithoutCategory(typewriter.Math))
```

### Override or exclude individual characters

```go
typewriter.New(
    typewriter.WithMapping("—", "--"),  // prefer -- over --- for em dash
    typewriter.WithMapping("×", ""),    // leave × unchanged (empty = pass through)
    typewriter.WithMapping("°", "deg"), // add a mapping not in builtins
)
```

### Convert Unicode bold/italic to markdown

```go
typewriter.New(
    typewriter.WithBold("**", "**"),
    typewriter.WithItalic("_", "_"),
)
```

### Convert Unicode bold/italic to HTML

```go
typewriter.New(
    typewriter.WithBold("<b>", "</b>"),
    typewriter.WithItalic("<i>", "</i>"),
)
```

### Superscripts and subscripts

```go
typewriter.New(
    typewriter.WithSuperscript("^", ""),  // E=mc² → E=mc^2
    typewriter.WithSubscript("", ""),     // H₂O  → H2O
)
```

## Using with the typographer

goldmark's typographer and this extension **cannot be meaningfully combined in a single
goldmark instance**. The typographer is an inline parser (priority 9999) that fires
during tokenisation — before any AST transformer runs. It converts ASCII → typographic
as `ast.String` nodes; the typewriter transformer walks only `ast.KindText` and cannot
undo that.

For consistent smart-typography output from mixed-source input, use a two-pass approach:

```go
import (
    tw "github.com/client9/typewriter"
    "github.com/yuin/goldmark/extension"
)

// Step 1: strip all typographic characters from the raw markdown source.
clean := tw.ReplaceBytes(src)

// Step 2: render with the typographer — consistent output regardless of input source.
md := goldmark.New(goldmark.WithExtensions(extension.Typographer))
md.Convert(clean, &buf)
```

## Related

- [github.com/client9/typewriter](https://github.com/client9/typewriter) — the core
  package; use directly for non-goldmark pipelines or raw-byte preprocessing

## License

[MIT](/LICENSE)

