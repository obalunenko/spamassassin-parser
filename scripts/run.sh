#!/usr/bin/env sh

./spamassassin-parser \
  --input_dir="${SPAMASSASSIN_INPUT}" \
  --output_dir="${SPAMASSASSIN_OUTPUT}" \
  --processed_dir="${SPAMASSASSIN_ARCHIVE}"