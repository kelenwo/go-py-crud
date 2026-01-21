"""
Tests for security utilities (password hashing and JWT).
"""
import pytest
from datetime import timedelta

from app.core.security import (
    hash_password,
    verify_password,
    validate_password,
    create_access_token,
    decode_access_token,
    PasswordError
)


class TestPasswordHashing:
    """Test password hashing functionality."""
    
    def test_hash_password_valid(self):
        """Test hashing a valid password."""
        password = "StrongPass123"
        hashed = hash_password(password)
        
        assert hashed is not None
        assert hashed != password
        assert len(hashed) > 0
    
    def test_hash_password_weak_too_short(self):
        """Test hashing a password that's too short."""
        with pytest.raises(PasswordError):
            hash_password("Short1")
    
    def test_hash_password_weak_no_uppercase(self):
        """Test hashing a password without uppercase."""
        with pytest.raises(PasswordError):
            hash_password("weakpass123")
    
    def test_hash_password_weak_no_lowercase(self):
        """Test hashing a password without lowercase."""
        with pytest.raises(PasswordError):
            hash_password("WEAKPASS123")
    
    def test_hash_password_weak_no_number(self):
        """Test hashing a password without number."""
        with pytest.raises(PasswordError):
            hash_password("WeakPassword")
    
    def test_verify_password_correct(self):
        """Test verifying correct password."""
        password = "TestPassword123"
        hashed = hash_password(password)
        
        assert verify_password(password, hashed) is True
    
    def test_verify_password_incorrect(self):
        """Test verifying incorrect password."""
        password = "TestPassword123"
        hashed = hash_password(password)
        
        assert verify_password("WrongPassword123", hashed) is False
    
    def test_verify_password_empty(self):
        """Test verifying empty password."""
        password = "TestPassword123"
        hashed = hash_password(password)
        
        assert verify_password("", hashed) is False


class TestPasswordValidation:
    """Test password validation."""
    
    def test_validate_password_valid(self):
        """Test validating a valid password."""
        validate_password("ValidPass123")  # Should not raise
    
    def test_validate_password_complex(self):
        """Test validating a complex password."""
        validate_password("C0mpl3xP@ssw0rd!")  # Should not raise
    
    def test_validate_password_too_short(self):
        """Test validating a password that's too short."""
        with pytest.raises(PasswordError):
            validate_password("Short1")
    
    def test_validate_password_no_uppercase(self):
        """Test validating a password without uppercase."""
        with pytest.raises(PasswordError):
            validate_password("lowercase123")
    
    def test_validate_password_no_lowercase(self):
        """Test validating a password without lowercase."""
        with pytest.raises(PasswordError):
            validate_password("UPPERCASE123")
    
    def test_validate_password_no_number(self):
        """Test validating a password without number."""
        with pytest.raises(PasswordError):
            validate_password("NoNumberPass")


class TestJWT:
    """Test JWT token functionality."""
    
    def test_create_access_token(self):
        """Test creating an access token."""
        data = {"sub": 1}
        token = create_access_token(data)
        
        assert token is not None
        assert isinstance(token, str)
        assert len(token) > 0
    
    def test_create_access_token_with_expiration(self):
        """Test creating a token with custom expiration."""
        data = {"sub": 1}
        token = create_access_token(data, expires_delta=timedelta(hours=1))
        
        assert token is not None
        assert isinstance(token, str)
    
    def test_decode_access_token_valid(self):
        """Test decoding a valid token."""
        data = {"sub": 1}
        token = create_access_token(data)
        
        decoded = decode_access_token(token)
        
        assert decoded is not None
        assert decoded["sub"] == 1
        assert "exp" in decoded
        assert "iat" in decoded
    
    def test_decode_access_token_invalid(self):
        """Test decoding an invalid token."""
        decoded = decode_access_token("invalid.token.here")
        
        assert decoded is None
    
    def test_decode_access_token_empty(self):
        """Test decoding an empty token."""
        decoded = decode_access_token("")
        
        assert decoded is None
    
    def test_decode_access_token_expired(self):
        """Test decoding an expired token."""
        data = {"sub": 1}
        token = create_access_token(data, expires_delta=timedelta(seconds=-1))
        
        decoded = decode_access_token(token)
        
        assert decoded is None
