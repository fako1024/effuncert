package effuncert

import (
	"fmt"
	"math"

	"github.com/fako1024/numerics"
	"github.com/fako1024/numerics/root"
)

const (
	epsilon                  = 1e-9
	maxQuadraticRootFindingN = 1000
)

const (

	// OneSigma denotes a one sigma standard deviation equivalent
	OneSigma = 0.6826895475

	// TwoSigma denotes a two sigma standard deviation equivalent
	TwoSigma = 0.9544997215

	// ThreeSigma denotes a three sigma standard deviation equivalent
	ThreeSigma = 0.9973001480

	// FourSigma denotes a four sigma standard deviation equivalent
	FourSigma = 0.99993669986724854

	// FiveSigma denotes a five sigma standard deviation equivalent
	FiveSigma = 0.99999940395355225
)

// Estimator denotes a numeric estimator instance for a Bernoulli experiment and
// its uncertainty based on a binomial probability distribution
type Estimator struct {
	NSuccess, NTrial float64 // Number of successes & trials
	Mode             float64 // Mode / expectation value of the estimator
	Integral         float64 // Integral of the PDF
	Variance         float64 // Variance / classical uncertainty of the estimator

	confidence                float64 // Confidence interval for the uncertainty estimation
	lowInterval, highInterval float64 // Valus holding the low / high relative uncertainty intervals
	isEstimated               bool    // Inficator if interval has been estimated
}

// New instantiates a new estimator based on a set of trails / successes
// and functional options (if any)
func New(nSuccess, nTrial uint64, options ...func(*Estimator)) *Estimator {

	obj := &Estimator{
		NSuccess:   float64(nSuccess), // Number of successful trials
		NTrial:     float64(nTrial),   // Total number of trials
		confidence: OneSigma,          // Sigma confidence
	}

	// Determine mode and classical variance
	obj.Mode = obj.NSuccess / obj.NTrial                             // Mode / Classical result
	obj.Variance = math.Sqrt(obj.Mode * (1 - obj.Mode) / obj.NTrial) // Variance / Classical uncertainty

	// Execute functional options (if any), see options.go for implementation
	for _, option := range options {
		option(obj)
	}

	return obj
}

// String returns a human-readable string representing the estimator result
func (e *Estimator) String() string {

	// Calculate / get low and high intervals
	lowInterval, highInterval := e.IntervalRelative()

	// Set the precision digits based on the set precision
	precisionDigits := len(fmt.Sprintf("%0.3f", math.Min(lowInterval, highInterval)))

	return fmt.Sprintf("(%.[4]*[1]f -%.[4]*[2]f +%.[4]*[3]f)", e.Mode, lowInterval, highInterval, precisionDigits)
}

// Quantile returns a quantile based on a probability
func (e *Estimator) Quantile(confidence float64) float64 {

	// Handle numerically impossible cases
	if confidence < 0. || confidence > 1. {
		return math.NaN()
	}

	// Handle special, numerically unstable cases
	if e.NTrial == 0 || confidence < epsilon {
		return 0.
	}
	if (1. - confidence) < epsilon {
		return 1.
	}

	// Determine the initial result seed based on the mode of the distribution and
	// stabilize edge cases
	initialEstimate := e.Mode
	if initialEstimate < epsilon {
		initialEstimate = math.Min(0.1, 1./float64(e.NTrial))
	} else if initialEstimate == 1 {
		initialEstimate = math.Max(0.9, 1.-(1./float64(e.NTrial)))
	}

	// For large values use a linear root finding method (as it is more stable)
	if e.NTrial > maxQuadraticRootFindingN {
		return root.Bisect(func(x float64) float64 {
			return numerics.BetaIncompleteRegular(x, 1.+float64(e.NSuccess), 1.-float64(e.NSuccess)+float64(e.NTrial)) - confidence
		}, 0., 1.)
	}

	// For smaller values use a quadratic root finding method (as it is faster and more precise)
	return root.Find(func(x float64) float64 {
		return numerics.BetaIncompleteRegular(x, 1.+float64(e.NSuccess), 1.-float64(e.NSuccess)+float64(e.NTrial)) - confidence
	}, func(x float64) float64 {
		return numerics.Binomial(x, float64(e.NSuccess), float64(e.NTrial)) / numerics.Beta(1.+float64(e.NSuccess), 1.-float64(e.NSuccess)+float64(e.NTrial))
	}, initialEstimate, root.WithLimits(0., 1.), root.WithHeuristics())
}

// Interval returns the absolute lower and upper quantiles for the uncertainty estimation
func (e *Estimator) Interval() (lowQuantile float64, highQuantile float64) {

	// If the result is not valid, return accordingly
	if math.IsNaN(e.Mode) || math.IsNaN(e.Integral) || math.IsNaN(e.Variance) {
		return math.NaN(), math.NaN()
	}

	switch {

	// Special case: Number of successes is 0
	case e.NSuccess == 0:
		lowQuantile = 0.
		highQuantile = e.Quantile(e.confidence)

	// Special case: Number of successes equals number of trials
	case e.NSuccess == e.NTrial:
		lowQuantile = e.Quantile(1. - e.confidence)
		highQuantile = 1.

	// Default case
	default:
		lowQuantile = e.Quantile(0.5 * (1. - e.confidence))
		highQuantile = e.Quantile(1.0 - 0.5*(1.-e.confidence))
	}

	return
}

// IntervalRelative returns the relative lower and upper quantiles for the uncertainty estimation
func (e *Estimator) IntervalRelative() (lowQuantile float64, highQuantile float64) {

	// Check if estimation has to be performed
	if !e.isEstimated {

		// Determine the absolute interval quantiles
		tempLow, tempHigh := e.Interval()

		//fmt.Println("Interval", e.Mode, tempLow, tempHigh)
		e.lowInterval, e.highInterval = e.Mode-tempLow, tempHigh-e.Mode
	}

	return e.lowInterval, e.highInterval
}
