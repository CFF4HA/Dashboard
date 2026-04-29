<div style="align-text: center; width: 100%; margin: auto;">
  <img src="static/github-ff-banner.png"/>
</div>

> A web-app designed to facilitate ingredient and product analysis for people with fragrance sensitivities.

## Contents

- [Auxiliary Content](#auxiliary-content)
  - [Installation Guide](./INSTALL.md)
- [Release Notes (v1.0.0)](#release-notes-v100)
  - [Features](#features)
  - [Known Issues (Bugs)](#known-issues-bugs)
  - [Outstanding Features](#outstanding-features)
- [Release Notes (v1.1.0)](#release-notes-v110)

## Auxiliary Content

The following list of content is content that we've specified in a separate document for conciseness. 

- [Installation Guide](./INSTALL.md)

## Release Notes (v1.0.0)

The v1.0.0 release of the Fragrance Free web-app correlates to the version of the app shown in the Capstone Expo. It is mostly feature complete, but lacks in a few areas which were still in-progress at that point in time but were completed for the third checkpoint (May 7th, 2025). 

### Features

The following is a list of features the web-app contained by this version. 

#### Ingredient Retrieval from PubChem

* Ingredients are pulled from PubChem and populated into a database specified at runtime from the command line (when the app is started), more on this in the [Installation Guide]().
  * Allows for ownership of the data ingested from the PubChem API.
* Requests made through the PubChem API respect PubChem's [Reponsible Usage](https://pubchem.ncbi.nlm.nih.gov/docs/pug-rest-tutorial) policy.
  * A built-in limiter ensures that requests to PubChem do not exceed 5 requests per second. 
* Specific Ingredient data from PubChem is transformed and standardized for analysis. The currently used data points from PubChem include:
  * Toxicity > Toxicological Information > Signs & Symptoms
  * Toxicity > Toxicological Information > Adverse Effects
  * Safety and Hazards > Hazards Identification > GHS Classification > GHS Hazard Statements
  * Safety and Hazards > Regulatory Information > California Safe Cosmetics Program (CSCP) Reportable Ingredient
  * Names and Identifiers > Synonyms > Depositor-Supplied Synonyms
* The user interface provides functionality for manually re-parsing ingredient information from PubChem.
* The user interface informs the user when an ingredient couldn't be parsed.

#### Product Creation

* Products can be created from ingredient lists, with careful attention to convenience functions that make this process as seamless as possible (i.e., Copy+Paste Auto-Parsing of Ingredients, etc). 
  * More of this in the [Usage Guide]().
* Ingredients which do not exist in the ingredient table at product creation time are automatically retrieved and populated in the product table.

#### Product Comparison

* Products can be compared to better understand relationships between various products. The following attributes are currently being compared:
  * Shared Ingredients
  * Shared Hazards
  * Shared Effects
  * Shared Symptoms
  * Shared Regulatory Announcements
  * Unique Ingredients (per Product)
  * Unique Hazards (per Product)
  * Unique Effects (per Product)
  * Unique Symptoms (per Product)
  * Unique Regulatory Announcements (per Product)
* Products being compared can be set as "Good" and "Bad" in the product comparison page to facilitate root cause analysis; allowing users to find which ingredients may be aggravating their sensitivities.
  * Ingredients found to be "Bad" for a user can be tagged with a "Bad" tag from the investigation page. 

#### Automated Tagging (Ingredient, Product Analysis)

* Users can create "Tag Rules" which auto-apply tags to ingredients based on the existence (or lack thereof) of some text in an ingredient's label. 
  * Tag Rules can be manually run by the user if they feel the database is out of sync. 
  * Tag Rules, once created, are re-run for all users when new ingredients are created.
  * Tags placed on ingredients bubble up to the product's they compose, allowing users to filter products by the tags applied to their ingredients. 
  * Products created are also automatically tagged with the tags of the ingredients in the product.

#### Basic User Personalization

* Users can favorite products and ingredients for easier retrieval of information.
* Users can edit products post-creation if they later wish to add an "Origin" (i.e., link to where the product can be purchased, etc).
* Users can create their own tags, and their set of automated tagging rules. 

### Known Issues (Bugs)

* The "Automated Tagging" Rule Table exposes a "checkbox" icon to disable certain tag rules from running, this currently does not actually change whether the rule is on/off. 
* Removing a product from the product comparison page's "Compare Products" list (showing the list of products being compared) does not immediately re-render the shared attributes.
* Editing a product's ingredient list and _removing_ an ingredient does not actually persist (meaning the ingredient remains visible).

### Outstanding Features

The following features, as explained earlier, were not yet complete by the time the Expo occured. These features will be included in a future release, `v1.1.0`, on or before May 7th, 2026.

* `Chrome Extension`: An auxiliary component which will allow users to create products without leaving a retailer's site. 
* `Exporting PDFs for Ingredients`: This feature will allow users to export ingredient specific information in a PDF for future refernce in an offline setting.
* `Exporting PDFs for Products`: This feature will allow users to export product specific information in a PDF for future refernce in an offline setting.
* `Ingredient Information Changes Alerting`: A feature that will allow a user to be alerted when an ingredient's information changes (i.e., because the data source changed, or new informationw as added). 
* `Admin, Data Source Configuration`: This feature will allow administrators to configure the frequency that certain data sources are re-ingested, and whether data should be ingested from them all-together.
* `Admin, System Wide Notifications`: This feature will allow administrative users to declare system wide notifications.
* `Admin, Site Usage and Metrics Dashboard`: This feature will expose usage statistics to admin enabled users.
* `Admin, Role Based Access Controls`: This feature will allow an administrative users to assign and remove specific system "roles" to registered users.

## Release Notes (v1.1.0)

> This will be completed and edited as we get closer to the May 7th presentation deadline. Thank you for your patience. 