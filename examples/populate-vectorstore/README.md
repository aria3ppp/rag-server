### init python environment
```
python3 -m venv ./examples/populate-vectorstore/.venv
```

### source python environment
```
source ./examples/populate-vectorstore/.venv/bin/activate
```

### install dependencies
```
pip install --requirement=./examples/populate-vectorstore/requirements.txt
```

### populate vecotrestore with chunks of maximum 2000 bytes from webpage url

#### Artificial Intelligence
```
python3 ./examples/populate-vectorstore/populate-vectorstore.py --max-chunk-bytes 2000 https://en.wikipedia.org/wiki/Artificial_intelligence
```

#### Cyrus the Great
```
python3 ./examples/populate-vectorstore/populate-vectorstore.py --max-chunk-bytes 2000 https://en.wikipedia.org/wiki/Cyrus_the_Great
```
