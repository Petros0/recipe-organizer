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
func (e *StrategyExecutor) Execute(url string, logger *Logger) (*Recipe, error) {
	var lastErr error

	for i, strategy := range e.strategies {
		if logger != nil {
			logger.Info("strategy", "Attempting to fetch recipe", map[string]interface{}{
				"strategy": strategy.Name(),
			})
		}

		recipe, err := strategy.Fetch(url)
		if err == nil {
			if logger != nil {
				logger.Info("strategy", "Recipe fetched successfully", map[string]interface{}{
					"strategy": strategy.Name(),
				})
			}
			return recipe, nil
		}

		lastErr = err

		// Check if we should try the next strategy
		isLastStrategy := i == len(e.strategies)-1
		if !isLastStrategy && strategy.CanRetry(err) {
			if logger != nil {
				logger.Info("strategy", "Strategy failed with retryable error, trying next", map[string]interface{}{
					"strategy": strategy.Name(),
					"error":    err.Error(),
				})
			}
			continue
		}

		// Either it's the last strategy or the error is not retryable
		if logger != nil {
			logger.Error("strategy", "Strategy failed", map[string]interface{}{
				"strategy": strategy.Name(),
				"error":    err.Error(),
			})
		}
		break
	}

	return nil, lastErr
}
