
.PHONY: post.ingredient

URL="localhost:8888"

post.ingredient:
	@curl -sX POST $(URL)/ingredients --data @post_ingredient.json

post.ingredient.error:
	@curl -sX POST $(URL)/ingredients --data @post_ingredient_error.json

delete.ingredient:
	@curl -s $(URL)/ingredients | \
		jq '.ingredients[] | select(.name=="ingredient_10") | .id' | \
		xargs -I@ curl -sX DELETE $(URL)/ingredients/@
