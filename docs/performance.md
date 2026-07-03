# Performance

`go-ruby-abbrev/abbrev` is the pure-Go library that
[`rbgo`](https://github.com/go-embedded-ruby/ruby) binds for Ruby's `abbrev`. This
page records the **methodology** for a comparative benchmark of that module
against the reference Ruby runtimes, part of the ecosystem-wide per-module
parity suite.

## What is measured

The **same** Ruby script — an `Abbrev.abbrev` computation over a representative word list — is run under every runtime. `rbgo`'s
number reflects **this pure-Go library doing the work**; every other column is
that interpreter's own stdlib. So the comparison is the **Ruby-visible
operation**, apples-to-apples across interpreters. The script prints a
deterministic checksum and its output is checked **byte-identical to MRI** before
timing.

- **Method:** best-of-5 wall time (best, not mean, to suppress scheduler noise);
  single-shot processes, no warm-up beyond the script's own loop.
- **Runtimes:** `ruby` (MRI, the oracle) and `ruby --yjit`; `jruby` (OpenJDK);
  `truffleruby` (GraalVM CE Native).
- The benchmark script and harness live in rbgo's repo under
  [`bench/modules/`](https://github.com/go-embedded-ruby/ruby/tree/main/bench/modules)
  (`abbrev.rb` + `run.sh`). Reproduce:
  `RBGO=./rbgo TRUFFLE=truffleruby bash bench/modules/run.sh 5`.

## Result

## Result (best of 5, ms)

| Runtime | time | vs MRI |
| --- | ---: | ---: |
| **rbgo** (go-ruby-abbrev) | 290 | 1.32× |
| MRI (ruby 4.0.5) | 220 | 1.00× |
| MRI + YJIT | 180 | 0.82× |
| JRuby 10.1.0.0 | 1330 | 6.05× |
| TruffleRuby 34.0.1 | 380 | 1.73× |

rbgo runs on **go-ruby-abbrev** at near parity with MRI (1.32x) on this unambiguous-abbreviation-table workload.

!!! note "Honest framing"
    JRuby and TruffleRuby are timed **cold, single-shot**, so they carry JVM /
    Graal startup on every run — read them as one-shot `ruby file.rb` costs, the
    same way `rbgo` and MRI are measured, not as steady-state JIT numbers. Rows
    that complete in well under ~200 ms carry the most relative noise; treat
    their ratios as order-of-magnitude. These are **real measured numbers** from
    the 2026-06-30 run (Apple M-series; `ruby 4.0.5 +PRISM`, `jruby 10.1.0.0`,
    `truffleruby 34.0.1`) — nothing is fabricated or cherry-picked.

## Library-level benchmark (Go API vs runtimes) — 2026-07-03

This section measures the **pure-Go library directly, through its Go API** — not
the `rbgo` interpreter path recorded above. It isolates the library primitive
from Ruby-interpreter dispatch, answering the parity question head-on: *is the
pure-Go implementation as fast as the reference runtime's own `abbrev`?* The
**same workload, same inputs, same iteration counts** run through the Go library
and through each reference runtime's stdlib; outputs were checked
**byte-identical to MRI** before any timing.

- **Host:** Apple M4 Max (`Mac16,5`, arm64), macOS — **date 2026-07-03**.
- **Runtimes:** Go 1.26.4 · MRI `ruby 4.0.5 +PRISM` · MRI + YJIT · JRuby 10.1.0.0
  (OpenJDK 25) · TruffleRuby 34.0.1 (GraalVM CE Native).
- **Method:** each process runs 3 untimed warm-up passes, then 25 timed passes of
  a fixed inner loop, timed with a monotonic clock; the **best** pass is reported
  as **ns/op** (lower is better). `vs MRI` < 1.00× means *faster than MRI*.
  Interpreter start-up is outside the timed region, so these are operation costs,
  not `ruby file.rb` process costs.

#### abbrev — `Abbrev.abbrev(words)` over a fixed 33-word list

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 11959.5 | 0.30× |
| MRI | 40097.5 | 1.00× |
| MRI + YJIT | 29772.0 | 0.74× |
| JRuby | 26713.0 | 0.67× |
| TruffleRuby | 5858.4 | 0.15× |

`Abbrev.abbrev` builds the unambiguous-abbreviation → word table over a fixed list
of 33 words drawn from several prefix-sharing families (`ruby`/`rules`/`ruler`/…,
`car`/`card`/`care`/`cargo`/…, `program`/`programmer`/…, `data`/`database`/…). The
pure-Go library **beats both MRI and MRI + YJIT**: **0.30× MRI** (~3.4× faster)
and **11959.5 / 29772.0 = 0.40× YJIT** (~2.5× faster than YJIT). MRI's `abbrev` is
pure Ruby (no C fast path), so the interpreter walks every prefix of every word,
allocating and hashing each — exactly where a compiled Go map loop pulls ahead;
YJIT narrows but does not close the gap. **TruffleRuby is fastest here** (0.15×
MRI): its Graal JIT compiles the pure-Ruby prefix loop to native code that, on
this steady-state warm loop, edges out the Go map churn. JRuby (0.67×) also beats
MRI once warm. The go-ruby column is the pure-Go library; every other column is
that interpreter's own `abbrev` stdlib doing the equivalent work; every output was
checked **byte-identical to MRI** against the live oracle before timing.

**go vs YJIT verdict:** go-ruby **beats YJIT** on this op (0.40× YJIT, ~2.5×
faster).

!!! note "Reproduce"
    The harness is committed under
    [`benchmarks/`](https://github.com/go-ruby-abbrev/docs/tree/main/benchmarks):
    a self-contained Go driver (`go/`, pins the published library via `go.mod`),
    the equivalent `ruby/abbrev.rb` workload, and `run.sh`. Run
    `bash benchmarks/run.sh`; env `OUTER`/`WARM` tune the pass budget and
    `RUBY`/`JRUBY`/`TRUFFLERUBY` select the runtime binaries.

!!! warning "Warm-up budget & noise — honest framing"
    Numbers reflect a **fixed warm-process budget** (3 warm-up + 25 timed passes
    in one process). The JVM/GraalVM JITs (JRuby, TruffleRuby) may need a larger
    warm-up to reach steady state, so their columns can shift with the budget;
    on this workload both had *already* reached a warm regime and beat MRI, but
    treat their exact ratios as budget-dependent. Every number here is a **real
    measured value** from the dated run above — nothing is fabricated, estimated,
    or cherry-picked.
