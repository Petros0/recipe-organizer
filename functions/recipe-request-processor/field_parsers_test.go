package handler

import (
	"testing"
)

func TestParseImage(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected []string
	}{
		{
			name:     "single string",
			input:    "https://example.com/image.jpg",
			expected: []string{"https://example.com/image.jpg"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: nil,
		},
		{
			name:     "array of strings",
			input:    []interface{}{"img1.jpg", "img2.jpg", "img3.jpg"},
			expected: []string{"img1.jpg", "img2.jpg", "img3.jpg"},
		},
		{
			name:     "array with empty strings filtered",
			input:    []interface{}{"img1.jpg", "", "img2.jpg"},
			expected: []string{"img1.jpg", "img2.jpg"},
		},
		{
			name: "array of ImageObjects",
			input: []interface{}{
				map[string]interface{}{"@type": "ImageObject", "url": "img1.jpg"},
				map[string]interface{}{"@type": "ImageObject", "url": "img2.jpg"},
			},
			expected: []string{"img1.jpg", "img2.jpg"},
		},
		{
			name: "mixed array with strings and ImageObjects",
			input: []interface{}{
				"direct.jpg",
				map[string]interface{}{"url": "object.jpg"},
			},
			expected: []string{"direct.jpg", "object.jpg"},
		},
		{
			name: "single ImageObject",
			input: map[string]interface{}{
				"@type": "ImageObject",
				"url":   "https://example.com/single.jpg",
			},
			expected: []string{"https://example.com/single.jpg"},
		},
		{
			name: "ImageObject without url",
			input: map[string]interface{}{
				"@type": "ImageObject",
			},
			expected: nil,
		},
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
		{
			name:     "unsupported type (number)",
			input:    123,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseImage(tt.input)

			if len(result) != len(tt.expected) {
				t.Errorf("len = %d, want %d", len(result), len(tt.expected))
				return
			}

			for i, img := range result {
				if img != tt.expected[i] {
					t.Errorf("image[%d] = %q, want %q", i, img, tt.expected[i])
				}
			}
		})
	}
}

func TestParseInstructions(t *testing.T) {
	tests := []struct {
		name          string
		input         interface{}
		expectedCount int
		checkFirst    string
	}{
		{
			name:          "nil input",
			input:         nil,
			expectedCount: 0,
		},
		{
			name: "array of HowToStep",
			input: []interface{}{
				map[string]interface{}{"@type": "HowToStep", "text": "Step 1"},
				map[string]interface{}{"@type": "HowToStep", "text": "Step 2"},
				map[string]interface{}{"@type": "HowToStep", "text": "Step 3"},
			},
			expectedCount: 3,
			checkFirst:    "Step 1",
		},
		{
			name: "HowToSection with nested steps",
			input: []interface{}{
				map[string]interface{}{
					"@type": "HowToSection",
					"name":  "Preparation",
					"itemListElement": []interface{}{
						map[string]interface{}{"@type": "HowToStep", "text": "Prep step 1"},
						map[string]interface{}{"@type": "HowToStep", "text": "Prep step 2"},
					},
				},
			},
			expectedCount: 1,
		},
		{
			name: "single HowToStep object (not array)",
			input: map[string]interface{}{
				"@type": "HowToStep",
				"text":  "Single step",
			},
			expectedCount: 1,
			checkFirst:    "Single step",
		},
		{
			name:          "empty array",
			input:         []interface{}{},
			expectedCount: 0,
		},
		{
			name: "array with non-object items filtered",
			input: []interface{}{
				"string item",
				map[string]interface{}{"@type": "HowToStep", "text": "Valid step"},
				123,
			},
			expectedCount: 1,
			checkFirst:    "Valid step",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseInstructions(tt.input)

			if len(result) != tt.expectedCount {
				t.Errorf("count = %d, want %d", len(result), tt.expectedCount)
				return
			}

			if tt.checkFirst != "" && len(result) > 0 {
				if result[0].Text != tt.checkFirst {
					t.Errorf("first instruction text = %q, want %q", result[0].Text, tt.checkFirst)
				}
			}
		})
	}
}

func TestParseInstructionItem(t *testing.T) {
	tests := []struct {
		name            string
		input           interface{}
		wantNil         bool
		wantType        string
		wantText        string
		wantNestedCount int
	}{
		{
			name: "HowToStep",
			input: map[string]interface{}{
				"@type": "HowToStep",
				"text":  "Mix ingredients",
				"name":  "Step 1",
			},
			wantType: "HowToStep",
			wantText: "Mix ingredients",
		},
		{
			name: "HowToSection with nested items",
			input: map[string]interface{}{
				"@type": "HowToSection",
				"name":  "Cooking",
				"itemListElement": []interface{}{
					map[string]interface{}{"@type": "HowToStep", "text": "Cook step 1"},
					map[string]interface{}{"@type": "HowToStep", "text": "Cook step 2"},
				},
			},
			wantType:        "HowToSection",
			wantNestedCount: 2,
		},
		{
			name:    "nil input",
			input:   nil,
			wantNil: true,
		},
		{
			name:    "string input",
			input:   "not an object",
			wantNil: true,
		},
		{
			name:    "number input",
			input:   42,
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseInstructionItem(tt.input)

			if tt.wantNil {
				if result != nil {
					t.Errorf("Expected nil, got %+v", result)
				}
				return
			}

			if result == nil {
				t.Fatal("Expected result but got nil")
			}

			if result.Type != tt.wantType {
				t.Errorf("Type = %q, want %q", result.Type, tt.wantType)
			}

			if tt.wantText != "" && result.Text != tt.wantText {
				t.Errorf("Text = %q, want %q", result.Text, tt.wantText)
			}

			if len(result.ItemListElement) != tt.wantNestedCount {
				t.Errorf("ItemListElement count = %d, want %d", len(result.ItemListElement), tt.wantNestedCount)
			}
		})
	}
}

func TestParseAuthor(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		wantNil  bool
		wantName string
		wantType string
		wantURL  string
	}{
		{
			name:     "string author",
			input:    "John Doe",
			wantName: "John Doe",
		},
		{
			name: "Person object",
			input: map[string]interface{}{
				"@type": "Person",
				"name":  "Jane Smith",
				"url":   "https://example.com/jane",
			},
			wantName: "Jane Smith",
			wantType: "Person",
			wantURL:  "https://example.com/jane",
		},
		{
			name: "Organization object",
			input: map[string]interface{}{
				"@type": "Organization",
				"name":  "Test Kitchen",
			},
			wantName: "Test Kitchen",
			wantType: "Organization",
		},
		{
			name: "array of authors - takes first",
			input: []interface{}{
				map[string]interface{}{"@type": "Person", "name": "First Author"},
				map[string]interface{}{"@type": "Person", "name": "Second Author"},
			},
			wantName: "First Author",
			wantType: "Person",
		},
		{
			name:    "nil input",
			input:   nil,
			wantNil: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantNil: true,
		},
		{
			name: "object without name",
			input: map[string]interface{}{
				"@type": "Person",
				"url":   "https://example.com",
			},
			wantNil: true,
		},
		{
			name:    "empty array",
			input:   []interface{}{},
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseAuthor(tt.input)

			if tt.wantNil {
				if result != nil {
					t.Errorf("Expected nil, got %+v", result)
				}
				return
			}

			if result == nil {
				t.Fatal("Expected result but got nil")
			}

			if result.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", result.Name, tt.wantName)
			}

			if tt.wantType != "" && result.Type != tt.wantType {
				t.Errorf("Type = %q, want %q", result.Type, tt.wantType)
			}

			if tt.wantURL != "" && result.URL != tt.wantURL {
				t.Errorf("URL = %q, want %q", result.URL, tt.wantURL)
			}
		})
	}
}

func TestParseNutrition(t *testing.T) {
	tests := []struct {
		name         string
		input        interface{}
		wantNil      bool
		wantCalories string
		wantProtein  string
	}{
		{
			name: "full nutrition info",
			input: map[string]interface{}{
				"@type":               "NutritionInformation",
				"calories":            "250 kcal",
				"proteinContent":      "10g",
				"fatContent":          "15g",
				"carbohydrateContent": "30g",
			},
			wantCalories: "250 kcal",
			wantProtein:  "10g",
		},
		{
			name: "partial nutrition info",
			input: map[string]interface{}{
				"@type":    "NutritionInformation",
				"calories": "100 calories",
			},
			wantCalories: "100 calories",
		},
		{
			name:    "nil input",
			input:   nil,
			wantNil: true,
		},
		{
			name:    "string input",
			input:   "not an object",
			wantNil: true,
		},
		{
			name:    "empty object",
			input:   map[string]interface{}{},
			wantNil: false, // Returns empty Nutrition struct
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseNutrition(tt.input)

			if tt.wantNil {
				if result != nil {
					t.Errorf("Expected nil, got %+v", result)
				}
				return
			}

			if result == nil {
				t.Fatal("Expected result but got nil")
			}

			if tt.wantCalories != "" {
				if result.Calories == nil {
					t.Error("Expected Calories but got nil")
				} else if *result.Calories != tt.wantCalories {
					t.Errorf("Calories = %q, want %q", *result.Calories, tt.wantCalories)
				}
			}

			if tt.wantProtein != "" {
				if result.ProteinContent == nil {
					t.Error("Expected ProteinContent but got nil")
				} else if *result.ProteinContent != tt.wantProtein {
					t.Errorf("ProteinContent = %q, want %q", *result.ProteinContent, tt.wantProtein)
				}
			}
		})
	}
}

func TestParseStringOrArray(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected []string
	}{
		{
			name:     "single string",
			input:    "Main Course",
			expected: []string{"Main Course"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: nil,
		},
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
		{
			name:     "array of strings (interface)",
			input:    []interface{}{"Main Course", "Dinner"},
			expected: []string{"Main Course", "Dinner"},
		},
		{
			name:     "array of strings (typed)",
			input:    []string{"Italian", "Mediterranean"},
			expected: []string{"Italian", "Mediterranean"},
		},
		{
			name:     "array with empty strings filtered",
			input:    []interface{}{"Main Course", "", "Dinner"},
			expected: []string{"Main Course", "Dinner"},
		},
		{
			name:     "array with only empty strings",
			input:    []interface{}{"", ""},
			expected: nil,
		},
		{
			name:     "single item array",
			input:    []interface{}{"Dessert"},
			expected: []string{"Dessert"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseStringOrArray(tt.input)

			if len(tt.expected) == 0 && len(result) == 0 {
				return // Both nil/empty, OK
			}

			if len(result) != len(tt.expected) {
				t.Errorf("Length mismatch: got %d, want %d", len(result), len(tt.expected))
				return
			}

			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("Index %d: got %q, want %q", i, v, tt.expected[i])
				}
			}
		})
	}
}
