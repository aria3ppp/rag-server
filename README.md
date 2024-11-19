### run vector store manually
```shell
env $(grep -v '^#' .env.example | xargs) go run ./cmd/vectorstore 
```