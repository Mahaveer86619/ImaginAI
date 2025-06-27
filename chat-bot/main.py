# A simple FastAPI application to demonstrate the Python AI service
from fastapi import FastAPI

print("Starting the Python AI service...")

app = FastAPI()

@app.post("/health")
def health_check():
    return {"status": "healthy"}
