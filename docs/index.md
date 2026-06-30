# go-ruby-abbrev documentation

**Ruby's abbrev — unambiguous prefix abbreviations — in pure Go, MRI-compatible, no cgo.**

`go-ruby-abbrev/abbrev` is a faithful, pure-Go (zero cgo) reimplementation of Ruby's abbrev unambiguous-prefix library,
matching reference Ruby (MRI) byte-for-byte. The module path is
`github.com/go-ruby-abbrev/abbrev`.

It is a **standalone, reusable library**: the module is importable by any Go
program, and it is the backend bound into
[go-embedded-ruby](https://github.com/go-embedded-ruby/ruby) by `rbgo` as a
native module — just like [go-ruby-regexp](https://github.com/go-ruby-regexp)
and [go-ruby-erb](https://github.com/go-ruby-erb). The dependency runs the other
way: this library has **no dependency on the Ruby runtime**.

!!! success "Status: complete — MRI byte-exact"
    A faithful, byte-for-byte port of upstream `abbrev.rb` (Akinori MUSHA): a single **`Abbrev(words, prefix...)`** covering both **`Abbrev.abbrev`** and the **`Array#abbrev`** core extension. It computes every **unambiguous prefix** plus each full word, with the optional **literal-string** prefix filter, and reproduces MRI's edge cases — ambiguous prefixes dropped, empty and duplicate words, and multibyte words split on **characters**. Validated by a **differential oracle** against the system `ruby` / `Abbrev.abbrev` at 100% coverage, `gofmt` + `go vet` clean, CI green across the six 64-bit Go targets and three OSes.

## Quick taste

```go
abbrev.Abbrev([]string{"ruby", "rules"})
// => {"rub":"ruby", "ruby":"ruby", "rule":"rules", "rules":"rules"}
// "r" and "ru" are omitted: each is a prefix of both words (ambiguous).

abbrev.Abbrev([]string{"car", "cone"}, "ca")
// => {"ca":"car", "car":"car"}   (the optional prefix keeps only words starting with it)
```

## Repositories

| Repo | What it is |
| --- | --- |
| [`abbrev`](https://github.com/go-ruby-abbrev/abbrev) | the library — Ruby's abbrev unambiguous-prefix library in pure Go |
| [`docs`](https://github.com/go-ruby-abbrev/docs) | this documentation site (MkDocs Material, versioned with mike) |
| [`go-ruby-abbrev.github.io`](https://github.com/go-ruby-abbrev/go-ruby-abbrev.github.io) | the organization landing page (Hugo) |
| [`brand`](https://github.com/go-ruby-abbrev/brand) | logo and brand assets |


## Principles

- **Pure Go, `CGO_ENABLED=0`** — trivial cross-compilation, a single static
  binary, no C toolchain.
- **MRI byte-exact.** Output matches reference Ruby exactly, not approximately,
  validated by a differential oracle against the `ruby` binary.
- **Standalone & reusable.** No dependency on the Ruby runtime — the dependency
  runs the other way.
- **100% test coverage** is the target, enforced as a CI gate.

## Where to go next

- [Why pure Go](why.md) — why this slice of Ruby is deterministic enough to live
  as a standalone, interpreter-independent Go library.
- [Usage & API](api.md) — the public surface and worked examples.
- [Roadmap](roadmap.md) — what is done and what is downstream by design.

Source lives at [github.com/go-ruby-abbrev/abbrev](https://github.com/go-ruby-abbrev/abbrev).
