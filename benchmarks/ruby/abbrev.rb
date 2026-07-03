# frozen_string_literal: true
# SPDX-License-Identifier: BSD-3-Clause
require "abbrev"
require_relative "_harness"

# Fixed representative word list, identical to the Go driver: a few dozen words
# in several families that share prefixes, so the unambiguous-abbreviation table
# is non-trivial.
WORDS = %w[
  ruby rules ruler rubygems rubocop
  cat car card care cargo carpet
  go golang gopher gone good goose
  test testing tester tested testify
  program programmer programming progress
  data database datacenter date dating
].freeze

# canon serializes an abbrev table to a single deterministic line so the Go and
# Ruby drivers can be checked byte-identical before timing.
def canon(hash)
  hash.keys.sort.map { |k| "#{k}=#{hash[k]}" }.join(";")
end

# --dump: print the canonical abbrev table for the cross-runtime output check,
# then exit without timing anything.
if ARGV[0] == "--dump"
  puts canon(Abbrev.abbrev(WORDS))
  exit
end

bench("abbrev", 2000) { Abbrev.abbrev(WORDS) }
