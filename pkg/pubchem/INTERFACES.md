# Pubchem: Interfaces 

This document describes the expected interfaces the user can use to interact with the PubChem API.

## General Overview

This library will return http.Requests crafted for use by the golang's HTTP Client class.

## Methods 

- `func pubchem.GetCompoundId(name string) int` ([more info]())
- `func pubchem.GetSubstanceId(name string) int`
- `func pubchem.GetCompoundAsJSON(id int) json.RawMessage`
- `func pubchem.GetSubstanceAsJSON(id int) json.RawMessage`

### Detailed Information

#### GetCompoundId 

This function uses the `https://pubchem.ncbi.nlm.nih.gov/rest/pug/compound/name/$NAME/cids/JSON` endpoint, where `$NAME` is the string
representation of the name.

#### GetSubstanceId 

This function uses the `https://pubchem.ncbi.nlm.nih.gov/rest/pug/substance/name/$NAME/sids/JSON` endpoint where `$NAME` is the string 
representation of the name.

#### GetCompoundAsJSON

This function uses the `https://pubchem.ncbi.nlm.nih.gov/rest/pug_view/data/compound/$ID/JSON` endpoint where the `$ID` is the integer representation 
of the id for the relevant compound.

#### GetSubstanceAsJSON 

This function uses the `https://pubchem.ncbi.nlm.nih.gov/rest/pug_view/data/substance/$ID/JSON` endpoint where the `$ID` is the integer representation 
of the id for the relevant compound.
