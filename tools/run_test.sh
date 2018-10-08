#!/bin/bash

root="$(dirname "$0")/.."

for name in *.tao; do
    "$root"/bin/tao < "$name" | diff -q - "${name%.*}.out"
done
