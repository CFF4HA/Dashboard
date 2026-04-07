# Imports
import pubchem
from collections import Counter, defaultdict
import models

"""
Takes in the names of an ingredient and queries it in Pubchem using our Ingredient function.
Parameters: Name of an ingredient
Returns: List of all Labels associated with the ingredient, and a List comprised of the Types
of Labels
"""
def getLabelTypesList(ing) -> tuple[list, list]:
    _, labels = pubchem.Ingredient(ing)
    types = [label.Type for label in labels]
    return labels, types

"""
Takes in a list of ingredient names and queries them in Pubchem using Ingredient function
Parameters: List of names of ingredients
Returns: A List comprised of each ingredient's List of Labels, and a List of the Types of Labels
present for all ingredients
"""
def massGetLabelTypesList(ings) -> tuple[list[list], list]:
    listLabels = []
    types = None
    for ingredient in ings:
        result = pubchem.Ingredient(ingredient)
        if result is None:
            print(f"Warning: No data found for {ingredient}")
            continue
        _, labels = result
        listLabels.append(labels)
        if (types is None):
            types = [label.Type for label in labels]
        else:
            moreTypes = [label.Type for label in labels]
            for type in moreTypes:
                types.append(type)
    types = sorted(types)
    return listLabels, types

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

# Compares number of defined symptoms for two lists of ingredients.
def massCompareSymptomNumber(ingList1, ingList2, types1, types2):
    numSymptoms1 = types1['symptom']
    numSymptoms2 = types2['symptom']
    print("Comparing Ingredient List 1: "+ingList1+"\n")
    print("with Ingredient List 2: "+ingList2+"\n")
    if (numSymptoms1 > numSymptoms2):
        print("The first ingredient list has more defined symptoms.")
        print(str(numSymptoms1)+" > "+str(numSymptoms2))
    elif (numSymptoms1 < numSymptoms2):
        print("The second ingredient list has more defined symptoms.")
        print(str(numSymptoms2)+" > "+str(numSymptoms1))
    else:
        print("Both ingredient lists have the same number of defined symptoms.")
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

# Compares the number of effect entries for two lists of ingredients.
def massCompareEffectNumber(ingList1, ingList2, types1, types2):
    numEffects1 = types1['effect']
    numEffects2 = types2['effect']
    print("Comparing Ingredient List 1: "+ingList1+"\n")
    print("with Ingredient List 2: "+ingList2+"\n")
    if (numEffects1 > numEffects2):
        print("The first ingredient list has more defined effects.")
        print(str(numEffects1)+" > "+str(numEffects2))
    elif (numEffects1 < numEffects2):
        print("The second ingredient list has more defined effects.")
        print(str(numEffects2)+" > "+str(numEffects1))
    else:
        print("Both ingredient lists have the same number of defined effects.")
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

# Compares the number of defined hazards for two lists of ingredients.
def massCompareHazardNumber(ingList1, ingList2, types1, types2):
    numHazards1 = types1['hazard']
    numHazards2 = types2['hazard']
    print("Comparing Ingredient List 1: "+ingList1+"\n")
    print("with Ingredient List 2: "+ingList2+"\n")
    if (numHazards1 > numHazards2):
        print("The first ingredient list has more defined hazards.")
        print(str(numHazards1)+" > "+str(numHazards2))
    elif (numHazards1 < numHazards2):
        print("The second ingredient list has more defined hazards.")
        print(str(numHazards2)+" > "+str(numHazards1))
    else:
        print("Both ingredient lists have the same number of defined hazards.")
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

# Compares the number of regulation entries for two lists of ingredients.
def massCompareRegulationNumber(ingList1, ingList2, types1, types2):
    numRegulations1 = types1['regulation']
    numRegulations2 = types2['regulation']
    print("Comparing Ingredient List 1: "+ingList1+"\n")
    print("with Ingredient List 2: "+ingList2+"\n")
    if (numRegulations1 > numRegulations2):
        print("The first ingredient list has more defined regulations.")
        print(str(numRegulations1)+" > "+str(numRegulations2))
    elif (numRegulations1 < numRegulations2):
        print("The second ingredient list has more defined regulations.")
        print(str(numRegulations2)+" > "+str(numRegulations1))
    else:
        print("Both ingredient lists have the same number of defined regulations.")
        print(str(numRegulations1)+" = "+str(numRegulations2))

# Takes in a list of Labels and returns a list of all symptoms defined in that list.
def getSymptomList(label1) -> set[str]:
    payloads_by_type = defaultdict(list)
    for label in label1:
        payloads_by_type[label.Type].append(label.Payload)
    symptoms = payloads_by_type['symptom']
    return set(symptoms)

# Takes in a list of Ingredient objects and returns Set of their combined symptoms
def getMassSymptomList(ings : list[models.Ingredient]) -> set[str]:
    symptomSet = set()
    for ing in ings:
        for label in ing.Labels:
            if label.Type == 'symptom':
                symptomSet.add(label.Payload)
    return symptomSet

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

"""
Compares symptoms from two lists of ingredient names. Prints out symptoms that are:
1. Common to both ingredients lists.
2. Unique to the first ingredient list.
3. Unique to the second ingredient list.
"""
def massCompareSymptoms(ingList1, ingList2):
    symptomList1 = []
    symptomList2 = []
    for ing in ingList1:
        symptoms = getSymptomList(getLabelTypesList(ing)[0])
        symptomList1.extend(symptoms)
    for ing in ingList2:
        symptoms = getSymptomList(getLabelTypesList(ing)[0])
        symptomList2.extend(symptoms)

    set1 = set(symptomList1)
    set2 = set(symptomList2)

    commonSymptoms = set1 & set2
    uniqueSymptoms1 = set1 - set2
    uniqueSymptoms2 = set2 - set1

    print("Common Symptoms:", commonSymptoms if commonSymptoms else "None")
    print("Unique to Ingredient List 1: ", uniqueSymptoms1 if uniqueSymptoms1 else "None")
    print("Unique to Ingredient List 2: ", uniqueSymptoms2 if uniqueSymptoms2 else "None")

# Takes in a list of Labels and returns a list of all hazards defined in that list.
def getHazardList(labelList) -> list:
    payloads_by_type = defaultdict(list)
    for label in labelList:
        payloads_by_type[label.Type].append(label.Payload)
    hazards = payloads_by_type['hazard']
    return hazards

# Takes in a list of Ingredient objects and returns Set of their combined hazards
def getMassHazardList(ings : list[models.Ingredient]) -> set[str]:
    hazardSet = set()
    for ing in ings:
        for label in ing.Labels:
            if label.Type == 'hazard':
                hazardSet.add(label.Payload)
    return hazardSet

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

def massCompareHazards(ingList1, ingList2):
    hazardList1 = []
    hazardList2 = []

    for ing in ingList1:
        hazards = getHazardList(getLabelTypesList(ing)[0])
        hazardList1.extend(hazards)

    for ing in ingList2:
        hazards = getHazardList(getLabelTypesList(ing)[0])
        hazardList2.extend(hazards)

    set1 = set(hazardList1)
    set2 = set(hazardList2)

    commonHazards = set1 & set2
    uniqueHazards1 = set1 - set2
    uniqueHazards2 = set2 - set1

    print("Common Hazards:", commonHazards if commonHazards else "None")
    print("Unique to Ingredient List 1:", uniqueHazards1 if uniqueHazards1 else "None")
    print("Unique to Ingredient List 2:", uniqueHazards2 if uniqueHazards2 else "None")

# Takes in a list of Labels and returns a list of all effects defined in that list.
def getEffectList(labelList) -> list:
    payloads_by_type = defaultdict(list)
    for label in labelList:
        payloads_by_type[label.Type].append(label.Payload)
    effects = payloads_by_type['effect']
    return effects

# Takes in a list of Ingredient objects and returns Set of their combined effects
def getMassEffectList(ings : list[models.Ingredient]) -> set[str]:
    effectSet = set()
    for ing in ings:
        for label in ing.Labels:
            if label.Type == 'effect':
                effectSet.add(label.Payload)
    return effectSet
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

"""
Comopares the combined effect lists from two lists of ingredient names.
"""
def massCompareEffects(ingList1, ingList2):
    effectList1 = []
    effectList2 = []

    for ing in ingList1:
        effects = getEffectList(getLabelTypesList(ing)[0])
        effectList1.extend(effects)

    for ing in ingList2:
        effects = getEffectList(getLabelTypesList(ing)[0])
        effectList2.extend(effects)

    set1 = set(effectList1)
    set2 = set(effectList2)

    commonEffects = set1 & set2
    uniqueEffects1 = set1 - set2
    uniqueEffects2 = set2 - set1

    print("Common Effects:", commonEffects if commonEffects else "None")
    print("Unique to Ingredient List 1: ", uniqueEffects1 if uniqueEffects1 else "None")
    print("Unique to Ingredient List 2: ", uniqueEffects2 if uniqueEffects2 else "None")

# Takes in a list of Labels and returns a list of all regulations defined in that list.
def getRegulationList(labelList) -> list:
    payloads_by_type = defaultdict(list)
    for label in labelList:
        payloads_by_type[label.Type].append(label.Payload)
    regulations = payloads_by_type['regulation']
    return regulations

# Takes in a list of Ingredient objects and returns a list of their combined regulations
def getMassRegulationList(ings : list[models.Ingredient]) -> set[str]:
    regSet = set()
    for ing in ings:
        for label in ing.Labels:
            if label.Type == 'regulation':
                regSet.add(label.Payload)
    return regSet

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

def massCompareRegulations(ingList1, ingList2):
    regulationList1 = []
    regulationList2 = []

    for ing in ingList1:
        regulations = getRegulationList(getLabelTypesList(ing)[0])
        regulationList1.extend(regulations)

    for ing in ingList2:
        regulations = getRegulationList(getLabelTypesList(ing)[0])
        regulationList2.extend(regulations)

    set1 = set(regulationList1)
    set2 = set(regulationList2)

    commonRegulations = set1 & set2
    uniqueRegulations1 = set1 - set2
    uniqueRegulations2 = set2 - set1

    print("Common Regulations:", commonRegulations if commonRegulations else "None")
    print("Unique to Ingredient List 1:", uniqueRegulations1 if uniqueRegulations1 else "None")
    print("Unique to Ingredient List 2:", uniqueRegulations2 if uniqueRegulations2 else "None")

