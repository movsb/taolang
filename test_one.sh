#!/bin/bash

set -e

make tao > /dev/null

./bin/tao "$@"

