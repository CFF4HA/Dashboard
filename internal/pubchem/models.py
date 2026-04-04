from dataclasses import dataclass


@dataclass
class Label:
    Type: str
    Payload: str
    Original: str


@dataclass
class Name:
    Text: str
    IngredientId: str


@dataclass
class Ingredient:
    Id: str
