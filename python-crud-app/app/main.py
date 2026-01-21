"""
FastAPI application entry point.
"""
from datetime import datetime
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from app.core.config import settings
from app.core.database import init_db
from app.middleware.ratelimit import RateLimiter, RateLimitMiddleware
from app.routers import auth, users


# Initialize database
init_db()

# Create FastAPI application
app = FastAPI(
    title="Python CRUD API",
    description="A production-ready RESTful API with authentication and CRUD operations",
    version="1.0.0",
    docs_url="/docs",
    redoc_url="/redoc"
)

# CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=settings.cors_origins_list,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Rate limiting middleware
rate_limiters = {
    "/api/auth/register": RateLimiter(requests=3, window=60),  # 3 requests per minute
    "/api/auth/login": RateLimiter(requests=5, window=60),     # 5 requests per minute
    "/api/users": RateLimiter(requests=100, window=60),        # 100 requests per minute
}

app.add_middleware(RateLimitMiddleware, rate_limiters=rate_limiters)

# Include routers
app.include_router(auth.router, prefix="/api")
app.include_router(users.router, prefix="/api")


@app.get("/health")
def health_check():
    """
    Health check endpoint.
    
    Returns:
        Health status and timestamp
    """
    return {
        "status": "healthy",
        "time": datetime.utcnow().isoformat() + "Z"
    }


@app.get("/")
def root():
    """
    Root endpoint.
    
    Returns:
        Welcome message
    """
    return {
        "message": "Welcome to Python CRUD API",
        "docs": "/docs",
        "health": "/health"
    }


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(
        "app.main:app",
        host="0.0.0.0",
        port=settings.PORT,
        reload=settings.ENV == "development"
    )
