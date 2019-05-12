package estimator

// WithPrecision sets a specific estimator precision
func WithPrecision(precision int) func(*Estimator) {
	return func(e *Estimator) {
		e.precision = precision
	}
}

// WithConfidence sets a specific estimator confidence (standard deviation equivalent)
func WithConfidence(confidence float64) func(*Estimator) {
	return func(e *Estimator) {
		e.confidence = confidence
	}
}
