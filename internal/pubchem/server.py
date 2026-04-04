import argparse
from contextlib import asynccontextmanager
from fastapi import FastAPI, HTTPException, Depends
from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession
from sqlalchemy.orm import sessionmaker
from sqlalchemy import text
import uvicorn

parser = argparse.ArgumentParser(description="Ingredient Server")
parser.add_argument("--port", type=int, default=8000, help="Server port")
parser.add_argument("--db-url", type=str, required=True, help="Postgres URL")
args, _ = parser.parse_known_args()


AsyncSessionLocal = None


@asynccontextmanager
async def lifespan(app: FastAPI):
    global AsyncSessionLocal

    # SQLAlchemy requires the asyncpg driver specified in the URL
    async_db_url = args.db_url.replace(
        "postgresql://", "postgresql+asyncpg://")
    engine = create_async_engine(async_db_url, echo=False)
    AsyncSessionLocal = sessionmaker(
        engine, class_=AsyncSession, expire_on_commit=False)

    yield

    await engine.dispose()

app = FastAPI(lifespan=lifespan)

# Dependency to inject the database session into your routes


async def get_db():
    async with AsyncSessionLocal() as session:
        yield session


@app.get("/ingredient")
async def process_ingredient(name: str, db: AsyncSession = Depends(get_db)):
    if not name:
        raise HTTPException(
            status_code=400, detail="Name parameter is required")

    # Boilerplate example: Execute a simple raw SQL query to test the connection
    try:
        result = await db.execute(text("SELECT 1"))
        db_alive = result.scalar() == 1
    except Exception as e:
        raise HTTPException(
            status_code=500, detail=f"Database connection failed: {str(e)}")

    return {
        "status": "success",
        "name": name,
        "database_connected": db_alive,
        "message": "Ready to implement ORM models or complex queries"
    }

if __name__ == "__main__":
    uvicorn.run("main:app", host="0.0.0.0", port=args.port, reload=True)
