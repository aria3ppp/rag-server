# RAG Server Deployment Guide

## Table of Contents
- [How to Run the Server](#how-to-run-the-server)
  - [Install Models](#install-models)
  - [Run RAG Server via Docker](#run-rag-server-via-docker)
  - [Populate Vector Store](#populate-vectorstore)
  - [Test the RAG Server](#test-the-rag-server)

## How to Run the Server

### Install Models

#### Prerequisites
- Ensure `ninja` tool is on your PATH  
  **Ubuntu example:**
  ```bash
  sudo apt install ninja-build
  ```

#### Installation Steps
1. Run the installation script:
   ```bash
   ./install-models.sh
   ```

**Important Notes:**
- Modify `install-models.sh` to specify which models to install
- For models with HuggingFace restrictions:
  - Environment variable method:
    ```bash
    HUGGINGFACE_TOKEN=hf_tokenblahblahblah ./install-models.sh
    ```
  - CLI argument method (overrides environment variable):
    ```bash
    ./install-models.sh --hf-token hf_tokenblahblahblah
    ```

### Run RAG Server via Docker
```bash
docker compose up --build -d --wait rag
```

### Populate Vectorstore
Follow the [vectorstore population guide](examples/populate-vectorstore/README.md) to load online article content.

#### Python Environment Setup
```bash
python3 -m venv ./examples/populate-vectorstore/.venv
source ./examples/populate-vectorstore/.venv/bin/activate
pip install --requirement=./examples/populate-vectorstore/requirements.txt
```

#### Content Population Examples
**Artificial Intelligence:**
```bash
python3 ./examples/populate-vectorstore/populate-vectorstore.py \
  --max-chunk-bytes 2000 \
  https://en.wikipedia.org/wiki/Artificial_intelligence
```

**Cyrus the Great:**
```bash
python3 ./examples/populate-vectorstore/populate-vectorstore.py \
  --max-chunk-bytes 2000 \
  https://en.wikipedia.org/wiki/Cyrus_the_Great
```

### Test the RAG Server
Follow the [chat client guide](examples/rag-chat/README.md) to deploy a test client.

#### Client Setup
1. Generate JS client stubs:
   ```bash
   docker run --network=host --rm -v ${PWD}/examples/rag-chat/:/local \
     openapitools/openapi-generator-cli generate \
     -i http://localhost:8000/v1/rag.swagger.json \
     -g javascript \
     -o /local/js-client \
     --additional-properties=usePromises=true,useES6=true
   ```

2. Verify directory structure:
   ```
   ./examples/rag-chat/
   ├── index.html
   └── js-client
       ├── ...
       └── ...
   ```

3. Install serve tool:
   ```bash
   npm install -g serve
   ```

4. Start web server:
   ```bash
   serve -l 3000 ./examples/rag-chat/
   ```

5. Access chat client:
   ```
   http://localhost:3000
   ```