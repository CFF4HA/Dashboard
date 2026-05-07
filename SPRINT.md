The following document shows the exact back-end procedures 
required for each sprint story. 

## Sprint 1 

> User Story No. 1: As a standard platform user, I want to personally tailor my experience by annotating ingredients with hazard ratings and notes

### Relevant Endpoints 

- `POST /ingredient/note` creates a note.
- `GET /ingredient/note/get` gets the note for an ingredient.
- `DELETE /ingredient/note/remove` deletes a note from an ingredient.
- `POST /ingredient/tag` tags an ingredient manually.
- `DELETE /ingredient/tag/remove` removes a tag from an ingredient.

> User Story No. 2: As a standard platform user, I want to receive a notification from the system when a given ingredient’s origin data has been modified

- `POST /ingredient/monitor` adds this ingredient to a user's monitoring list.
- `GET /ingredient/monitor/get` retrieves the user's monitoring list.
- `DELETE /ingredient/monitor/remove` removes the ingredient monitor from the relevant table.

> User Story No. 3: As an administrative user, I want to manually configure web scraping destinations and targets through an administrative panel to account for changing website layouts

- `PUT /ingredient/scrape/config/create` adds a scrape config to the system.
- `GET /ingredient/scrape/config/get` retrieves a list of scrape configurations.
- `DELETE /ingredient/scrape/config/remove` removes a scrape configuration from the system.

> User Story No. 4: As an administrative user, I want to declare global system-wide notifications (i.e., status banner alerts) if a notice about a specific chemical must go out

- `PUT /notification/create` creates a notification in the system.
- `DELETE /notification/remove` removes a notification in the system.
- `GET /notification/get` gets a list of all the notifications in the system.
- `GET /notification/get/enabled` gets a list of all the enabled notifications in the system.

## Sprint 2

> User Story No. 1: As a standard platform user, I want to compare multiple ingredient lists (i.e., Products) and understand the key differences in an aesthetically pleasing manner

- `POST /product/compare` compares the list of provided products and returns various attributes (i.e., shared and unique attributes, etc).

> User Story No. 2: As a standard platform user, when comparing ingredient lists, I want to be able to drill down and view ingredient specific information without leaving the product comparison page

- `GET /ingredient/retrieve` pulls the ingredient from the available databases (PubChem, etc).
- `GET /ingredient/get` pulls all ingredients from OUR database.
- `GET /ingredient/get/name` pulls the ingredients with a given name from OUR database.
- `GET /ingredient/get/id` pulls the ingredients with a given ID from OUR database.
- `DELETE /ingredient/remove` removes ingredients from OUR database.

> User Story No. 3: As a standard platform user, I want to be able to save ingredient lists as products so I can reference them in the future

- `PUT /product/create` creates a product in the database.
- `GET /product/scrape` scrapes a product using the LLM and chromium engine.

> User Story No. 4: As a standard platform user, I want to be able to use other people’s saved products in a privacy preserving manner

- `GET /product/get` gets a list of all products in the database available to the user.
- `GET /product/get/name` gets a list of all products in the database by name.
- `GET /product/get/user` gets a list of all products in the database by user.
- `GET /product/get/id` gets a list of all products in the database by id.
- `GET /product/get/ingredients` gets a list of all products in the database by ingredients.
- `GET /product/get/tag` gets a list of products in the database by tag.
- `GET /product/get/operation` gets a list of products using a boolean operation.

## Sprint 3

> User Story No. 4: As an administrative user, I would like to see web and database usage metrics and audit logs from a password protected administrative page

- `GET /usage/metrics/latest` retrieves the latest usage metric for the system.
