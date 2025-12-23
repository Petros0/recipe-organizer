# Feature Specification: URL-Based Recipe Import

**Feature Branch**: `001-url-recipe-import`  
**Created**: 2025-12-17  
**Status**: Done  
**Input**: URL-based recipe import feature for extracting recipe data from websites using schema.org Recipe structured data with LLM fallback

## Clarifications

### Session 2025-12-17

- Q: How should real-time status updates for extraction progress be delivered to the user? → A: Appwrite real-time subscription
- Q: Multiple recipes on page handling? → A: Skipped (not a current concern)
- Q: What UI pattern should display extraction progress? → A: Animated skeleton/shimmer preview of recipe card
- Q: What layout should the recipe preview use? → A: Full-detail scrollable view (hero image → metadata → ingredients → instructions)
- Q: What UI pattern should display extraction errors? → A: Inline error banner within skeleton card area (shimmer transforms to error state)

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Import Recipe from URL with Structured Data (Priority: P1)

As a user, I want to paste a recipe URL from a popular cooking website and have the app automatically extract all recipe details so I can save the recipe to my collection without manual data entry.

**Why this priority**: This is the core feature of the app. Users who discover recipes online need a quick way to save them. Most major recipe websites (AllRecipes, Food Network, BBC Good Food, etc.) use schema.org Recipe markup, making this the primary extraction path.

**Independent Test**: Can be tested by pasting a URL from a major recipe website with JSON-LD markup and verifying all recipe data is correctly extracted and displayed.

**Acceptance Scenarios**:

1. **Given** a user has a recipe URL from a website with schema.org Recipe data, **When** they submit the URL in the app, **Then** the system extracts and displays the recipe name, image(s), ingredients, instructions, cooking times, servings, and nutrition information.
2. **Given** a user submits a valid recipe URL, **When** the extraction completes successfully, **Then** the user sees a preview of the extracted recipe data in a full-detail scrollable view (hero image → metadata → ingredients → instructions) before saving.
3. **Given** a user is viewing extracted recipe data, **When** they confirm the import, **Then** the recipe is saved to their collection with all extracted fields.

---

### User Story 2 - Import Recipe with LLM Fallback (Priority: P2)

As a user, I want to import recipes from websites that don't have structured data, so I can save recipes from any source including blogs and personal cooking sites.

**Why this priority**: Many recipe blogs and smaller cooking sites don't implement schema.org markup. This fallback ensures the app works universally, expanding the range of importable recipes significantly.

**Independent Test**: Can be tested by submitting a URL from a website without JSON-LD recipe markup and verifying the LLM successfully extracts recipe components from the page content.

**Acceptance Scenarios**:

1. **Given** a user submits a URL from a website without schema.org Recipe data, **When** the primary extraction fails, **Then** the system automatically attempts LLM-based extraction from the page content.
2. **Given** LLM extraction is processing, **When** the user is waiting, **Then** they see a progress indicator with estimated completion time.
3. **Given** LLM extraction completes, **When** the user views the results, **Then** they can review and edit the extracted data before saving.

---

### User Story 3 - Handle Failed or Partial Imports (Priority: P3)

As a user, when a recipe import fails or returns incomplete data, I want clear feedback and the option to manually fill in missing information so I don't lose the recipe entirely.

**Why this priority**: Robust error handling ensures users aren't frustrated when websites have bot protection or unusual layouts. Allowing manual completion preserves user effort.

**Independent Test**: Can be tested by submitting a URL that triggers an extraction error and verifying the user receives actionable error messages with recovery options.

**Acceptance Scenarios**:

1. **Given** a user submits a URL that cannot be accessed (404, bot protection, etc.), **When** the extraction fails, **Then** the user sees an inline error banner within the skeleton card area (shimmer transforms to error state) with a clear message explaining the issue and suggested actions.
2. **Given** extraction returns partial data (e.g., missing nutrition info), **When** the preview is shown, **Then** missing fields are clearly indicated and the user can manually enter the missing information.
3. **Given** a user has partially extracted data, **When** they choose to save anyway, **Then** the recipe is saved with available data and missing fields marked as incomplete.

---

### User Story 4 - Source Attribution (Priority: P4)

As a user, I want imported recipes to retain attribution to the original source so I can reference the original website and give proper credit.

**Why this priority**: Attribution is important for copyright respect and allows users to return to the original source. It's essential but not blocking for core functionality.

**Independent Test**: Can be tested by importing a recipe and verifying the saved recipe includes source URL and author information.

**Acceptance Scenarios**:

1. **Given** a recipe is successfully imported, **When** the user views the saved recipe, **Then** they see the original source URL and author name (if available).
2. **Given** a saved recipe has source attribution, **When** the user taps on the source, **Then** they are offered the option to open the original URL in their browser.

---

### Edge Cases

**Deferred to post-MVP:**
- What happens when the URL points to a non-recipe page (e.g., homepage, category listing)?
- How does the system handle websites with aggressive bot protection that block the headless browser?
- What happens when a recipe page contains multiple recipes (e.g., main dish with side dish)?
- How does the system handle non-English recipe websites?
- How does the system handle very long recipes with 50+ ingredients or steps?
- What happens when recipe images are blocked or CDN-protected?

**Handled in MVP:**
- What happens when the user submits a malformed or invalid URL? → FR-002 validates URL format; error shown via FR-010/FR-018.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST accept a URL input from the user for recipe extraction.
- **FR-002**: System MUST validate that the submitted URL is well-formed before processing.
- **FR-003**: System MUST attempt to extract recipe data using schema.org Recipe JSON-LD as the primary method.
- **FR-004**: System MUST fall back to LLM-based extraction when schema.org data is not available.
- **FR-005**: System MUST extract the following recipe fields when available: name, image(s), ingredients list, instructions/steps, prep time, cook time, total time, servings/yield, nutrition information, and author/source. **Required fields**: name, ingredients list (minimum 1 ingredient). All other fields are optional.
- **FR-006**: System MUST handle websites with bot protection by using a headless browser with appropriate headers.
- **FR-007**: System MUST display extracted recipe data for user review before saving.
- **FR-008**: System MUST allow users to edit extracted data before saving. Note: Editing is available in LLM fallback flow (US2) and partial import flow (US3). US1 preview is read-only with option to proceed to edit if needed.
- **FR-009**: System MUST preserve the original source URL and author attribution with saved recipes.
- **FR-010**: System MUST provide clear, user-friendly error messages when extraction fails.
- **FR-011**: System MUST indicate which fields could not be extracted (partial success scenario).
- **FR-012**: System MUST handle network timeouts gracefully with appropriate feedback.
- **FR-015**: System MUST provide real-time extraction status updates via Appwrite real-time subscription.
- **FR-013**: System MUST support extraction from both HTTP and HTTPS URLs.
- **FR-014**: System MUST sanitize extracted content to prevent display of malicious scripts or markup.
- **FR-016**: System MUST display extraction progress using animated skeleton/shimmer preview of the recipe card structure.
- **FR-017**: System MUST display the recipe preview in a full-detail scrollable view with vertical layout: hero image → metadata (timing, servings) → ingredients list → instructions.
- **FR-018**: System MUST display extraction errors as an inline error banner within the skeleton card area, transforming the shimmer animation into an error state.

### Key Entities

- **Recipe**: The core entity representing a saved recipe with name, description, images, ingredients, instructions, timing information, servings, nutrition facts, source URL, and author.
- **Ingredient**: Individual ingredient items with quantity, unit, and ingredient name parsed from the ingredients list.
- **Instruction Step**: Individual cooking steps with step number and instruction text.
- **Nutrition Information**: Nutritional data including calories, fat, protein, carbohydrates, and other available macro/micronutrients.
- **Recipe Source**: Attribution information including the original URL, website name, and author name.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can successfully import a recipe from a URL in under 30 seconds for sites with structured data (tested on 4G mobile network, mid-tier device).
- **SC-002**: 95% of major recipe websites with schema.org markup result in complete recipe extraction.
- **SC-003**: 80% of recipe pages without structured data are successfully parsed via LLM fallback.
- **SC-004**: Users can complete the full import flow (paste URL → review → save) in under 2 minutes.
- **SC-005**: Less than 5% of import attempts result in complete failure with no usable data extracted.
- **SC-006**: 90% of users successfully import their first recipe on the first attempt.

## Assumptions

- Popular recipe websites (AllRecipes, Food Network, BBC Good Food, Serious Eats, etc.) consistently implement schema.org Recipe markup.
- The LLM service has reasonable rate limits sufficient for user import volume.
- Users import recipes one at a time (batch import is out of scope for this feature).
- Recipe pages typically contain a single primary recipe per URL.
- Standard HTTP timeouts of 30 seconds are acceptable for user experience.
- Mobile networks may have higher latency; the UI should remain responsive during extraction.

## Out of Scope

- Importing recipes from social media platforms (Instagram, TikTok, YouTube Shorts) - noted for future roadmap.
- Batch importing multiple recipes at once.
- Optical character recognition (OCR) from images of recipes.
- Recipe editing features beyond the import preview (full editing is a separate feature).
- Recipe synchronization across devices (handled by core app infrastructure).
- User accounts and authentication (handled by core app infrastructure).
