#!/bin/bash

one() {
    ./bin/tao < "$name" | diff -q - "${name%%.*}.out"
}

for name in tests/*.tao; do
    one
done
