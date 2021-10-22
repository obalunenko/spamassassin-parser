#!/bin/sh

set -e

echo "current user $(whoami)"

./spamassassin-parser \
  --input="${SPAMASSASSIN_INPUT}" \
  --result="${SPAMASSASSIN_RESULT}" \
  --archive="${SPAMASSASSIN_ARCHIVE}" \
  --errors="${SPAMASSASSIN_RECEIVE_ERRORS}" \
  --log_level="${SPAMASSASSIN_LOG_LEVEL}" \
  --log_format="${SPAMASSASSIN_LOG_FORMAT}" \
  --log_sentry_dsn="${SPAMASSASSIN_LOG_SENTRY_DSN}" \
  --log_sentry_trace="${SPAMASSASSIN_LOG_SENTRY_TRACE}" \
  --log_sentry_trace_level="${SPAMASSASSIN_LOG_SENTRY_TRACE_LEVEL}"
