<!-- SPDX-License-Identifier: BSD-3-Clause -->
# `go-ruby-abbrev` library-level benchmark harness

Reproducible, cross-runtime benchmark of the **pure-Go `go-ruby-abbrev/abbrev`
library** against the reference Ruby runtimes (MRI, MRI + YJIT, JRuby,
TruffleRuby). It measures the **library primitive** through its Go API, isolated
from the rbgo interpreter, so the numbers answer: *is the pure-Go implementation
as fast as the reference runtime's own `abbrev`?*

## Layout

- `go/`            — self-contained Go driver; `go.mod` pins the published library.
- `ruby/abbrev.rb` — the equivalent workload; `ruby/_harness.rb` is the shared timer.
- `run.sh`         — runs every available runtime and prints one Markdown table per
  sub-benchmark (ns/op + ratio vs MRI).

## Run

```sh
bash benchmarks/run.sh
```

Environment knobs: `OUTER` (timed passes, default 25), `WARM` (untimed warm-up
passes, default 3), and `RUBY`/`JRUBY`/`TRUFFLERUBY` to select runtime binaries.

## Method

Each process runs `WARM` untimed passes (to let the JVM/GraalVM JITs warm up),
then `OUTER` timed passes of a fixed inner loop, timed with a monotonic clock;
the **best** pass is reported as **ns/op**. Interpreter start-up is outside the
timed region. The Go driver and the Ruby script build the **identical** fixed
word list and compute `Abbrev.abbrev` over it; before any timing, `run.sh` checks
the Go abbrev→word table **byte-identical to MRI** (each driver's `--dump` mode
emits the sorted table) and aborts on any mismatch. Results are published, dated,
in `../docs/performance.md`.
