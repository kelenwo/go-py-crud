"""
Rate limiting middleware to prevent brute force attacks.
"""
import time
from collections import defaultdict
from typing import Dict, List
from fastapi import Request, HTTPException, status
from starlette.middleware.base import BaseHTTPMiddleware


class RateLimiter:
    """Simple in-memory rate limiter."""
    
    def __init__(self, requests: int, window: int):
        """
        Initialize rate limiter.
        
        Args:
            requests: Maximum number of requests allowed
            window: Time window in seconds
        """
        self.requests = requests
        self.window = window
        self.clients: Dict[str, List[float]] = defaultdict(list)
    
    def is_allowed(self, client_id: str) -> bool:
        """
        Check if client is allowed to make a request.
        
        Args:
            client_id: Client identifier (usually IP address)
            
        Returns:
            True if allowed, False if rate limit exceeded
        """
        now = time.time()
        
        # Remove old requests outside the window
        self.clients[client_id] = [
            req_time for req_time in self.clients[client_id]
            if now - req_time < self.window
        ]
        
        # Check if limit exceeded
        if len(self.clients[client_id]) >= self.requests:
            return False
        
        # Add current request
        self.clients[client_id].append(now)
        return True
    
    def cleanup(self):
        """Clean up old entries."""
        now = time.time()
        for client_id in list(self.clients.keys()):
            self.clients[client_id] = [
                req_time for req_time in self.clients[client_id]
                if now - req_time < self.window
            ]
            if not self.clients[client_id]:
                del self.clients[client_id]


class RateLimitMiddleware(BaseHTTPMiddleware):
    """Middleware for rate limiting."""
    
    def __init__(self, app, rate_limiters: Dict[str, RateLimiter]):
        """
        Initialize middleware.
        
        Args:
            app: FastAPI application
            rate_limiters: Dictionary mapping path patterns to rate limiters
        """
        super().__init__(app)
        self.rate_limiters = rate_limiters
    
    async def dispatch(self, request: Request, call_next):
        """Process request with rate limiting."""
        client_ip = request.client.host if request.client else "unknown"
        path = request.url.path
        
        # Find matching rate limiter
        limiter = None
        for pattern, rate_limiter in self.rate_limiters.items():
            if pattern in path:
                limiter = rate_limiter
                break
        
        # Apply rate limiting if limiter found
        if limiter and not limiter.is_allowed(client_ip):
            raise HTTPException(
                status_code=status.HTTP_429_TOO_MANY_REQUESTS,
                detail="Rate limit exceeded. Please try again later."
            )
        
        response = await call_next(request)
        return response
