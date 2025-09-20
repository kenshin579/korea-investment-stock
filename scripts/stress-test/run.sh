#!/usr/bin/env bash

echo "running..."
source .venv/bin/activate && python scripts/stress-test/stress_runner.py --config config.yaml --duration 5 --bulk-size 3