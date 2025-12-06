package handler

// toRecipeResponse transforms the internal Recipe struct to the custom API response format
func toRecipeResponse(url string, recipe *Recipe) *RecipeResponse {
	response := &RecipeResponse{
		URL:          url,
		Instructions: []string{},
		Ingredients:  []string{},
	}

	// Extract recipe details
	response.Recipe = RecipeDetails{
		Name: recipe.Name,
	}

	// Handle optional description
	if recipe.Description != nil {
		response.Recipe.Description = *recipe.Description
	}

	// Extract first image
	if len(recipe.Image) > 0 {
		response.Recipe.Image = recipe.Image[0]
	}

	// Handle optional time fields
	if recipe.PrepTime != nil {
		response.Recipe.PrepTime = *recipe.PrepTime
	}
	if recipe.CookTime != nil {
		response.Recipe.CookTime = *recipe.CookTime
	}
	if recipe.TotalTime != nil {
		response.Recipe.TotalTime = *recipe.TotalTime
	}

	// Extract author name
	if recipe.Author != nil {
		response.Recipe.Author = recipe.Author.Name
	}

	// Copy ingredients
	if len(recipe.RecipeIngredient) > 0 {
		response.Ingredients = recipe.RecipeIngredient
	}

	// Flatten instructions to string array
	response.Instructions = flattenInstructions(recipe.RecipeInstructions)

	return response
}

// flattenInstructions extracts text from RecipeInstruction objects, handling nested HowToSections
func flattenInstructions(instructions []RecipeInstruction) []string {
	var result []string

	for _, inst := range instructions {
		// If this is a HowToSection with nested items, extract from those
		if len(inst.ItemListElement) > 0 {
			result = append(result, flattenInstructions(inst.ItemListElement)...)
		} else if inst.Text != "" {
			result = append(result, inst.Text)
		}
	}

	return result
}

