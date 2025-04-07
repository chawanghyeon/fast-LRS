FROM python:3.11-slim

WORKDIR /app

COPY app.py .

RUN pip install --no-cache-dir fastapi uvicorn psycopg2-binary pydantic[email]

CMD ["uvicorn", "app:app", "--host", "0.0.0.0", "--port", "8000", "--workers", "1", "--no-access-log"]
