#!/usr/bin/env bash

set -euo pipefail

nix build .#assets
(cd www && git clean -fxd)
cp result/www/* www/
