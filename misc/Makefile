
.PHONY: post.ingredient

URL="localhost:8888"

post.ingredient:
	@curl -X POST $(URL)/ingredients --data @post_ingredient.json

post.ingredient.error:
	@curl -X POST $(URL)/ingredients --data @post_ingredient_error.json
