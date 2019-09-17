package effuncert

// WithConfidence sets a specific estimator confidence (standard deviation equivalent)
func WithConfidence(confidence float64) func(*Estimator) {
	return func(e *Estimator) {
		e.confidence = confidence
	}
}
