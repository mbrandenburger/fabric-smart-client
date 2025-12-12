#!/usr/bin/env bash
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

set -euo pipefail

GOGC=100 go test -bench=. -benchmem -count=10 -timeout=20m -cpu=1,2,4,6,8,10,12,14,16,18,20,22,24,26,28,30,32,34,36,38,40,42,48,64 -run=^$ ./... > plots/benchmark_gc_100.txt
GOGC=8000 go test -bench=. -benchmem -count=10 -timeout=20m -cpu=1,2,4,6,8,10,12,14,16,18,20,22,24,26,28,30,32,34,36,38,40,42,48,64 -run=^$ ./... > plots/benchmark_gc_8000.txt

