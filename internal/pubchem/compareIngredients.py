# Imports
import pubchem
from collections import Counter, defaultdict

# Get ingredient names; will be used through the file.
ing1 = input("What is the first ingredient?\n")
ing2 = input("What is the second ingredient?\n")

"""
Takes in the names of two ingredients and queries them in Pubchem using our Ingredient function.
Parameters: Names of the two ingredients
Returns: List of all Labels associated with the ingredient, and a count of how many times
each label type appears.
"""
def getLabelTypesList(ing1, ing2) -> tuple[list, list, dict, dict]:
    _, firstLabels = pubchem.Ingredient(ing1)
    _, secondLabels = pubchem.Ingredient(ing2)
    types1 = [label.Type for label in firstLabels]
    types1 = dict(Counter(label.Type for label in firstLabels))
    types2 = dict(Counter(label.Type for label in secondLabels))
    return firstLabels, secondLabels, types1, types2

# Compares the number of unique symptom entries for each ingredient.
def compareSymptomNumber(ing1, ing2, types1, types2):
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

# Compares the number of unique effect entries of each ingredient.
def compareEffectNumber(ing1, ing2, types1, types2):
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

# Compares the nubmer of unique hazard entries of each ingredient.
def compareHazardNumber(types1, types2):
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

# Compares the number of unique regulation entries of each ingredient.
def compareRegulationNumber(types1, types2):
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

# Takes in a list of Labels and returns a list of all symptoms defined in that list.
def getSymptomList(label1) -> list:
    payloads_by_type = defaultdict(list)
    for label in label1:
        payloads_by_type[label.Type].append(label.Payload)
    symptoms = payloads_by_type['symptom']
    return symptoms

"""
Compares two symptom lists from getSymptomList and prints out symptoms that are: 
1. Common to both ingredients listed.
2. Unique to the first ingredient listed.
3. Unique to the second ingredient listed.
Parameters: Two ingredient names and their respective lists of Symptoms from getSymptomList().
"""
def compareSymptoms(ing1, ing2, list1, list2):
    set1 = set(list1)
    set2 = set(list2)

    commonSymptoms = set1 & set2
    uniqueSymptoms1 = set1 - set2
    uniqueSymptoms2 = set2 - set1

    print("Common Symptoms:", commonSymptoms if commonSymptoms else "None")
    print("Unique to "+ing1+": ", uniqueSymptoms1 if uniqueSymptoms1 else "None")
    print("Unique to "+ing2+": ", uniqueSymptoms2 if uniqueSymptoms2 else "None")

# Takes in a list of Labels and returns a list of all hazards defined in that list.
def getHazardList(labelList) -> list:
    payloads_by_type = defaultdict(list)
    for label in labelList:
        payloads_by_type[label.Type].append(label.Payload)
    hazards = payloads_by_type['hazard']
    return hazards

"""
Compares two hazard lists from getHazardList and prints out hazards that are: 
1. Common to both ingredients listed.
2. Unique to the first ingredient listed.
3. Unique to the second ingredient listed.
Parameters: Two ingredient names and their respective lists of Hazards from getHazardList().
"""
def compareHazards(ing1, ing2, list1, list2):
    set1 = set(list1)
    set2 = set(list2)

    commonHazards = set1 & set2
    uniqueHazards1 = set1 - set2
    uniqueHazards2 = set2 - set1

    print("Common Hazards:", commonHazards if commonHazards else "None")
    print("Unique to "+ing1+": ", uniqueHazards1 if uniqueHazards1 else "None")
    print("Unique to "+ing2+": ", uniqueHazards2 if uniqueHazards2 else "None")

# Takes in a list of Labels and returns a list of all effects defined in that list.
def getEffectList(labelList) -> list:
    payloads_by_type = defaultdict(list)
    for label in labelList:
        payloads_by_type[label.Type].append(label.Payload)
    effects = payloads_by_type['effect']
    return effects

"""
Compares two effect lists from getEffectList and prints out effects that are: 
1. Common to both ingredients listed.
2. Unique to the first ingredient listed.
3. Unique to the second ingredient listed.
Parameters: Two ingredient names and their respective lists of Effects from getEffectList().
"""
def compareEffects(ing1, ing2, list1, list2):
    set1 = set(list1)
    set2 = set(list2)

    commonEffects = set1 & set2
    uniqueEffects1 = set1 - set2
    uniqueEffects2 = set2 - set1

    print("Common Effects:", commonEffects if commonEffects else "None")
    print("Unique to "+ing1+": ", uniqueEffects1 if uniqueEffects1 else "None")
    print("Unique to "+ing2+": ", uniqueEffects2 if uniqueEffects2 else "None")

# Takes in a list of Labels and returns a list of all regulations defined in that list.
def getRegulationList(labelList) -> list:
    payloads_by_type = defaultdict(list)
    for label in labelList:
        payloads_by_type[label.Type].append(label.Payload)
    regulations = payloads_by_type['regulation']
    return regulations

"""
Compares two regulation lists from getRegulationList and prints out regulations that are: 
1. Common to both ingredients listed.
2. Unique to the first ingredient listed.
3. Unique to the second ingredient listed.
Parameters: Two ingredient names and their respective lists of Regulations from getRegulationList().
"""
def compareRegulations(ing1, ing2, list1, list2):
    set1 = set(list1)
    set2 = set(list2)

    commonRegulations = set1 & set2
    uniqueRegulations1 = set1 - set2
    uniqueRegulations2 = set2 - set1

    print("Common Regulations:", commonRegulations if commonRegulations else "None")
    print("Unique to "+ing1+": ", uniqueRegulations1 if uniqueRegulations1 else "None")
    print("Unique to "+ing2+": ", uniqueRegulations2 if uniqueRegulations2 else "None")