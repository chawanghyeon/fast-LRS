#!/bin/sh

echo "=== Running K6 Benchmarks ==="

for target in FastAPI Go NodeJs; do
  echo ""
  echo "--- Testing $target ---"
  TARGET=$target k6 run benchmark-all.js
done

echo ""
echo "✅ All benchmarks complete!"
