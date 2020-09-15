# Numeric estimation of statistical uncertainties for Bernoulli experiments
This package performs a numeric estimation of quantiles and uncertainties for a Bernoulli experiment or and estimator for an efficiency (k/n successful coin flip trials) in pure Go.

[![Github Release](https://img.shields.io/github/release/fako1024/effuncert.svg)](https://github.com/fako1024/effuncert/releases)
[![GoDoc](https://godoc.org/github.com/fako1024/effuncert?status.svg)](https://godoc.org/github.com/fako1024/effuncert/)
[![Go Report Card](https://goreportcard.com/badge/github.com/fako1024/effuncert)](https://goreportcard.com/report/github.com/fako1024/effuncert)
![Build/Test Status](https://github.com/fako1024/effuncert/workflows/Go/badge.svg)

Different Frequentist (e.g. Likelihood methods) and non-Frequentist approachs can be used to determine the uncertainties for the Trigger efficiency. In this case, a Bayesian approach was chosen to account for the asymmetry introduced by the Binomial distributed variables and to solve the problematic of boundary effects as well.
The underlying Binomial probability distribution function and quantiles are modeled using an incomplete, regularized Beta function (the corresponding integral / cumulative PDF) and applying a numeric root finding method to determine arbitrary quantiles.

## Features
- Estimation of classical mode and variance of Bernoulli experiment / efficiency calculation
- Extraction / estimation of quantiles for the underlying probability distribution function
- Extraction asymmetric uncertainty intervals equivalents for any confidence level / interval

## Installation
```bash
go get -u github.com/fako1024/effuncert
```

## API summary

The API of the package is fairly straight-forward. The following functions are exposed:
```Go
// Estimator denotes a numeric estimator instance for a Bernoulli experiment and
// its uncertainty based on a binomial probability distribution
type Estimator struct {
	NSuccess, NTrial float64 // Number of successes & trials
	Mode             float64 // Mode / expectation value of the estimator
	Integral         float64 // Integral of the PDF
	Variance         float64 // Variance / classical uncertainty of the estimator
}

// New instantiates a new estimator based on a set of trails / successes
// and functional options (if any)
func New(nSuccess, nTrial uint64, options ...func(*Estimator)) *Estimator

// String returns a human-readable string representing the estimator result
func (e *Estimator) String() string

// Quantile returns a quantile based on a probability
func (e *Estimator) Quantile(confidence float64) float64

// Interval returns the absolute lower and upper quantiles for the uncertainty estimation
func (e *Estimator) Interval() (lowQuantile float64, highQuantile float64)

// IntervalRelative returns the relative lower and upper quantiles for the uncertainty estimation
func (e *Estimator) IntervalRelative() (lowQuantile float64, highQuantile float64)
```

## Example
```Go
// 3 successful out of 8 Bernoulli experiments
nSuccess, nTrial := 3, 8

// Instantiate an estimator
e := effuncert.New(nSuccess, nTrial,
	effuncert.WithConfidence(effuncert.OneSigma)
)

// Print the result in a well-formatted way, making use of the String() method
// Will print: (0.37500 -0.12882 +0.17969)
fmt.Println(e)

// Print the median
// Will print: 0.39308483281062906
fmt.Println(e.Quantile(0.5))
```

An additional example binary code can be found in the `examples/effuncert` folder.
