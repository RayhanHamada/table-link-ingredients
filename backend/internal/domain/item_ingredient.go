package domain

// ItemIngredient is the domain entity for the join table tm_item_ingredient.
type ItemIngredient struct {
	UUIDItem       string `json:"uuid_item"`
	UUIDIngredient string `json:"uuid_ingredient"`
}
