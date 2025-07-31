from fastapi import FastAPI
from fastapi.responses import JSONResponse
import uvicorn
import os

# Create FastAPI instance
app = FastAPI(
    title="Hello World API",
    description="A simple FastAPI Hello World application",
    version="1.0.0",
    docs_url="/docs",
    redoc_url="/redoc"
)

# Health check endpoint
@app.get("/health")
async def health_check():
    """Health check endpoint for load balancers and monitoring"""
    return {"status": "healthy", "message": "Service is running"}

# Root endpoint
@app.get("/")
async def read_root():
    """Root endpoint returning a welcome message"""
    return {"message": "Hello World!", "version": "1.0.0"}

# Hello endpoint with path parameter
@app.get("/hello/{name}")
async def say_hello(name: str):
    """Personalized hello message"""
    return {"message": f"Hello, {name}!", "name": name}

# Hello endpoint with query parameter
@app.get("/hello")
async def say_hello_query(name: str = "World"):
    """Hello message with optional query parameter"""
    return {"message": f"Hello, {name}!", "name": name}

if __name__ == "__main__":
    # This allows running the app directly with python main.py
    port = int(os.getenv("PORT", 8000))
    uvicorn.run(
        "main:app",
        host="0.0.0.0",
        port=port,
        reload=True if os.getenv("ENVIRONMENT") == "development" else False
    )
