# goldmark-typewriter

A [goldmark](https://github.com/yuin/goldmark) extension that converts typographic
("smart") Unicode characters back to their plain ASCII equivalents.

Built on [github.com/client9/typewriter](https://github.com/client9/typewriter). All
`Category` constants and option functions are re-exported so you only need one import.

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
    typewriter.New(typewriter.WithoutCategory(typewriter.Math)),
))
```

## What it converts

All categories are active by default.

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
| Spaces | NBSP, thin, en, em, figure, hair spaces | plain space |

## Prose vs code content

| | Prose | Code spans / fenced blocks |
|---|---|---|
| goldmark extension | converted | **preserved** |
| `typewriter.ReplaceBytes` (preprocessor) | converted | converted |

The extension preserves code content because goldmark's HTML renderer reads code spans
and fenced blocks directly from the original source bytes — AST-level replacement is not
possible there. This is a constraint of goldmark's architecture, not a policy choice.

To also normalise code content (for example, smart quotes inside a shell command pasted
from a blog), use `typewriter.ReplaceBytes` on the raw source before parsing:

```go
import tw "github.com/client9/typewriter"

clean := tw.ReplaceBytes(src)
md.Convert(clean, &buf)
```

## Configuration

### Enable only specific categories

`WithCategory` sets the active categories to exactly what you pass, replacing the default:

```go
// Only convert dashes and ellipses.
typewriter.New(typewriter.WithCategory(typewriter.Dashes | typewriter.Ellipsis))
```

### Disable specific categories

`WithoutCategory` removes categories from the active set:

```go
typewriter.New(typewriter.WithoutCategory(typewriter.Math))
```

Options compose left-to-right:

```go
typewriter.New(
    typewriter.WithCategory(typewriter.CategoryAll),
    typewriter.WithoutCategory(typewriter.Math | typewriter.Bullets),
)
```

### Override or exclude individual mappings

```go
typewriter.New(
    typewriter.WithMapping("—", "--"),  // prefer -- over --- for em dash
    typewriter.WithMapping("×", ""),    // leave × unchanged (empty = pass through)
    typewriter.WithMapping("°", "deg"), // add a mapping not in builtins
)
```

## Using with the typographer

goldmark's [typographer extension](https://github.com/yuin/goldmark#built-in-extensions)
converts ASCII punctuation to typographic Unicode. typewriter is its complement.

They **cannot be meaningfully combined in a single goldmark instance**: the typographer is
an inline parser and always fires during tokenisation, before any AST transformer runs.
In a single instance the typographer fires first, converting ASCII → typographic; the
typewriter transformer then has nothing left to do (and cannot undo `ast.String` nodes
the typographer produced).

For consistent smart-typography output from mixed-source input — where some content
already contains curly quotes and some does not — use a two-pass approach:

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
  conversion package; use directly when you don't need goldmark integration
