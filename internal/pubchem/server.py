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
import compareIngredients as comp
import json
import string

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

    ingredient = models.Ingredient(Id=(uuid.uuid4()), Labels=labels)

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

@app.get("/ingredient/compare")
async def compare_products(product1: str, product2: str, db: AsyncSession = Depends(get_db)):
    product1 = product1.split(",")
    product2 = product2.split(",")
    ingredients1 : list[models.Ingredient] = []
    ingredients2 : list[models.Ingredient] = []
    for ing in product1:
        if ing in cache:
            ingredients1.append(cache[ing])
        else:
            _, labels = pubchem.Ingredient(ing)
            ingredient = models.Ingredient(Id=(uuid.uuid4()), Labels=labels)
            ingredients1.append(ingredient)
            cache[ing] = ingredient
    for ing in product2:
        if ing in cache:
            ingredients2.append(cache[ing])
        else:
            _, labels = pubchem.Ingredent(ing)
            ingredient = models.Ingredient(Id=(uuid.uuid4()), Labels=labels)
            ingredients2.append(ingredient)
            cache[ing] = ingredient
    # Symptom comparisons
    symptoms1 = comp.getMassSymptomList(ingredients1)
    symptoms2 = comp.getMassSymptomList(ingredients2)
    commonSymptoms = symptoms1 & symptoms2
    commonSymptomCount = len(commonSymptoms)
    uniqueSymptoms1 = symptoms1 - symptoms2
    uniqueSymptom1Count = len(uniqueSymptoms1)
    uniqueSymptoms2 = symptoms2 - symptoms1
    uniqueSymptom2Count = len(uniqueSymptoms2)

    # Hazard comparisons
    hazards1 = comp.getMassHazardList(ingredients1)
    hazards2 = comp.getMassHazardList(ingredients2)
    commonHazards = hazards1 & hazards2
    commonHazardCount = len(commonHazards)
    uniqueHazards1 = hazards1 - hazards2
    uniqueHazard1Count = len(uniqueHazards1)
    uniqueHazards2 = hazards2 - hazards1
    uniqueHazard2Count = len(uniqueHazards2)

    # Effect comparisons
    effects1 = comp.getMassEffectList(ingredients1)
    effects2 = comp.getMassEffectList(ingredients2)
    commonEffects = effects1 & effects2
    commonEffectCount = len(commonEffects)
    uniqueEffects1 = effects1 - effects2
    uniqueEffect1Count = len(uniqueEffects1)
    uniqueEffects2 = effects2 - effects1
    uniqueEffect2Count = len(uniqueEffects2)

    # Regulation comparisons
    reg1 = comp.getMassRegulationList(ingredients1)
    reg2 = comp.getMassRegulationList(ingredients2)
    commonRegs = reg1 - reg2
    commonRegCount = len(commonRegs)
    uniqueReg1 = reg1 - reg2
    uniqueReg1Count = len(uniqueReg1)
    uniqueReg2 = reg2 - reg1
    uniqueReg2Count = len(uniqueReg2)

    # Time to put it all into a json format
    allInfo = {
        "symptoms": {
            "common": {
                "data": list(commonSymptoms),
                "count": commonSymptomCount
            },
            "unique1": {
                "data": list(uniqueSymptoms1),
                "count": uniqueSymptom1Count
            },
            "unique2": {
                "data": list(uniqueSymptoms2),
                "count": uniqueSymptom2Count
            }
        },
        "hazards": {
            "common": {
                "data": list(commonHazards),
                "count": commonHazardCount
            },
            "unique1": {
                "data": list(uniqueHazards1),
                "count": uniqueHazard1Count
            },
            "unique2": {
                "data": list(uniqueHazards2),
                "count": uniqueHazard2Count
            }
        },
        "effects": {
            "common": {
                "data": list(commonEffects),
                "count": commonEffectCount
            },
            "unique1": {
                "data": list(uniqueEffects1),
                "count": uniqueEffect1Count
            },
            "unique2": {
                "data": list(uniqueEffects2),
                "count": uniqueEffect2Count
            }
        },
        "regulations": {
            "common": {
                "data": list(commonRegs),
                "count": commonRegCount
            },
            "unique1": {
                "data": list(uniqueReg1),
                "count": uniqueReg1Count
            },
            "unique2": {
                "data": list(uniqueReg2),
                "count": uniqueReg2Count
            }
        }
    }
    return allInfo

if __name__ == "__main__":
    uvicorn.run("server:app", host="0.0.0.0", port=args.port, reload=True)
