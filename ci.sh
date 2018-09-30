#!/bin/sh
set -e

export BINDIR="$(mktemp -d)"
curl -sSf https://raw.githubusercontent.com/bloom42/rocket/master/install.sh | sh

"${BINDIR}/rocket" "$@"
