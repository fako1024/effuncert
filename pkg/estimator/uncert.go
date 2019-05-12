package estimator

import (
	"fmt"
	"math"
)

const epsilon = 1e-9

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
	NSuccess, NTrial uint64  // Number of successes & trials
	Mode             float64 // Mode / expectation value of the estimator
	Integral         float64 // Integral of the PDF
	Variance         float64 // Variance / classical uncertainty of the estimator

	precision       int       // Precision of the estimation (i.e. number of bins in PDF)
	precisionDigits int       // Number of digits for adaptive string precision
	confidence      float64   // Confidence interval for the uncertainty estimation
	isEstimated     bool      // Flag indicating if the estimation has been completed
	pdf, bins       []float64 // value and binning slices for the PDF histogram
}

// New instantiates a new estimator based on a set of trails / successes
// and functional options (if any)
func New(nSuccess, nTrial uint64, options ...func(*Estimator)) *Estimator {

	mode := float64(nSuccess) / float64(nTrial)

	obj := &Estimator{
		NSuccess:   nSuccess,                                       // Number of successful trials
		NTrial:     nTrial,                                         // Total number of trials
		Mode:       mode,                                           // Mode / Classical result
		Variance:   math.Sqrt(mode * (1 - mode) / float64(nTrial)), // Variance / Classical uncertainty
		confidence: OneSigma,                                       // Sigma confidence
		precision:  1000,                                           // Precision (number of bins in PDF)
	}

	// Execute functional options (if any), see options.go for implementation
	for _, option := range options {
		option(obj)
	}

	// Set the precision digits based on the set precision
	obj.precisionDigits = len(fmt.Sprintf("%d", obj.precision))

	return obj
}

// String returns a human-readable string representing the estimator result
func (e *Estimator) String() string {
	lowInterval, highInterval := e.IntervalRelative()

	return fmt.Sprintf("(%.[4]*[1]f -%.[4]*[2]f +%.[4]*[3]f)", e.Mode, lowInterval, highInterval, e.precisionDigits)
}

// Quantile returns a quantile based on a probability
func (e *Estimator) Quantile(confidence float64) float64 {

	// If the result has not yet been determined, do it now
	if !e.isEstimated {
		e.estimate()
	}

	// Handle numerically impossible cases
	if confidence < 0. || confidence > 1. {
		return math.NaN()
	}

	// Handle special, numerically unstable cases
	if e.NTrial == 0 || e.Integral < epsilon || confidence < epsilon {
		return 0.
	}
	if (1. - confidence) < epsilon {
		return 1.
	}

	// Prepare some temporary variables
	currentIntegral := 0.
	interval := confidence * e.Integral

	// Find the bin for the requested quantile in the PDF histogram
	for i := 0; i < e.precision; i++ {
		currentIntegral += e.pdf[i]
		if currentIntegral > interval {

			// Return the bin center of the found bin
			return (e.bins[i] + e.bins[i+1]) / 2.
		}
	}

	// If none is found, return NaN (should not even happen)
	return math.NaN()
}

// Interval returns the absolute lower and upper quantiles for the uncertainty estimation
func (e *Estimator) Interval() (lowQuantile float64, highQuantile float64) {

	// If a matching lookup value exists, return it right away
	if e.confidence == OneSigma && e.NTrial <= maxLookupNTrial {
		if lookupResult, ok := lookupTable[experiment{e.NSuccess, e.NTrial}]; ok {
			return lookupResult.lowInterval, lookupResult.highInterval
		}
	}

	// If the result has not yet been determined, do it now
	if !e.isEstimated {
		e.estimate()
	}

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

		// Prepare some temporary variables
		lowBin, highBin := -1, -1
		currentIntegral := 0.
		interval := e.confidence * e.Integral

		// Find the bin corresponding to the mode in the PDF histogram and the respective
		// integral up until that point
		for i := 0; i < e.precision; i++ {
			if e.Mode >= e.bins[i] && e.Mode <= e.bins[i+1] {
				highBin, lowBin = i, i
				currentIntegral += e.pdf[i]
				break
			}
		}

		// Expand symmetrically outwards from the mode until the requested confidence
		// interval is reached
		for currentIntegral < interval {
			if e.pdf[highBin+1] >= e.pdf[lowBin-1] {
				highBin++
				currentIntegral += e.pdf[highBin]
			} else {
				lowBin--
				currentIntegral += e.pdf[lowBin]
			}
		}

		// Set low and high quantiles of the confidence interval boundaries accordingly
		lowQuantile = (e.bins[lowBin] + e.bins[lowBin+1]) / 2.
		highQuantile = (e.bins[highBin] + e.bins[highBin+1]) / 2.
	}

	return
}

// IntervalRelative returns the relative lower and upper quantiles for the uncertainty estimation
func (e *Estimator) IntervalRelative() (lowQuantile float64, highQuantile float64) {

	// Determine the absolute interval quantiles
	tempLow, tempHigh := e.Interval()

	return e.Mode - tempLow, tempHigh - e.Mode
}
