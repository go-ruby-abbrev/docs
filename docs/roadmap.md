# Roadmap

`go-ruby-abbrev/abbrev` is grown **test-first**, each capability differential-tested against MRI
rather than built in isolation. Ruby's `abbrev` — the
deterministic, interpreter-independent prefix-abbreviation algorithm — is
**complete**.

| Stage | What | Status |
| --- | --- | --- |
| Unambiguous abbreviations | Every prefix that identifies exactly one word plus each full word; a prefix shared by two or more words never appears. | **Done** |
| Optional prefix filter | `Abbrev(words, prefix)` keeps only words starting with `prefix`, treated as a literal string (MRI anchors it as `/\A<quoted>/`), not a pattern. | **Done** |
| Both Ruby entry points | The Go shape of `Abbrev.abbrev(words[, prefix])` and the `Array#abbrev` core extension, collapsed into one variadic `Abbrev`. | **Done** |
| MRI edge cases | Empty words map to themselves with no prefixes; duplicate words collapse; multibyte words split on characters, not bytes, matching `String#[]`. | **Done** |
| Differential oracle & coverage | Runtime-free cases drive coverage to 100%; a differential oracle runs the same inputs through `ruby` (`Abbrev.abbrev`) and asserts byte-identical maps. gofmt + go vet clean, green across six arches and three OSes. | **Done** |

## Documented out-of-scope boundaries

These are **deliberate**, recorded so the module's surface is unambiguous:

- **No `Regexp` prefix.** MRI's optional `pattern` may be a `Regexp`; this library models the common literal String-prefix case (what `rbgo` binds). A host that needs the regexp form filters the word list before calling `Abbrev`.
- **No interpreter.** The library implements the deterministic algorithm; it never runs arbitrary Ruby. Wiring it into a live object model is the consumer's job — that is why `rbgo` binds this module rather than the reverse.
- **Reference is reference Ruby (MRI).** Byte-for-byte conformance targets MRI's `Abbrev.abbrev`, pinned by the differential oracle.
- **Standalone & reusable.** The module has no dependency on the Ruby runtime; the dependency runs the other way.

See [Usage & API](api.md) for the surface and [Why pure Go](why.md) for the
deterministic/interpreter split.
