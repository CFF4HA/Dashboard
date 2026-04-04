# import types
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


def Ingredient(name):
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

                print(f"Synonyms: {synonyms}")

    toxicity = findSectionByHeading(sections, "TOCHeading", "Toxicity")
    if toxicity:
        toxicology_info = findSectionByHeading(
            toxicity.get("Section", []), "TOCHeading", "Toxicological Information")

        if toxicology_info:
            signs_symptoms = findSectionByHeading(
                toxicology_info.get("Section", []), "TOCHeading", "Signs and Symptoms")
            adverse_effects = findSectionByHeading(
                toxicology_info.get("Section", []), "TOCHeading", "Adverse Effects")
            ongoing_test_stat = findSectionByHeading(
                toxicology_info.get("Section", []), "TOCHeading", "Ongoing Test Status")

            print(f"signs and symptoms: {signs_symptoms}")
            print(f"adverse effects: {adverse_effects}")
            print(f"ongoing test status: {ongoing_test_stat}")

    safety = findSectionByHeading(sections, "TOCHeading", "Safety and Hazards")
    if safety:
        regulatory_info = findSectionByHeading(
            safety.get("Section", []), "TOCHeading", "Regulatory Information")

        if regulatory_info:
            cscp = findSectionByHeading(regulatory_info.get(
                "Information"), "Name", "California Safe Cosmetics Program (CSCP) Reportable Ingredient")

            print(f"CSCP: {cscp}")

    return None
