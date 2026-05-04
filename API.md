# API 

## Product Related Endpoints

* /product/create/manual (POST)
  
  This endpoint creates a product. 
  _Required Query Parameters_
  
  1. Name: string, the name of the product (required) [name='name'].
  2. Origin: string, a link to the product (optional) [name='origin'].
  3. Ingredient List: string, a comma separated list of ingredients (required) [name='ingredient_list']

* /product/create/scraping (POST)

  This endpoint creates a product from a URL, using the 
  scraping engine.
  _Required Query Parameters_
  
  1. URL: string, the URL to the product (required) [name='url'].

* /product/edit (PUT)

  This endpoint is used to update an existing product. It will overwrite 
  all parameters that it receives.
  _Required Query Parameters_

  1. ID: string, the id of the product that is being edited (required) [name='id'].
  2. Name: string, the name of the product (required) [name='name'].
  3. Origin: string, a link to the product (required) [name='origin'].
  4. Ingredient List: string, a comma separated list of ingredients (required) [name='ingredient_list'].

* /product/remove (DELETE)

  This endpoint is used to delete a product from the database.
  _Required Query Parameters_
  
  1. ID: string, the id of the product that is being deleted (required) [name='id'].

* /product/get (GET)

  This endpoint is used to get a list of products with no filter from the database.
  _Required Query Parameters_
  
  1. Cursor: string, a server-side provided "pointer" or "marker" to only load a specific number of results at a time.

* /product/get/id (GET)

  This endpoint is used to get a single product by its ID.
  _Required Query Parameters_

  1. ID: string, the id of the product to retrieve (required) [name='id'].

* /product/get/name (GET)

  This endpoint is used to get a list of products by name.
  _Required Query Parameters_

  1. Name: string, the name (or partial name) to search for (required) [name='name'].
  2. Cursor: string, pagination cursor (optional) [name='cursor'].

* /product/get/tags (GET)

  This endpoint is used to get a list of products by tags.
  _Required Query Parameters_

  1. Tags: string, a comma separated list of tag IDs to filter by (required) [name='tags'].
  2. Cursor: string, pagination cursor (optional) [name='cursor'].

* /product/get/ingredients (GET)

  This endpoint is used to get a list of products by ingredients.
  _Required Query Parameters_

  1. Ingredients: string, a comma separated list of ingredient IDs to filter by (required) [name='ingredients'].
  2. Cursor: string, pagination cursor (optional) [name='cursor'].

## Tag Related Endpoints 

* /tag/rule/create (POST)

  This endpoint is used to create a tagging rule.
  _Required Query Parameters_:

  1. Name: string, the name of the tagging rule - in case a user wishes to name them (optional).
  2. Description: string, the description given to the tagging rule (optional)
  3. Tag: string, the name of the tags to apply if the conditions are met.
  4. Conditions: string, a boolean algebra compatible statement containing strings to match against.

* /tag/rule/edit (PUT)

  This endpoint is used to edit an existing tagging rule. It will overwrite all parameters that it receives.
  _Required Query Parameters_

  1. ID: string, the id of the tagging rule to edit (required) [name='id'].
  2. Name: string, the updated name of the tagging rule (required) [name='name'].
  3. Description: string, the updated description of the tagging rule (required) [name='description'].
  4. Tag: string, the updated tag to apply if the conditions are met (required) [name='tag'].
  5. Conditions: string, the updated boolean algebra compatible statement to match against (required) [name='conditions'].

* /tag/rule/remove (DELETE)

  This endpoint is used to delete a tagging rule from the database.
  _Required Query Parameters_

  1. ID: string, the id of the tagging rule to delete (required) [name='id'].

* /tag/rule/get (GET)

  This endpoint is used to get a list of all tagging rules.
  _Required Query Parameters_

  1. Cursor: string, pagination cursor (optional) [name='cursor'].

* /tag/rule/get/tags (GET)

  This endpoint is used to get all tagging rules that apply a given tag.
  _Required Query Parameters_

  1. Tag: string, the name of the tag to filter rules by (required) [name='tag'].
  2. Cursor: string, pagination cursor (optional) [name='cursor'].

* /tag/rule/get/name (GET)

  This endpoint is used to get tagging rules by name.
  _Required Query Parameters_

  1. Name: string, the name (or partial name) to search for (required) [name='name'].
  2. Cursor: string, pagination cursor (optional) [name='cursor'].

* /tag/create (POST)

  This endpoint is used to create a tag.
  _Required Query Parameters_

  1. Name: string, the name of the tag (required) [name='name'].
  2. Description: string, a description of the tag (optional) [name='description'].

* /tag/edit (PUT)

  This endpoint is used to edit an existing tag.
  _Required Query Parameters_

  1. ID: string, the id of the tag to edit (required) [name='id'].
  2. Name: string, the updated name of the tag (optional) [name='name'].
  3. Description: string, the updated description of the tag (optional) [name='description'].

* /tag/remove (DELETE)

  This endpoint is used to delete a tag from the database.
  _Required Query Parameters_

  1. ID: string, the id of the tag to delete (required) [name='id'].

* /tag/get (GET)

  This endpoint will return a list of tags.
  _Required Query Parameters_

  1. Cursor: string, pagination cursor (optional) [name='cursor'].

* /tag/get/name (GET)

  This endpoint is used to get a list of tags by name.
  _Required Query Parameters_

  1. Name: string, the name (or partial name) to search for (required) [name='name'].
  2. Cursor: string, pagination cursor (optional) [name='cursor'].

* /tag/get/id (GET)

  This endpoint is used to get a single tag by its ID.
  _Required Query Parameters_

  1. ID: string, the id of the tag to retrieve (required) [name='id'].

* /tag/get/ingredient (GET)

  This endpoint will be used to get the tags from an ingredient.
  _Required Query Parameters_

  1. ID: string, the id of the ingredient for whose tags will be retrieved (required) [name='id'].
  2. Cursor: string, pagination cursor (optional) [name='cursor'].

* /tag/get/product (GET)

  This endpoint will be used to get the tags from a product.
  _Required Query Parameters_

  1. ID: string, the id of the product whose tags will be retrieved (required) [name='id'].
  2. Cursor: string, pagination cursor (optional) [name='cursor'].

## Ingredient Related Endpoints 

* /ingredient/pull (POST)

  This endpoint is equivalent to the ingredient creation endpoint. This will work 
  by running the relevant backend functions to "sync" the ingredient, requires a "primary_name".
  _Required Query Parameters_

  1. Name: string, the primary name of the ingredient to sync (required) [name='primary_name'].

* /ingredient/edit (PUT)

  This endpoint is used to edit the internal representation of an ingredient.
  _Required Query Parameters_

  1. Name: string, the primary_name of the ingredient.
  2. Tags: []string, a comma separated list of tag IDs.
  3. Ingredients: []string, a comma separated list of ingredient IDs.
  4. Synonyms: []string, a comma separated list of names by which the ingredient is known.
  5. Labels: []string, a comma separated list of labels the ingredient contains.

* /ingredient/remove (DELETE)

  This endpoint is used to remove an ingredient from the database.
  _Required Query Parameters_

  1. id: string, the id of the ingredient to remove (required) [name='id'].

* /ingredient/tag (POST)

  This endpoint is used to tag an ingredient. Requires a string representation of the tag to use.
  _Required Query Parameters_

  1. ID: string, the id of the ingredient to tag (required) [name='id'].
  2. Tag: string, the name of the tag to apply (required) [name='tag'].

* /ingredient/tag/remove (DELETE)

  This endpoint is used to remove a tag from an ingredient.
  _Required Query Parameters_

  1. id: string, the id of the ingredient (required) [name='id'].
  2. Tag: string, the name of the tag to remove (required) [name='tag'].

* /ingredient/get (GET)

  This endpoint is used to get a list of all ingredients with no filter.
  _Required Query Parameters_

  1. Cursor: string, pagination cursor (optional) [name='cursor'].

* /ingredient/get/id (GET)

  This endpoint is used to get a single ingredient by its ID.
  _Required Query Parameters_

  1. ID: string, the id of the ingredient to retrieve (required) [name='id'].

* /ingredient/get/name (GET)

  This endpoint is used to get a list of ingredients by name.
  _Required Query Parameters_

  1. Name: string, the name (or partial name) to search for (required) [name='name'].
  2. Cursor: string, pagination cursor (optional) [name='cursor'].

* /ingredient/get/labels (GET)

  This endpoint is used to get a list of ingredients by label.
  _Required Query Parameters_

  1. Labels: string, a comma separated list of labels to filter by (required) [name='labels'].
  2. Cursor: string, pagination cursor (optional) [name='cursor'].

* /ingredient/get/tags (GET)

  This endpoint is used to get a list of ingredients by tag.
  _Required Query Parameters_

  1. Tags: string, a comma separated list of tag IDs to filter by (boolean algebra compatible) (required) [name='tags'].
  2. Cursor: string, pagination cursor (optional) [name='cursor'].

* /ingredient/get/products (GET)

  This endpoint returns the ingredients in a product, where the product is an ID.
  _Required Query Parameters_

  1. ID: string, the id of the product whose ingredients will be retrieved (required) [name='id'].
  2. Cursor: string, pagination cursor (optional) [name='cursor'].
