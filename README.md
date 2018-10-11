This is a RESTful Pokédex API written in Go and returns responses in plain text format. Only accepts GET requests.

## ----- Pokédex API Help -----

+ Display this help message: / or /help

## --- Listing Pokemons, Types and Moves ---

+ List all pokemons, moves and types: /list
+ List all pokemons: /list?pokemons
+ List all types: /list?types
+ List all moves: /list?moves
+ List all pokemons for a given type: /list?type={type}
+ Get valid attributes to sort Pokemons by: /list?type={type}&sortby
+ List all pokemons for a given type and sort them by an attribute: /list?type={type}&sortby={attribute}
+ List all pokemons for a given type and sort them by an attribute in reversed order: /list?type={type}&sortby={attribute}&reversed

## --- Getting information about a Pokemon, Type or Move ---

+ /{resourceName}
