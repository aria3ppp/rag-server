#!/bin/bash

# set -Eefuvx
# set -Eefux
set -Eefu

source "$(dirname "$0")/scripts/llama_utils.sh"

# embeddings
download CompendiumLabs/bge-small-en-v1.5-gguf bge-small-en-v1.5-f32.gguf
# download Ralriki/multilingual-e5-large-instruct-GGUF multilingual-e5-large-instruct-F16.gguf

# rerankers
# download gpustack/jina-reranker-v1-turbo-en-GGUF jina-reranker-v1-turbo-en-FP16.gguf
# download gpustack/jina-reranker-v2-base-multilingual-GGUF jina-reranker-v2-base-multilingual-FP16.gguf
# download gpustack/bge-reranker-v2-m3-GGUF bge-reranker-v2-m3-FP16.gguf

# llms
download hugging-quants/Llama-3.2-1B-Instruct-Q4_K_M-GGUF llama-3.2-1b-instruct-q4_k_m.gguf

# Download PyTorch models then convert and quantize them

MODEL_BGE_RERANKER_BASE_REPO=BAAI/bge-reranker-base
download "$MODEL_BGE_RERANKER_BASE_REPO" \
  pytorch_model.bin \
  config.json \
  tokenizer.json \
  tokenizer_config.json \
  special_tokens_map.json \
  sentencepiece.bpe.model

# Convert and quantize to Q4_K_M (or any other type you prefer)
convert_and_quantize "$MODEL_BGE_RERANKER_BASE_REPO" Q4_K_M

# MODEL_BGE_SMALL_EN_V1_5_REPO=BAAI/bge-small-en-v1.5
# download "$MODEL_BGE_SMALL_EN_V1_5_REPO"

# MODEL_LLAMA_3_1_8B_INSTRUCT_REPO=meta-llama/Meta-Llama-3.1-8B-Instruct
# download "$MODEL_LLAMA_3_1_8B_INSTRUCT_REPO" \
#     model-00001-of-00004.safetensors \
#     model-00002-of-00004.safetensors \
#     model-00003-of-00004.safetensors \
#     model-00004-of-00004.safetensors \
#     model.safetensors.index.json \
#     tokenizer.json \
#     tokenizer_config.json \
#     special_tokens_map.json \
#     config.json \
#     generation_config.json

# # Convert and quantize to multiple types
# convert_and_quantize "$MODEL_LLAMA_3_1_8B_INSTRUCT_REPO" F16 Q8_0 Q4_K_M