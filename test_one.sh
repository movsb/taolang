#!/bin/bash

set -e

make build > /dev/null

./tao "$@"
