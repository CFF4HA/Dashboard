import argparse
from contextlib import asynccontextmanager
from fastapi import FastAPI, HTTPException, Depends
from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession
from sqlalchemy.orm import sessionmaker
from sqlalchemy import text
import models
import uuid
import pubchem
import uvicorn

parser = argparse.ArgumentParser(description="Ingredient Server")
parser.add_argument("--port", type=int, default=8000, help="Server port")
parser.add_argument("--db-url", type=str, required=True, help="Postgres URL")
args, _ = parser.parse_known_args()

cache = {}

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

    if name in cache:
        return cache[name]

    names, labels = pubchem.Ingredient(name)
    ingredient = models.Ingredient(Id=(uuid.uuid4()))

    try:
        # First insert, fetching the returned ID
        result = await db.execute(text("INSERT INTO ingredients (id, failed) VALUES (:id, :failed) RETURNING id"), {
            "id": ingredient.Id, "failed": names is None or labels is None})
        if result.scalar() is None:
            raise HTTPException(
                status_code=500, detail="Failed to insert ingredient into database")

        if names is None or labels is None:
            # insert the name so that we can cache the miss on purpose
            result = await db.execute(text("INSERT INTO names (text, ingredient_id) VALUES (:text, :ingredient_id) RETURNING text"), {
                "text": name,
                "ingredient_id": ingredient.Id})
            if result.scalar() is None:
                raise HTTPException(
                    status_code=500, detail="Failed to insert name into database")
        else:
            # Cache the original name as well, since it's a valid synonym
            if name not in names:
                names.append(name)

            for name in names:
                # Added RETURNING id to verify success
                result = await db.execute(text("INSERT INTO names (text, ingredient_id) VALUES (:text, :ingredient_id) RETURNING text"), {
                    "id": uuid.uuid4(), "text": name, "ingredient_id": ingredient.Id})
                if result.scalar() is None:
                    raise HTTPException(
                        status_code=500, detail="Failed to insert name into database")

            for label in labels:
                label_id = uuid.uuid4()
                result = await db.execute(text("INSERT INTO labels (id, type, payload, origin) VALUES (:id, :type, :payload, :origin) RETURNING id"), {
                    "id": label_id, "type": label.Type, "payload": label.Payload, "origin": label.Origin})
                if result.scalar() is None:
                    raise HTTPException(
                        status_code=500, detail="Failed to insert label into database")

                # Returning ingredient_id (or label_id) since join tables often don't have a singular primary key
                result = await db.execute(text("INSERT INTO ingredient_labels (ingredient_id, label_id) VALUES (:ingredient_id, :label_id) RETURNING ingredient_id"), {
                    "ingredient_id": ingredient.Id, "label_id": label_id})
                if result.scalar() is None:
                    raise HTTPException(
                        status_code=500, detail="Failed to associate label with ingredient in database")

        await db.commit()

    except Exception as e:
        await db.rollback()
        raise e

    cache[name] = ingredient


if __name__ == "__main__":
    uvicorn.run("server:app", host="0.0.0.0", port=args.port, reload=True)
