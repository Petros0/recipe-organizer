package handler

// FetchStrategy defines the interface for recipe fetching strategies
type FetchStrategy interface {
	// Fetch attempts to fetch and parse a recipe from the given URL
	Fetch(url string) (*Recipe, error)

	// CanRetry determines if the given error is retryable by the next strategy
	CanRetry(err error) bool

	// Name returns the strategy name for logging purposes
	Name() string
}

// StrategyExecutor executes fetch strategies in order until one succeeds
type StrategyExecutor struct {
	strategies []FetchStrategy
}

// NewStrategyExecutor creates a new executor with the given strategies
func NewStrategyExecutor(strategies ...FetchStrategy) *StrategyExecutor {
	return &StrategyExecutor{strategies: strategies}
}

// Execute tries each strategy in order until one succeeds or all fail
// The logger function is called to log progress messages
func (e *StrategyExecutor) Execute(url string, logger func(...interface{})) (*Recipe, error) {
	var lastErr error

	for i, strategy := range e.strategies {
		logger("Attempting to fetch recipe using " + strategy.Name() + "...")

		recipe, err := strategy.Fetch(url)
		if err == nil {
			return recipe, nil
		}

		lastErr = err

		// Check if we should try the next strategy
		isLastStrategy := i == len(e.strategies)-1
		if !isLastStrategy && strategy.CanRetry(err) {
			logger(strategy.Name() + " failed with retryable error, trying next strategy...")
			continue
		}

		// Either it's the last strategy or the error is not retryable
		break
	}

	return nil, lastErr
}
