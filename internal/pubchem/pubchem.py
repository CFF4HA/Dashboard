import models
import warnings
import pubchempy as pcp
import requests


def findSectionByHeading(sections, key, heading):
    for section in sections:
        if section.get(key) == heading:
            return section

    return None


def findValueFromObjList(objects, key, name):
    for obj in objects:
        if obj.get(key) == name:
            return obj.get("Value", {}).get("StringWithMarkup", [{}])[0].get("String")

    return None


def jsonFromCompoundName(compound):
    # Grab the compounds associated with that name, then get the CID
    compounds = pcp.get_compounds(compound, "name")
    if not compounds:  # Look if there is actually a compound with that name
        warnings.warn("No chemical compound with that name found on PubChem.")
        return None
    cid = compounds[0].cid

    # Use CID to get json of all data on compound
    url = f"https://pubchem.ncbi.nlm.nih.gov/rest/pug_view/data/compound/{cid}/JSON"
    response = requests.get(url)

    if response.status_code != 200:  # Make sure we get some kind of response
        warnings.warn("Failed to fetch data")
        return None

    return response.json()


def Ingredient(name) -> (list[str], list[models.Label]):
    ingredient_names = []
    ingredient_labels = []

    data = jsonFromCompoundName(name)
    if not data:
        return None
    sections = data.get("Record", {}).get("Section", [])

    # Get the individual sections we care about.
    names = findSectionByHeading(
        sections, "TOCHeading", "Names and Identifiers")
    if names:
        synonyms = findSectionByHeading(
            names.get("Section"), "TOCHeading", "Synonyms")

        if synonyms:
            depositors = findSectionByHeading(synonyms.get(
                "Section", []), "TOCHeading", "Depositor-Supplied Synonyms")

            if depositors:
                synonyms = depositors.get("Information", [])[0].get(
                    "Value", {}).get("StringWithMarkup", [{}])

                for synonym in synonyms:
                    ingredient_names.append(synonym.get("String"))

                print(f"Synonyms: {ingredient_names}")
    else:
        ingredient_names.append(name)

    toxicity = findSectionByHeading(sections, "TOCHeading", "Toxicity")
    if toxicity:
        toxicology_info = findSectionByHeading(
            toxicity.get("Section", []), "TOCHeading", "Toxicological Information")

        if toxicology_info:
            signs_symptoms = findSectionByHeading(
                toxicology_info.get("Section", []), "TOCHeading", "Signs and Symptoms")
            adverse_effects = findSectionByHeading(
                toxicology_info.get("Section", []), "TOCHeading", "Adverse Effects")

            if signs_symptoms:
                for symptom_type in signs_symptoms.get("Information", []):
                    for s in symptom_type.get("Value", {}).get("StringWithMarkup", [{}]):
                        context = s.get("String", "")
                        ingredient_labels.append(models.Symptom(context))

            if adverse_effects:
                for effect_type in adverse_effects.get("Information", []):
                    for e in effect_type.get("Value", {}).get("StringWithMarkup", [{}]):
                        context = e.get("String", "")
                        ingredient_labels.append(models.Effect(context))

    safety = findSectionByHeading(sections, "TOCHeading", "Safety and Hazards")
    if safety:
        hazards = findSectionByHeading(
            safety.get("Section", []), "TOCHeading", "Hazards Identification")
        if hazards:
            ghs = findSectionByHeading(
                hazards.get("Section", []), "TOCHeading", "GHS Classification")

            if ghs:
                classifications = findSectionByHeading(
                    ghs.get("Information", []), "Name", "GHS Hazard Statements").get("Value", {}).get("StringWithMarkup", [])

                # This is where we actually add it to the label list
                for classification in classifications:
                    context = classification.get("String", "")
                    ingredient_labels.append(models.Hazard(context))

        regulatory_info = findSectionByHeading(
            safety.get("Section", []), "TOCHeading", "Regulatory Information")

        if regulatory_info:
            cscp = findSectionByHeading(regulatory_info.get(
                "Information"), "Name", "California Safe Cosmetics Program (CSCP) Reportable Ingredient")

            if cscp:
                for c in cscp.get("Value", {}).get("StringWithMarkup", []):
                    context = c.get("String", "")

                    ingredient_labels.append(models.RegulatoryStatus(context))

    return ingredient_names, ingredient_labels
