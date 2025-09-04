#! /bin/bash

if [ -z "$VIRTUAL_ENV" ]; then
  echo "Error: No virtual environment found. Please activate a python 3.12 uv virtual environment."
  exit 1
fi

git clone https://github.com/vllm-project/vllm.git
cd vllm
uv pip install -r requirements/cpu.txt --extra-index-url https://pypi.org/simple/
uv pip install -e .

vllm serve Qwen/Qwen3-0.6B \
  --max-model-len 8192 \
  --max-num-batched-tokens 8192 \
  --guided-decoding-backend guidance \
  --reasoning-parser deepseek_r1 \
  --enable-auto-tool-choice \
  --tool-call-parser hermes
