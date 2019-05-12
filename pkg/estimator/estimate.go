package estimator

import "math"

// histParams denotes a set of parameters that describe a PDF histogram
type histParams struct {
	step      float64
	shift     float64
	limitLow  float64
	limitHigh float64
}

// init (re)sets basic parameters of the estimator and marks it as not estimated
func (e *Estimator) init() {

	e.isEstimated = false

	// Initialize the estimator histograms
	e.pdf = make([]float64, e.precision, e.precision)
	e.bins = make([]float64, e.precision+1, e.precision+1)
}

// estimate performs populates the underlying PDF histograms and the results of
// the estimator
func (e *Estimator) estimate() {

	// Make sure the estimator is marked as complete / estimated at the end of the method
	defer func() {
		e.isEstimated = true
	}()

	// If an invalid situation has been requested, set numeric values to NaN
	if e.NSuccess > e.NTrial || e.NTrial == 0 {
		e.Mode = math.NaN()
		e.Integral = math.NaN()
		e.Variance = math.NaN()
		return
	}

	// Initialize the estimator
	e.init()

	// Obtain parameters for PDF and binning histograms
	histParams := e.getHistParams()

	// Prepare some temporary variables
	tempValue, currentBin := 0., histParams.limitLow
	nTrial := float64(e.NTrial)
	nSuccess := float64(e.NSuccess)

	switch {

	// Special case: Number of successes is 0
	case e.NSuccess == 0:
		for i := 0; i < e.precision; i++ {

			// Calculate binomial probability for current bin
			tempValue = math.Exp(nTrial*math.Log(1.-currentBin) - histParams.shift)

			// Update running variables
			e.pdf[i] = tempValue
			e.bins[i] = currentBin
			e.Integral += tempValue
			currentBin += histParams.step
		}

	// Special case: Number of successes equals number of trials
	case e.NSuccess == e.NTrial:
		for i := 0; i < e.precision; i++ {

			// Calculate binomial probability for current bin
			tempValue = math.Exp(nSuccess*math.Log(currentBin) - histParams.shift)

			// Update running variables
			e.pdf[i] = tempValue
			e.bins[i] = currentBin
			e.Integral += tempValue
			currentBin += histParams.step
		}

	// Default case
	default:
		for i := 0; i < e.precision; i++ {

			// Calculate binomial probability for current bin
			tempValue = math.Exp(nSuccess*math.Log(currentBin) +
				(nTrial-nSuccess)*math.Log(1.-currentBin) - histParams.shift)

			// Update running variables
			e.pdf[i] = tempValue
			e.bins[i] = currentBin
			e.Integral += tempValue
			currentBin += histParams.step
		}
	}

	// Update rightmost boundary bin of the PDF
	e.bins[e.precision] = currentBin
}

// getHistParams defines an optimal set of parameters for the PDF histogram in
// oder to achieve a semi-constant absolute precision (essentially "zooming" in)
// on the relevant part of the PDF
func (e *Estimator) getHistParams() (params histParams) {

	// Floating point equivalents of nTrial / nSuccess to avoid clutter from repeated
	// conversions below
	nTrial := float64(e.NTrial)
	nSuccess := float64(e.NSuccess)

	// Determine lower / upper limit on PDF range based on 5 * classical variance
	// The factor 5 is arbitrary / based on experience
	params.limitLow = math.Max(0., e.Mode-5.*e.Variance)
	params.limitHigh = math.Min(1., e.Mode+5.*e.Variance)

	// Special case: Number of successes is 0 (need more range towards upper part of pdf)
	if e.NSuccess == 0 {
		params.limitLow = 0.
		params.limitHigh = math.Min(1., (8. / nTrial))
	}

	// Special case : Number of successes equals number of trials (need more range towards lower part of pdf)
	if e.NSuccess == e.NTrial {
		params.limitLow = math.Max(0., (1. - 8./nTrial))
		params.limitHigh = 1.
	}

	// Determine step size / bin width for the PDF histogram
	params.step = math.Abs(params.limitHigh-params.limitLow) / float64(e.precision)
	if e.NSuccess != 0 && e.NSuccess != e.NTrial {
		params.shift = nSuccess*math.Log(e.Mode) + (nTrial-nSuccess)*math.Log(1.-e.Mode)
	}

	return
}
