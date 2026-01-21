"""
Authentication routes for user registration and login.
"""
from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session

from app.core.database import get_db
from app.core.security import hash_password, verify_password, create_access_token, PasswordError
from app.models.user import User
from app.schemas.user import UserCreate, UserLogin, AuthResponse, UserResponse


router = APIRouter(prefix="/auth", tags=["Authentication"])


@router.post("/register", response_model=AuthResponse, status_code=status.HTTP_201_CREATED)
def register(user_data: UserCreate, db: Session = Depends(get_db)):
    """
    Register a new user.
    
    Args:
        user_data: User registration data
        db: Database session
        
    Returns:
        Authentication response with token and user data
        
    Raises:
        HTTPException: If user already exists or password is weak
    """
    # Check if user already exists
    existing_user = db.query(User).filter(
        (User.email == user_data.email.lower()) | (User.username == user_data.username)
    ).first()
    
    if existing_user:
        raise HTTPException(
            status_code=status.HTTP_409_CONFLICT,
            detail="User with this email or username already exists"
        )
    
    # Hash password
    try:
        password_hash = hash_password(user_data.password)
    except PasswordError as e:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail=str(e)
        )
    
    # Create user
    user = User(
        username=user_data.username,
        email=user_data.email.lower(),
        password_hash=password_hash
    )
    
    db.add(user)
    db.commit()
    db.refresh(user)
    
    # Generate token (sub must be a string for python-jose)
    token = create_access_token(data={"sub": str(user.id)})
    
    return AuthResponse(
        token=token,
        user=UserResponse.model_validate(user)
    )


@router.post("/login", response_model=AuthResponse)
def login(credentials: UserLogin, db: Session = Depends(get_db)):
    """
    Login user and return JWT token.
    
    Args:
        credentials: User login credentials
        db: Database session
        
    Returns:
        Authentication response with token and user data
        
    Raises:
        HTTPException: If credentials are invalid
    """
    # Find user by email
    user = db.query(User).filter(User.email == credentials.email.lower()).first()
    
    if not user or not verify_password(credentials.password, user.password_hash):
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid email or password"
        )
    
    # Generate token (sub must be a string for python-jose)
    token = create_access_token(data={"sub": str(user.id)})
    
    return AuthResponse(
        token=token,
        user=UserResponse.model_validate(user)
    )
