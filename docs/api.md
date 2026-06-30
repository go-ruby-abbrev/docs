# Usage & API

The public API lives at the module root (`github.com/go-ruby-abbrev/abbrev`). It is **Ruby-shaped but Go-idiomatic**: a single `Abbrev` function collapses both Ruby entry points (`Abbrev.abbrev` and `Array#abbrev`) into one call with a variadic optional prefix, returning a plain `map[string]string`.

!!! success "Status: implemented"
    The library is built and importable as `github.com/go-ruby-abbrev/abbrev`, bound into
    `rbgo` as a native module; see [Roadmap](roadmap.md).

## Install

```sh
go get github.com/go-ruby-abbrev/abbrev
```

## Worked example

```go
abbrev.Abbrev([]string{"ruby", "rules"})
// => map[string]string{
//      "rub":   "ruby",  "ruby":  "ruby",
//      "rule":  "rules", "rules": "rules",
//    }
// "r" and "ru" are omitted: each is a prefix of both words (ambiguous).

abbrev.Abbrev([]string{"car", "cone"}, "ca")
// => map[string]string{"ca": "car", "car": "car"}
// The optional prefix keeps only words starting with it.
```

## Shape

```go
// Abbrev returns the unambiguous-abbreviation map for words: every prefix that
// identifies exactly one word, plus each full word. An optional literal prefix
// keeps only words starting with it.
func Abbrev(words []string, prefix ...string) map[string]string
```

This is the idiomatic Go shape of both Ruby entry points:

| Ruby                                | Go                                  |
| ----------------------------------- | ----------------------------------- |
| `Abbrev.abbrev(words)`              | `abbrev.Abbrev(words)`              |
| `Abbrev.abbrev(words, prefix)`      | `abbrev.Abbrev(words, prefix)`      |
| `words.abbrev` (`Array#abbrev`)     | `abbrev.Abbrev(words)`              |
| `words.abbrev(prefix)`              | `abbrev.Abbrev(words, prefix)`      |

## MRI fidelity

The port reproduces MRI exactly, including its edge cases:

- **Ambiguous prefixes are dropped.** A prefix shared by two or more words never
  appears; each full word always maps to itself.
- **Empty words** contribute no prefixes but still map to themselves:
  `Abbrev([]string{""})` → `{"": ""}`.
- **Duplicate words** collapse: `Abbrev([]string{"dog", "dog"})` → `{"dog": "dog"}`.
- **The prefix is a literal string, not a pattern** (MRI anchors it as
  `/\A<quoted>/`): `Abbrev([]string{"a.b", "axb"}, "a.")` → `{"a.b": "a.b", "a.": "a.b"}`.
- **Multibyte words split on characters**, not bytes, matching Ruby's
  `String#[]` semantics: `Abbrev([]string{"café", "cane"})` includes `"caf"`.

The optional Ruby `pattern` may also be a `Regexp`; this library models the
common String-prefix case (which is what `rbgo` binds), so a host that needs the
regexp form filters the word list before calling `Abbrev`.

## MRI conformance

Correctness is defined by reference Ruby. A **differential oracle** runs the same
inputs through the live `ruby` binary (`Abbrev.abbrev`) and asserts byte-identical
maps. The oracle skips itself where `ruby` is absent (the Windows lane, the qemu
cross-arch lanes) and is gated on `RUBY_VERSION >= "4.0"`, so the deterministic
suite alone keeps the gate green everywhere.

## Relationship to Ruby

`go-ruby-abbrev/abbrev` is **standalone and reusable**, and is the backend bound into
[go-embedded-ruby](https://github.com/go-embedded-ruby/ruby) by `rbgo` as a
native module — the same way [go-ruby-regexp](https://github.com/go-ruby-regexp)
and [go-ruby-erb](https://github.com/go-ruby-erb) are bound. The dependency runs
the other way: this library has no dependency on the Ruby runtime.
