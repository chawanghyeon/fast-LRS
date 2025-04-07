from fastapi import FastAPI, HTTPException
from pydantic import BaseModel, EmailStr, Field, validator
from typing import List
import uuid
import re
import psycopg2

app = FastAPI()

class User(BaseModel):
    name: str = Field(..., min_length=2, max_length=50)
    email: EmailStr
    age: int = Field(..., ge=0, le=120)
    bio: str = Field(..., min_length=10)
    interests: List[str] = Field(..., min_items=1)

    @validator("email")
    def validate_email_domain(cls, v):
        if not v.endswith("@example.com"):
            raise ValueError("Only @example.com emails are allowed")
        return v

    @validator("bio")
    def validate_bio_keywords(cls, v):
        if not re.search(r"(engineer|developer|programmer)", v.lower()):
            raise ValueError("Bio must include your profession")
        return v

conn = psycopg2.connect(
    host="postgres",
    dbname="benchmark",
    user="benchmark",
    password="benchmark"
)
conn.autocommit = True

@app.post("/data")
async def create_user(user: User):
    user_id = str(uuid.uuid4())
    with conn.cursor() as cur:
        cur.execute(
            "INSERT INTO users (id, name, email) VALUES (%s, %s, %s)",
            (user_id, user.name, user.email)
        )
    return {"id": user_id}
