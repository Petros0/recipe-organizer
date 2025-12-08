package handler

// parseImage parses image field which can be string or array
func parseImage(imageVal interface{}) []string {
	var images []string

	switch v := imageVal.(type) {
	case string:
		if v != "" {
			images = append(images, v)
		}
	case []interface{}:
		for _, item := range v {
			if str, ok := item.(string); ok && str != "" {
				images = append(images, str)
			} else if obj, ok := item.(map[string]interface{}); ok {
				// Could be ImageObject with url property
				if url := getString(obj, "url"); url != "" {
					images = append(images, url)
				}
			}
		}
	case map[string]interface{}:
		// ImageObject
		if url := getString(v, "url"); url != "" {
			images = append(images, url)
		}
	}

	return images
}

// parseInstructions parses recipeInstructions which can be HowToStep, HowToSection, or array
func parseInstructions(instructionsVal interface{}) []RecipeInstruction {
	var instructions []RecipeInstruction

	if instructionsVal == nil {
		return instructions
	}

	switch v := instructionsVal.(type) {
	case []interface{}:
		for _, item := range v {
			if inst := parseInstructionItem(item); inst != nil {
				instructions = append(instructions, *inst)
			}
		}
	default:
		if inst := parseInstructionItem(v); inst != nil {
			instructions = append(instructions, *inst)
		}
	}

	return instructions
}

// parseInstructionItem parses a single instruction item
func parseInstructionItem(item interface{}) *RecipeInstruction {
	obj, ok := item.(map[string]interface{})
	if !ok {
		return nil
	}

	inst := &RecipeInstruction{}
	inst.Type = getString(obj, "@type")
	inst.Text = getString(obj, "text")
	inst.Name = getString(obj, "name")
	inst.URL = getString(obj, "url")

	// Handle HowToSection with itemListElement
	if itemList, ok := obj["itemListElement"].([]interface{}); ok {
		for _, item := range itemList {
			if subInst := parseInstructionItem(item); subInst != nil {
				inst.ItemListElement = append(inst.ItemListElement, *subInst)
			}
		}
	}

	return inst
}

// parseAuthor parses author field
func parseAuthor(authorVal interface{}) *Person {
	if authorVal == nil {
		return nil
	}

	author := &Person{}

	switch v := authorVal.(type) {
	case string:
		author.Name = v
	case map[string]interface{}:
		author.Type = getString(v, "@type")
		author.Name = getString(v, "name")
		author.URL = getString(v, "url")
	case []interface{}:
		// Array of authors - take first
		if len(v) > 0 {
			return parseAuthor(v[0])
		}
	}

	if author.Name == "" {
		return nil
	}

	return author
}

// parseStringOrArray parses a field that can be either a string or array of strings
// This is useful for fields like recipeCategory and recipeCuisine which can be either format
func parseStringOrArray(val interface{}) []string {
	if val == nil {
		return nil
	}

	var result []string

	switch v := val.(type) {
	case string:
		if v != "" {
			result = append(result, v)
		}
	case []interface{}:
		for _, item := range v {
			if str, ok := item.(string); ok && str != "" {
				result = append(result, str)
			}
		}
	case []string:
		for _, str := range v {
			if str != "" {
				result = append(result, str)
			}
		}
	}

	return result
}

// parseNutrition parses nutrition field
func parseNutrition(nutritionVal interface{}) *Nutrition {
	if nutritionVal == nil {
		return nil
	}

	obj, ok := nutritionVal.(map[string]interface{})
	if !ok {
		return nil
	}

	nutrition := &Nutrition{}
	nutrition.Type = getString(obj, "@type")
	nutrition.Calories = getStringPtr(obj, "calories")
	nutrition.FatContent = getStringPtr(obj, "fatContent")
	nutrition.SaturatedFatContent = getStringPtr(obj, "saturatedFatContent")
	nutrition.CholesterolContent = getStringPtr(obj, "cholesterolContent")
	nutrition.SodiumContent = getStringPtr(obj, "sodiumContent")
	nutrition.CarbohydrateContent = getStringPtr(obj, "carbohydrateContent")
	nutrition.FiberContent = getStringPtr(obj, "fiberContent")
	nutrition.SugarContent = getStringPtr(obj, "sugarContent")
	nutrition.ProteinContent = getStringPtr(obj, "proteinContent")

	return nutrition
}
