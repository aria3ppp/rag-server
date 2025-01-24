#!/bin/bash

# Function to display usage
function usage {
  echo "Usage: $0 [--hf-token TOKEN]"
  echo "  --hf-token TOKEN: Hugging Face token for authentication (optional)."
  echo "  -h, --help: Show this help message and exit."
  echo "  Alternatively, set the HUGGINGFACE_TOKEN environment variable."
}

# Initialize HUGGINGFACE_TOKEN with the environment variable (or default to empty if unset)
HUGGINGFACE_TOKEN="${HUGGINGFACE_TOKEN:-}"

# Parse command-line arguments
while [[ $# -gt 0 ]]; do
  case "$1" in
    --hf-token)
      HUGGINGFACE_TOKEN="$2"  # Override with the command-line token
      shift 2
      ;;
    -h | --help)
      usage
      exit 0
      ;;
    *)
      echo "Unknown argument: $1"
      usage
      exit 1
      ;;
  esac
done

# Get the directory of this script (llama_utils.sh)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"

# Resolve the parent directory of the scripts directory
PARENT_DIR="$(dirname "$SCRIPT_DIR")"

# Set paths relative to the parent directory
MODELS_DIR="$PARENT_DIR/models"
LLAMA_CPP_DIR="$MODELS_DIR/.llama.cpp"

# Read the llama.cpp commit hash from .llama.cpp-version
LLAMA_CPP_VERSION=$(cat "$PARENT_DIR/.llama.cpp-version")

# Clone llama.cpp if not already done
if [ ! -d "$LLAMA_CPP_DIR" ]; then
  echo "Cloning llama.cpp..."
  git clone --depth 1 https://github.com/ggerganov/llama.cpp.git "$LLAMA_CPP_DIR"
  git -C "$LLAMA_CPP_DIR" fetch --unshallow
else
  echo "llama.cpp already cloned. Skipping clone."
fi

# Checkout llama.cpp to the specific commit
git -C "$LLAMA_CPP_DIR" checkout "$LLAMA_CPP_VERSION"

# Create a Python virtual environment
echo "Creating Python virtual environment..."
python3 -m venv "$LLAMA_CPP_DIR/.venv"
source "$LLAMA_CPP_DIR/.venv/bin/activate"

# Install dependencies from requirements.txt
echo "Installing Python dependencies..."
pip install -r "$LLAMA_CPP_DIR/requirements.txt"

# Check if Ninja is installed
if ! command -v ninja &> /dev/null; then
  echo "Error: Ninja is not installed. Please install Ninja before running this script."
  exit 1
fi

# Build llama.cpp using CMake with Ninja
echo "Building llama.cpp..."
mkdir -p "$LLAMA_CPP_DIR/build"
cmake -G Ninja -S "$LLAMA_CPP_DIR" -B "$LLAMA_CPP_DIR/build"
cmake --build "$LLAMA_CPP_DIR/build" --config Release

# Authenticate with Hugging Face if a token is provided
if [ -n "$HUGGINGFACE_TOKEN" ]; then
  echo "Authenticating with Hugging Face using provided token..."
  huggingface-cli login --token "$HUGGINGFACE_TOKEN"
fi

function download {
  local MODEL_REPO=$1
  shift  # Shift to get the list of files
  local FILES=("$@")  # Remaining arguments are the files to download

  echo "Downloading files from $MODEL_REPO..."
  local MODEL_DIR="$MODELS_DIR/$MODEL_REPO"
  huggingface-cli download "$MODEL_REPO" "${FILES[@]}" --repo-type model --local-dir "$MODEL_DIR"
}

function convert_and_quantize {
  local MODEL_REPO=$1
  shift  # Shift to get the list of quantization types
  local QUANT_TYPES=("$@")  # Remaining arguments are the quantization types

  local MODEL_DIR="$MODELS_DIR/$MODEL_REPO"
  local MODEL_NAME=$(basename "$MODEL_REPO")

  # Step 1: Check if any quantized files already exist
  local SKIP_CONVERSION=true
  for QUANT_TYPE in "${QUANT_TYPES[@]}"; do
    local OUTPUT_FILE="$MODEL_DIR/$MODEL_NAME-$QUANT_TYPE.gguf"
    if [ ! -f "$OUTPUT_FILE" ]; then
      SKIP_CONVERSION=false  # At least one quantized file is missing
      break
    fi
  done

  # Step 2: Convert to GGUF only if no quantized files exist
  local HIGH_PRECISION_FILE="$MODEL_DIR/$MODEL_NAME-auto.gguf"
  if [ "$SKIP_CONVERSION" = false ]; then
    if [ ! -f "$HIGH_PRECISION_FILE" ]; then
      echo "Converting $MODEL_REPO to GGUF (auto)..."
      python3 "$LLAMA_CPP_DIR/convert_hf_to_gguf.py" "$MODEL_DIR" --outtype auto --outfile "$HIGH_PRECISION_FILE"
    else
      echo "High-precision GGUF file already exists: $HIGH_PRECISION_FILE"
    fi
  else
    echo "All quantized files already exist. Skipping conversion."
  fi

  # Step 3: Quantize to each specified type
  for QUANT_TYPE in "${QUANT_TYPES[@]}"; do
    local OUTPUT_FILE="$MODEL_DIR/$MODEL_NAME-$QUANT_TYPE.gguf"
    if [ ! -f "$OUTPUT_FILE" ]; then
      echo "Quantizing $HIGH_PRECISION_FILE to $QUANT_TYPE..."
      "$LLAMA_CPP_DIR/build/bin/llama-quantize" "$HIGH_PRECISION_FILE" "$OUTPUT_FILE" "$QUANT_TYPE"
    else
      echo "File $OUTPUT_FILE already exists. Skipping quantization to $QUANT_TYPE."
    fi
  done

  # Step 4: Delete the intermediate high-precision file
  if [ -f "$HIGH_PRECISION_FILE" ]; then
    echo "Deleting intermediate high-precision file: $HIGH_PRECISION_FILE"
    rm "$HIGH_PRECISION_FILE"
  fi
}