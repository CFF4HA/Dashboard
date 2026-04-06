import pubchem
from collections import Counter

ing1 = input("What is the first ingredient?\n")
ing2 = input("What is the second ingredient?\n")

def getLabelTypesList(ing1, ing2) -> tuple[dict, dict]:
    _, firstLabels = pubchem.Ingredient(ing1)
    _, secondLabels = pubchem.Ingredient(ing2)
    types1 = [label.Type for label in firstLabels]
    types1 = dict(Counter(label.Type for label in firstLabels))
    types2 = dict(Counter(label.Type for label in secondLabels))
    return types1, types2

def compareSymptoms(types1, types2):
    numSymptoms1 = types1['symptom']
    numSymptoms2 = types2['symptom']
    if (numSymptoms1 > numSymptoms2):
        print(ing1+" has more defined symptoms than "+ing2)
        print(str(numSymptoms1)+" > "+str(numSymptoms2))
    elif (numSymptoms1 < numSymptoms2):
        print(ing2+" has more defined symptoms than "+ing1)
        print(str(numSymptoms2)+" > "+str(numSymptoms1))
    else:
        print(ing1+" and "+ing2+" have the same number of defined symptoms.")
        print(str(numSymptoms1)+" = "+str(numSymptoms2))

def compareEffects(types1, types2):
    numEffects1 = types1['effect']
    numEffects2 = types2['effect']
    if (numEffects1 > numEffects2):
        print(ing1+" has more defined effects than "+ing2)
        print(str(numEffects1)+" > "+str(numEffects2))
    elif (numEffects1 < numEffects2):
        print(ing2+" has more defined effects than "+ing1)
        print(str(numEffects2)+" > "+str(numEffects1))
    else:
        print(ing1+" and "+ing2+" have the same number of defined effects.")
        print(str(numEffects1)+" = "+str(numEffects2))

def compareHazards(types1, types2):
    numHazards1 = types1['hazard']
    numHazards2 = types2['hazard']
    if (numHazards1 > numHazards2):
        print(ing1+" has more defined hazards than "+ing2)
        print(str(numHazards1)+" > "+str(numHazards2))
    elif (numHazards1 < numHazards2):
        print(ing2+" has more defined hazards than "+ing1)
        print(str(numHazards2)+" > "+str(numHazards1))
    else:
        print(ing1+" and "+ing2+" have the same number of defined hazards.")
        print(str(numHazards1)+" = "+str(numHazards2))

def compareRegulations(types1, types2):
    numRegulations1 = types1['regulation']
    numRegulations2 = types2['regulation']
    if (numRegulations1 > numRegulations2):
        print(ing1+" has more defined regulations than "+ing2)
        print(str(numRegulations1)+" > "+str(numRegulations2))
    elif (numRegulations1 < numRegulations2):
        print(ing2+" has more defined regulations than "+ing1)
        print(str(numRegulations2)+" > "+str(numRegulations1))
    else:
        print(ing1+" and "+ing2+" have the same number of defined regulations.")
        print(str(numRegulations1)+" = "+str(numRegulations2))

types1, types2 = getLabelTypesList(ing1, ing2)
compareSymptoms(types1, types2)
compareEffects(types1, types2)
compareHazards(types1, types2)
compareRegulations(types1, types2)