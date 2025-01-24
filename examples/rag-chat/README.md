### generate js client stubs

```
docker run --network=host --rm -v ${PWD}/examples/rag-chat/:/local openapitools/openapi-generator-cli generate \
  -i http://localhost:8000/v1/rag.swagger.json \
  -g javascript \
  -o /local/js-client \
  --additional-properties=usePromises=true,useES6=true
```

### put `./js-client` beside the sample index.html
```
❯ tree ./examples/rag-chat/
.
├── index.html
└── js-client
    ├── ...
    │── ...
    └── ...
```

### install serve tool
```
npm install -g serve
```

### serve index.html
```
serve -l 3000 ./examples/rag-chat/
```

### open chat client
```
browse http://localhost:3000
```