"""
OmniRoute AI Gateway - FastAPI service for AI/ML model orchestration
Provides unified access to LLM providers, embeddings, and AI capabilities
"""
from fastapi import FastAPI, HTTPException, Depends, BackgroundTasks
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel, Field
from typing import Optional, List, Dict, Any
from enum import Enum
import httpx
import asyncio
import os
import json
import logging
from datetime import datetime
import uuid

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = FastAPI(
    title="OmniRoute AI Gateway",
    description="Unified AI/ML orchestration for the OmniRoute ecosystem",
    version="1.0.0"
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# =============================================================================
# Models
# =============================================================================

class AIProvider(str, Enum):
    OPENAI = "openai"
    ANTHROPIC = "anthropic"
    GEMINI = "gemini"
    LOCAL = "local"

class MessageRole(str, Enum):
    USER = "user"
    ASSISTANT = "assistant"
    SYSTEM = "system"

class Message(BaseModel):
    role: MessageRole
    content: str

class CompletionRequest(BaseModel):
    provider: AIProvider = AIProvider.OPENAI
    model: str = "gpt-4"
    messages: List[Message]
    temperature: float = 0.7
    max_tokens: int = 2000
    stream: bool = False
    tools: Optional[List[Dict[str, Any]]] = None
    context: Optional[Dict[str, Any]] = None

class CompletionResponse(BaseModel):
    id: str
    provider: str
    model: str
    content: str
    usage: Dict[str, int]
    created_at: datetime
    latency_ms: int

class EmbeddingRequest(BaseModel):
    text: str | List[str]
    model: str = "text-embedding-3-small"
    provider: AIProvider = AIProvider.OPENAI

class EmbeddingResponse(BaseModel):
    embeddings: List[List[float]]
    model: str
    usage: Dict[str, int]

class AnalysisRequest(BaseModel):
    analysis_type: str  # sentiment, entity_extraction, classification, summarization
    text: str
    options: Optional[Dict[str, Any]] = None

class RecommendationRequest(BaseModel):
    customer_id: str
    context: Dict[str, Any]  # purchase_history, preferences, etc.
    limit: int = 10

class ForecastRequest(BaseModel):
    product_id: Optional[str] = None
    category: Optional[str] = None
    horizon_days: int = 30
    include_confidence: bool = True

# =============================================================================
# AI Provider Clients
# =============================================================================

class OpenAIClient:
    def __init__(self):
        self.api_key = os.getenv("OPENAI_API_KEY")
        self.base_url = "https://api.openai.com/v1"
    
    async def complete(self, request: CompletionRequest) -> CompletionResponse:
        start_time = datetime.now()
        
        async with httpx.AsyncClient() as client:
            response = await client.post(
                f"{self.base_url}/chat/completions",
                headers={"Authorization": f"Bearer {self.api_key}"},
                json={
                    "model": request.model,
                    "messages": [{"role": m.role, "content": m.content} for m in request.messages],
                    "temperature": request.temperature,
                    "max_tokens": request.max_tokens,
                },
                timeout=60.0
            )
            response.raise_for_status()
            data = response.json()
        
        latency = int((datetime.now() - start_time).total_seconds() * 1000)
        
        return CompletionResponse(
            id=data["id"],
            provider="openai",
            model=data["model"],
            content=data["choices"][0]["message"]["content"],
            usage=data["usage"],
            created_at=datetime.now(),
            latency_ms=latency
        )
    
    async def embed(self, request: EmbeddingRequest) -> EmbeddingResponse:
        texts = request.text if isinstance(request.text, list) else [request.text]
        
        async with httpx.AsyncClient() as client:
            response = await client.post(
                f"{self.base_url}/embeddings",
                headers={"Authorization": f"Bearer {self.api_key}"},
                json={"model": request.model, "input": texts},
                timeout=30.0
            )
            response.raise_for_status()
            data = response.json()
        
        return EmbeddingResponse(
            embeddings=[e["embedding"] for e in data["data"]],
            model=data["model"],
            usage=data["usage"]
        )

class AnthropicClient:
    def __init__(self):
        self.api_key = os.getenv("ANTHROPIC_API_KEY")
        self.base_url = "https://api.anthropic.com/v1"
    
    async def complete(self, request: CompletionRequest) -> CompletionResponse:
        start_time = datetime.now()
        
        # Convert messages for Claude format
        system_msg = next((m.content for m in request.messages if m.role == MessageRole.SYSTEM), None)
        messages = [{"role": m.role, "content": m.content} for m in request.messages if m.role != MessageRole.SYSTEM]
        
        async with httpx.AsyncClient() as client:
            body = {
                "model": request.model or "claude-3-sonnet-20240229",
                "messages": messages,
                "max_tokens": request.max_tokens,
            }
            if system_msg:
                body["system"] = system_msg
                
            response = await client.post(
                f"{self.base_url}/messages",
                headers={
                    "x-api-key": self.api_key,
                    "anthropic-version": "2023-06-01"
                },
                json=body,
                timeout=60.0
            )
            response.raise_for_status()
            data = response.json()
        
        latency = int((datetime.now() - start_time).total_seconds() * 1000)
        
        return CompletionResponse(
            id=data["id"],
            provider="anthropic",
            model=data["model"],
            content=data["content"][0]["text"],
            usage={"prompt_tokens": data["usage"]["input_tokens"], "completion_tokens": data["usage"]["output_tokens"]},
            created_at=datetime.now(),
            latency_ms=latency
        )

# Global clients
openai_client = OpenAIClient()
anthropic_client = AnthropicClient()

def get_client(provider: AIProvider):
    if provider == AIProvider.OPENAI:
        return openai_client
    elif provider == AIProvider.ANTHROPIC:
        return anthropic_client
    else:
        raise HTTPException(status_code=400, detail=f"Unsupported provider: {provider}")

# =============================================================================
# API Endpoints
# =============================================================================

@app.get("/health")
async def health():
    return {"status": "healthy", "service": "ai-gateway", "timestamp": datetime.now().isoformat()}

@app.post("/api/v1/completions", response_model=CompletionResponse)
async def create_completion(request: CompletionRequest):
    """Generate AI completion using specified provider"""
    try:
        client = get_client(request.provider)
        return await client.complete(request)
    except httpx.HTTPError as e:
        logger.error(f"AI provider error: {e}")
        raise HTTPException(status_code=502, detail=f"AI provider error: {str(e)}")

@app.post("/api/v1/embeddings", response_model=EmbeddingResponse)
async def create_embeddings(request: EmbeddingRequest):
    """Generate text embeddings"""
    try:
        client = get_client(request.provider)
        return await client.embed(request)
    except httpx.HTTPError as e:
        logger.error(f"Embedding error: {e}")
        raise HTTPException(status_code=502, detail=f"Embedding error: {str(e)}")

@app.post("/api/v1/analyze")
async def analyze_text(request: AnalysisRequest):
    """Perform text analysis (sentiment, entities, classification, summarization)"""
    
    prompts = {
        "sentiment": f"Analyze the sentiment of this text. Return JSON with 'sentiment' (positive/negative/neutral), 'confidence' (0-1), and 'explanation':\n\n{request.text}",
        "entity_extraction": f"Extract entities from this text. Return JSON with 'entities' array containing 'text', 'type' (person/organization/location/date/product), and 'confidence':\n\n{request.text}",
        "classification": f"Classify this text into categories. Return JSON with 'categories' array of 'name' and 'confidence':\n\n{request.text}",
        "summarization": f"Summarize this text in 2-3 sentences:\n\n{request.text}"
    }
    
    prompt = prompts.get(request.analysis_type)
    if not prompt:
        raise HTTPException(status_code=400, detail=f"Unknown analysis type: {request.analysis_type}")
    
    completion_request = CompletionRequest(
        provider=AIProvider.OPENAI,
        model="gpt-4",
        messages=[Message(role=MessageRole.USER, content=prompt)],
        temperature=0.3
    )
    
    result = await openai_client.complete(completion_request)
    
    # Try to parse as JSON, otherwise return raw content
    try:
        parsed = json.loads(result.content)
        return {"analysis_type": request.analysis_type, "result": parsed}
    except json.JSONDecodeError:
        return {"analysis_type": request.analysis_type, "result": result.content}

@app.post("/api/v1/recommendations")
async def get_recommendations(request: RecommendationRequest):
    """Generate product recommendations for a customer"""
    
    # Build context-aware prompt
    prompt = f"""Based on the following customer context, recommend {request.limit} products.

Customer ID: {request.customer_id}
Context: {json.dumps(request.context, indent=2)}

Return recommendations as JSON array with: product_id, name, reason, confidence_score (0-1)"""

    completion_request = CompletionRequest(
        provider=AIProvider.OPENAI,
        model="gpt-4",
        messages=[Message(role=MessageRole.USER, content=prompt)],
        temperature=0.5
    )
    
    result = await openai_client.complete(completion_request)
    
    try:
        recommendations = json.loads(result.content)
        return {"customer_id": request.customer_id, "recommendations": recommendations}
    except json.JSONDecodeError:
        return {"customer_id": request.customer_id, "recommendations": [], "raw_response": result.content}

@app.post("/api/v1/forecast")
async def generate_forecast(request: ForecastRequest):
    """Generate demand forecast (calls forecasting-service internally)"""
    
    # In production, this would call the forecasting-service
    # Here we return a mock response
    return {
        "product_id": request.product_id,
        "category": request.category,
        "horizon_days": request.horizon_days,
        "forecast": [
            {"date": f"2026-01-{19+i:02d}", "predicted_demand": 100 + i * 5, "confidence_lower": 90 + i * 4, "confidence_upper": 110 + i * 6}
            for i in range(min(request.horizon_days, 14))
        ]
    }

@app.post("/api/v1/chat/order-assistant")
async def order_assistant(messages: List[Message], customer_id: Optional[str] = None):
    """Specialized order assistant for OmniRoute"""
    
    system_prompt = """You are OmniRoute's order assistant. Help customers with:
- Placing new orders
- Checking order status
- Product recommendations
- Pricing inquiries
- Delivery tracking

Be helpful, concise, and professional. Use Nigerian Naira (₦) for prices."""

    full_messages = [Message(role=MessageRole.SYSTEM, content=system_prompt)] + messages
    
    completion_request = CompletionRequest(
        provider=AIProvider.OPENAI,
        model="gpt-4",
        messages=full_messages,
        temperature=0.7
    )
    
    result = await openai_client.complete(completion_request)
    return {"response": result.content, "customer_id": customer_id}

# =============================================================================
# Entry Point
# =============================================================================

if __name__ == "__main__":
    import uvicorn
    port = int(os.getenv("PORT", "8140"))
    uvicorn.run(app, host="0.0.0.0", port=port)
