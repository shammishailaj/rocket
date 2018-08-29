#!/bin/sh
set -e

export BINDIR="$(mktemp -d)"
curl -sSf https://raw.githubusercontent.com/astrocorp42/rocket/master/install.sh | sh

"${BINDIR}/rocket" "$@"
