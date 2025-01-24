import argparse
import requests
from bs4 import BeautifulSoup
import nltk
from urllib.parse import urlparse

# Ensure NLTK punkt tokenizer is available
try:
    nltk.data.find('tokenizers/punkt')
except LookupError:
    nltk.download('punkt')

def get_webpage_content(url):
    """Fetch and parse content from any webpage"""
    try:
        response = requests.get(url, timeout=10)
        response.raise_for_status()
    except requests.exceptions.RequestException as e:
        print(f"Error fetching {url}: {e}")
        return []
    
    soup = BeautifulSoup(response.text, 'html.parser')
    paragraphs = []
    
    # Try to find main content areas first
    main_content = soup.find(['main', 'article']) or soup.find(id='content')
    content_source = main_content if main_content else soup
    
    for p in content_source.find_all('p'):
        text = p.get_text().strip()
        if text:
            paragraphs.append(text)
    
    return paragraphs

def split_into_chunks(text, max_bytes):
    """Split text into chunks that don't exceed max_bytes"""
    chunks = []
    current_chunk = []
    current_bytes = 0
    
    for word in text.split():
        word_bytes = len(word.encode('utf-8'))
        space_bytes = 1 if current_chunk else 0  # Account for space between words
        
        if current_bytes + word_bytes + space_bytes > max_bytes:
            if current_chunk:
                chunks.append(" ".join(current_chunk))
                current_chunk = []
                current_bytes = 0
            
            # Handle words larger than max_bytes
            while word_bytes > max_bytes:
                chunk_part = word[:max_bytes//4]  # Simple split for demonstration
                chunks.append(chunk_part)
                word = word[len(chunk_part):]
                word_bytes = len(word.encode('utf-8'))
            
            if word:
                current_chunk.append(word)
                current_bytes = word_bytes
        else:
            current_chunk.append(word)
            current_bytes += word_bytes + space_bytes
    
    if current_chunk:
        chunks.append(" ".join(current_chunk))
    
    return chunks

def chunk_text(paragraphs, max_bytes=2000):
    """Split paragraphs into byte-limited chunks"""
    chunks = []
    current_chunk = []
    current_bytes = 0
    
    def add_chunk(text, text_bytes):
        """Helper to manage chunk accumulation"""
        nonlocal current_bytes, current_chunk
        if current_bytes + text_bytes > max_bytes:
            if current_chunk:
                chunks.append(" ".join(current_chunk))
                current_chunk = []
                current_bytes = 0
        current_chunk.append(text)
        current_bytes += text_bytes + 1  # Account for space
    
    for para in paragraphs:
        para_bytes = len(para.encode('utf-8'))
        
        if para_bytes == 0:
            continue
        
        if para_bytes > max_bytes:
            # Split large paragraphs into sentences first
            sentences = nltk.sent_tokenize(para)
            for sentence in sentences:
                sent_bytes = len(sentence.encode('utf-8'))
                
                if sent_bytes > max_bytes:
                    # Split into smaller chunks
                    sub_chunks = split_into_chunks(sentence, max_bytes)
                    for chunk in sub_chunks:
                        add_chunk(chunk, len(chunk.encode('utf-8')))
                else:
                    add_chunk(sentence, sent_bytes)
        else:
            add_chunk(para, para_bytes)
    
    # Add remaining content
    if current_chunk:
        chunks.append(" ".join(current_chunk))
    
    return chunks

def insert_chunks(chunks, url, api_url):
    """Insert chunks into vector store with metadata"""
    parsed_url = urlparse(url)
    payload = {
        "texts": [{
            "text": chunk,
            "metadata": {
                "source": parsed_url.netloc,
                "path": parsed_url.path,
                "chunk_id": idx
            }
        } for idx, chunk in enumerate(chunks)]
    }
    
    try:
        response = requests.post(api_url, json=payload, timeout=30)
        response.raise_for_status()
        return response.json()
    except requests.exceptions.RequestException as e:
        print(f"Error inserting chunks: {e}")
        return {"error": str(e)}

def main():
    parser = argparse.ArgumentParser(description="Webpage content to vector store processor")
    parser.add_argument("url", help="URL of the webpage to process")
    parser.add_argument("--max-chunk-bytes", type=int, default=2000, 
                       help="Maximum chunk size in bytes (default: 2000)")
    parser.add_argument("--api-url", default="http://localhost:8080/api/v1/insert_texts",
                       help="Vector store API endpoint (default: http://localhost:8080/api/v1/insert_texts)")
    args = parser.parse_args()

    print(f"Processing: {args.url}")
    paragraphs = get_webpage_content(args.url)
    
    if not paragraphs:
        print("No content found")
        return
    
    chunks = chunk_text(paragraphs, args.max_chunk_bytes)
    print(f"Created {len(chunks)} chunks")
    
    # Validate chunk sizes
    for i, chunk in enumerate(chunks):
        chunk_bytes = len(chunk.encode('utf-8'))
        if chunk_bytes > args.max_chunk_bytes:
            print(f"Warning: Chunk {i} exceeds limit ({chunk_bytes}/{args.max_chunk_bytes} bytes)")
    
    result = insert_chunks(chunks, args.url, args.api_url)
    print("Insertion result:", result)

if __name__ == "__main__":
    main()