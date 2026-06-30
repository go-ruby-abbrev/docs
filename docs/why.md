# Why pure Go

`go-ruby-abbrev/abbrev` reimplements Ruby's `abbrev` standard library in **pure Go, with cgo
disabled**. The slice of Ruby it covers is **deterministic and
interpreter-independent**: computing the set of unambiguous prefix abbreviations
of a word list is a pure function of those words, with no live binding and no
evaluation of arbitrary Ruby. That is exactly the part that can — and should —
live as a standalone Go library, separate from the interpreter.

## Extracted from rbgo, reusable by anyone

This library is a faithful, byte-for-byte port of upstream `abbrev.rb`
(Akinori MUSHA), packaged as a reusable standalone library so that:

- any Go program can import `github.com/go-ruby-abbrev/abbrev` directly, with no Ruby runtime;
- the dependency runs the *other* way — `rbgo` binds this module as a native
  module (the same pattern as [go-ruby-regexp](https://github.com/go-ruby-regexp)
  and [go-ruby-erb](https://github.com/go-ruby-erb)), rather than this module
  depending on the interpreter;
- the behaviour is pinned by a **differential oracle** against the system
  `ruby` (`Abbrev.abbrev`), independent of any one consumer.

## Why pure Go matters here

Because the library is CGO-free and dependency-free, it:

- cross-compiles to every Go target with no C toolchain, and links into a single
  static binary;
- has **no dependency on the Ruby runtime** — the dependency runs the other way;
- can be differentially tested against the `ruby` binary wherever one is on
  `PATH`, while the cross-arch and Windows lanes (where `ruby` is absent) still
  validate the library itself.

See [Usage & API](api.md) for the surface and [Roadmap](roadmap.md) for what is
in scope.
