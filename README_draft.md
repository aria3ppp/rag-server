RAG server in go

api is a gRPC server proxied by envoy to support REST http api
uses sentence-embedding* to embed documents (sentences) to a vector database (usally weaviate db)
accept the user query questions and retrived the context of the query question from vector db to add context to llm
the llm is a llama3.1 8b

https://go.dev/blog/llmpowered

* sentence-embedding aka text-embedding or embedding model:
 - https://stackoverflow.blog/2023/11/09/an-intuitive-introduction-to-text-embeddings/
 - https://www.cloudflare.com/learning/ai/what-are-embeddings/

-----------------------------------------------------------------------------------------------------------------------

https://claude.ai/chat/8fc391fc-dfc0-451d-a5a4-6e2f7d301711

This example demonstrates a simple RAG server that indexes documents, performs similarity search, and generates responses using OpenAI's API.
To further improve this implementation, you could consider:
    1.Adding batch processing for document indexing
    2. Implementing caching for frequently accessed documents or queries
    3. Using more advanced retrieval techniques, such as re-ranking or hybrid search
    4. Implementing error handling and retries for API calls
    5. Adding a web server interface for easier interaction

-----------------------------------------------------------------------------------------------------------------------

Features:

1. Sign in / Sign up
2. simple html interface for user in order to interact with the app
3. Payment gateway to buy subscription