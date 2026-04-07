from dataclasses import dataclass

import re

ORIGIN_PUBCHEM = "PubChem API"


@dataclass
class LabelMatcher:
    Regex: str
    Payload: str


hazards = [
    LabelMatcher(Regex=r"skin irritation", Payload="Skin Irritation"),
    LabelMatcher(Regex=r"carcinogenic", Payload="Carcinogenic"),
    LabelMatcher(Regex=r"skin corrosion", Payload="Skin Corrosion"),
    LabelMatcher(Regex=r"allergic|allergen|allergy", Payload="Allergen"),
]

symptoms = [
    LabelMatcher(Regex=r"cough", Payload="Cough"),
    LabelMatcher(Regex=r"reddness", Payload="Redness"),
]

effects = [
    LabelMatcher(Regex=r"hepatotoxin", Payload="Hepatotoxin Exposure"),
    LabelMatcher(Regex=r"sensitizer|allergic|allergen",
                 Payload="Allergic Reaction"),
]

regulatory_statuses = [
    LabelMatcher(Regex=r"allergen", Payload="Fragrance Allergen"),
]


@ dataclass
class Label:
    Type: str
    Payload: str
    Origin: str


@ dataclass
class Name:
    Text: str
    IngredientId: str


@ dataclass
class Ingredient:
    Id: str
    Labels: list


def Hazard(context: str):
    # this is where we will use regex to detect the type
    # of hazard that we're dealing with.
    found = False
    for h in hazards:
        if re.search(h.Regex, context, re.IGNORECASE):
            found = True
            return Label(Type="hazard", Payload=h.Payload, Origin=ORIGIN_PUBCHEM)

    if not found:
        return Label(Type="hazard", Payload=context, Origin=ORIGIN_PUBCHEM)


def Symptom(context: str):
    # this is where we will use regex to detect the type
    # of symptom that we're dealing with.
    found = False
    for s in symptoms:
        if re.search(s.Regex, context, re.IGNORECASE):
            found = True
            return Label(Type="symptom", Payload=s.Payload, Origin=ORIGIN_PUBCHEM)

    if not found:
        return Label(Type="symptom", Payload=context, Origin=ORIGIN_PUBCHEM)


def Effect(context: str):
    # this is where we will use regex to detect the type
    # of effect that we're dealing with.
    found = False
    for e in effects:
        if re.search(e.Regex, context, re.IGNORECASE):
            found = True
            return Label(Type="effect", Payload=e.Payload, Origin=ORIGIN_PUBCHEM)

    if not found:
        return Label(Type="effect", Payload=context, Origin=ORIGIN_PUBCHEM)


def RegulatoryStatus(context: str):
    # this is where we will use regex to detect the type
    # of regulatory status that we're dealing with.
    found = False
    for r in regulatory_statuses:
        if re.search(r.Regex, context, re.IGNORECASE):
            found = True
            return Label(Type="regulation", Payload=r.Payload, Origin=ORIGIN_PUBCHEM)

    if not found:
        return Label(Type="regulation", Payload=context, Origin=ORIGIN_PUBCHEM)
