#!/bin/bash

for name in *.tao; do
    ../bin/tao < "$name" | diff -q - "${name%%.*}.out"
done
