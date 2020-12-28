package effuncert

import (
	"flag"
	"fmt"
	"math"
	"os"
	"testing"
)

type testCase struct {
	nTrial, nSuccess          uint64
	mode                      float64
	lowInterval, highInterval float64
}

var printDetails bool

var testCases = []testCase{
	{nTrial: 1, nSuccess: 0, mode: 0.00000000000000000000, lowInterval: 0.00000000000000000000, highInterval: 0.43669683783951662726},
	{nTrial: 1000, nSuccess: 0, mode: 0.00000000000000000000, lowInterval: 0.00000000000000000000, highInterval: 0.00114607066981073846},
	{nTrial: 34, nSuccess: 12, mode: 0.35294117647058825815, lowInterval: 0.07170524174383174909, highInterval: 0.08813397652966181717},
	{nTrial: 34, nSuccess: 20, mode: 0.58823529411764707842, lowInterval: 0.08704384738783355058, highInterval: 0.07718703154428951940},
	{nTrial: 34, nSuccess: 23, mode: 0.67647058823529415683, lowInterval: 0.08823190142073744635, highInterval: 0.06851677963473190580},
	{nTrial: 34, nSuccess: 24, mode: 0.70588235294117651630, lowInterval: 0.08800890841708530754, highInterval: 0.06500697230683705197},
	{nTrial: 19, nSuccess: 10, mode: 0.52631578947368418131, lowInterval: 0.11147421506782295708, highInterval: 0.10641671434632560267},
	{nTrial: 19, nSuccess: 15, mode: 0.78947368421052632748, lowInterval: 0.11959256558315511931, highInterval: 0.06393337736735893451},
	{nTrial: 19, nSuccess: 18, mode: 0.94736842105263152636, lowInterval: 0.10313257722118751580, highInterval: 0.01695170049169547610},
	{nTrial: 19, nSuccess: 19, mode: 1.00000000000000000000, lowInterval: 0.05577777428665087189, highInterval: 0.00000000000000000000},
	{nTrial: 40, nSuccess: 17, mode: 0.42499999999999998890, lowInterval: 0.07274121114258058629, highInterval: 0.07991723012675505666},
	{nTrial: 40, nSuccess: 29, mode: 0.72499999999999997780, lowInterval: 0.08021451331684192887, highInterval: 0.05868347400168372197},
	{nTrial: 40, nSuccess: 34, mode: 0.84999999999999997780, lowInterval: 0.07351099722083942467, highInterval: 0.04000509075583891239},
	{nTrial: 40, nSuccess: 34, mode: 0.84999999999999997780, lowInterval: 0.07351099722083942467, highInterval: 0.04000509075583891239},
	{nTrial: 40, nSuccess: 36, mode: 0.90000000000000002220, lowInterval: 0.06800251913412902471, highInterval: 0.02969462800788957857},
	{nTrial: 340, nSuccess: 200, mode: 0.58823529411764707842, lowInterval: 0.02713219283516221658, highInterval: 0.02609959668534478361},
	{nTrial: 340, nSuccess: 230, mode: 0.67647058823529415683, lowInterval: 0.02634357536204612327, highInterval: 0.02427836369847957698},
	{nTrial: 340, nSuccess: 240, mode: 0.70588235294117651630, lowInterval: 0.02586276831269918031, highInterval: 0.02345334203538562701},
	{nTrial: 190, nSuccess: 100, mode: 0.52631578947368418131, lowInterval: 0.03630936281377022956, highInterval: 0.03576054923423221954},
	{nTrial: 190, nSuccess: 150, mode: 0.78947368421052632748, lowInterval: 0.03254349006370826913, highInterval: 0.02650609061886510798},
	{nTrial: 190, nSuccess: 180, mode: 0.94736842105263152636, lowInterval: 0.02122698023950320145, highInterval: 0.01189189797565037843},
	{nTrial: 190, nSuccess: 190, mode: 1.00000000000000000000, lowInterval: 0.00599179204880251337, highInterval: 0.00000000000000000000},
	{nTrial: 400, nSuccess: 170, mode: 0.42499999999999998890, lowInterval: 0.02428413438726445550, highInterval: 0.02503077457909930192},
	{nTrial: 400, nSuccess: 290, mode: 0.72499999999999997780, lowInterval: 0.02340882827956103363, highInterval: 0.02116887018025126466},
	{nTrial: 400, nSuccess: 340, mode: 0.84999999999999997780, lowInterval: 0.01960806841929063626, highInterval: 0.01612351178325077683},
	{nTrial: 400, nSuccess: 340, mode: 0.84999999999999997780, lowInterval: 0.01960806841929063626, highInterval: 0.01612351178325077683},
	{nTrial: 400, nSuccess: 360, mode: 0.90000000000000002220, lowInterval: 0.01704196513242905997, highInterval: 0.01305938321106281386},
	{nTrial: 1, nSuccess: 1, mode: 1.00000000000000000000, lowInterval: 0.43669683783951640521, highInterval: 0.00000000000000000000},
	{nTrial: 1000, nSuccess: 1000, mode: 1.00000000000000000000, lowInterval: 0.00114607066981076144, highInterval: 0.00000000000000000000},
	{nTrial: 1000000000, nSuccess: 0, mode: 0.00000000000000000000, lowInterval: 0.00000000000000000000, highInterval: 0.00000000114232534543},
	{nTrial: 1000000000, nSuccess: 1, mode: 0.00000000100000000000, lowInterval: 0.00000000029423211142, highInterval: 0.00000000229600879923},
	{nTrial: 1000000000, nSuccess: 1000000000, mode: 1.00000000000000000000, lowInterval: 0.00000000114232534543, highInterval: 0.00000000000000000000},
	{nTrial: 1000000000, nSuccess: 500000000, mode: 0.50000000000000000000, lowInterval: math.NaN(), highInterval: math.NaN()},
}

var invalidCases = []testCase{
	{nTrial: 10, nSuccess: 11, mode: math.NaN(), lowInterval: math.NaN(), highInterval: math.NaN()},
	{nTrial: 0, nSuccess: 0, mode: math.NaN(), lowInterval: math.NaN(), highInterval: math.NaN()},
	{nTrial: 0, nSuccess: 1e9, mode: math.NaN(), lowInterval: math.NaN(), highInterval: math.NaN()},
}

func TestStringerInterface(t *testing.T) {
	estimator := New(1, 2)
	expectedString := "(0.50000 -0.24787 +0.24787)"

	if estimator.String() != expectedString {
		t.Fatalf("Unexpected formatted string, want \"%s\", have \"%s\"", expectedString, estimator)
	}
}

func TestOptions(t *testing.T) {

	expectedConfidence := TwoSigma

	estimator := New(1, 42,
		WithConfidence(expectedConfidence),
	)

	if estimator.confidence != expectedConfidence {
		t.Fatalf("Unexpected estimator confidence, want %.5f, have %.5f", expectedConfidence, estimator.confidence)
	}
}

func TestInvalid(t *testing.T) {
	for _, cs := range invalidCases {
		if err := estimate(cs); err != nil {
			t.Fatal(err)
		}
	}
}

func TestTable(t *testing.T) {
	for _, cs := range testCases {
		if err := estimate(cs); err != nil {
			t.Error(err)
		}
	}
}

func TestSymmetry(t *testing.T) {

	for nTrial := uint64(1); nTrial < uint64(100); nTrial++ {
		for nSuccess := uint64(0); nSuccess < nTrial; nSuccess++ {

			lowInterval1, highInterval1 := New(nSuccess, nTrial).IntervalRelative()
			lowInterval2, highInterval2 := New(nTrial-nSuccess, nTrial).IntervalRelative()

			if math.Abs(highInterval1-lowInterval2) > epsilon {
				t.Fatalf("Upper interval of left-sided distribution for %d/%d does not equal lower interval of right-sided distribution: %.10f vs. %.10f", nSuccess, nTrial, highInterval1, lowInterval2)
			}

			if math.Abs(highInterval2-lowInterval1) > epsilon {
				t.Fatalf("Lower interval of left-sided distribution for %d/%d does not equal upper interval of right-sided distribution: %.10f vs. %.10f", nSuccess, nTrial, lowInterval1, highInterval2)
			}
		}
	}
}

func TestLoopFine(t *testing.T) {

	maxTrials := uint64(250)
	for nTrial := uint64(1); nTrial < maxTrials; nTrial++ {
		for nSuccess := uint64(0); nSuccess < nTrial; nSuccess++ {

			estimator := New(nSuccess, nTrial)

			lowInterval, highInterval := estimator.IntervalRelative()
			if math.IsNaN(lowInterval) || math.IsNaN(highInterval) {
				t.Fatalf("Unexpected NaN for %d/%d: %.10f , %.10f", nSuccess, nTrial, lowInterval, highInterval)
			}

			if printDetails {
				fmt.Printf("Result for %d/%d: %s\n", nSuccess, nTrial, estimator)
			}
		}
	}
}

func TestLoopCoarse(t *testing.T) {

	maxTrials := uint64(100000)
	for nTrial := uint64(1); nTrial < maxTrials; nTrial += 11111 {
		for nSuccess := uint64(0); nSuccess < nTrial; nSuccess += 1111 {

			estimator := New(nSuccess, nTrial)

			lowInterval, highInterval := estimator.IntervalRelative()
			if math.IsNaN(lowInterval) || math.IsNaN(highInterval) {
				t.Fatalf("Unexpected NaN for %d/%d: %.10f , %.10f", nSuccess, nTrial, lowInterval, highInterval)
			}

			if printDetails {
				fmt.Printf("Result for %d/%d: %s\n", nSuccess, nTrial, estimator)
			}
		}
	}
}

func TestQuantile(t *testing.T) {

	expectedQuantiles := map[float64]float64{
		-1000.:       math.NaN(),
		-0.000000001: math.NaN(),
		0.00:         0.,
		0.01:         0.1812873965,
		0.02:         0.1891104197,
		0.03:         0.1941656096,
		0.04:         0.1980147777,
		0.05:         0.2011748462,
		0.06:         0.2038849390,
		0.07:         0.2062764704,
		0.08:         0.2084298630,
		0.09:         0.2103981247,
		0.10:         0.2122181507,
		0.11:         0.2139167193,
		0.12:         0.2155139213,
		0.13:         0.2170252435,
		0.14:         0.2184628943,
		0.15:         0.2198366821,
		0.16:         0.2211546167,
		0.17:         0.2224233316,
		0.18:         0.2236483900,
		0.19:         0.2248345082,
		0.20:         0.2259857249,
		0.21:         0.2271055290,
		0.22:         0.2281969594,
		0.23:         0.2292626824,
		0.24:         0.2303050536,
		0.25:         0.2313261676,
		0.26:         0.2323278980,
		0.27:         0.2333119302,
		0.28:         0.2342797882,
		0.29:         0.2352328575,
		0.30:         0.2361724037,
		0.31:         0.2370995882,
		0.32:         0.2380154821,
		0.33:         0.2389210771,
		0.34:         0.2398172960,
		0.35:         0.2407050009,
		0.36:         0.2415850008,
		0.37:         0.2424580581,
		0.38:         0.2433248943,
		0.39:         0.2441861952,
		0.40:         0.2450426153,
		0.41:         0.2458947823,
		0.42:         0.2467433002,
		0.43:         0.2475887532,
		0.44:         0.2484317087,
		0.45:         0.2492727200,
		0.46:         0.2501123293,
		0.47:         0.2509510700,
		0.48:         0.2517894694,
		0.49:         0.2526280509,
		0.50:         0.2534673364,
		0.51:         0.2543078486,
		0.52:         0.2551501132,
		0.53:         0.2559946615,
		0.54:         0.2568420322,
		0.55:         0.2576927746,
		0.56:         0.2585474504,
		0.57:         0.2594066369,
		0.58:         0.2602709298,
		0.59:         0.2611409462,
		0.60:         0.2620173277,
		0.61:         0.2629007447,
		0.62:         0.2637918997,
		0.63:         0.2646915320,
		0.64:         0.2656004230,
		0.65:         0.2665194010,
		0.66:         0.2674493482,
		0.67:         0.2683912073,
		0.68:         0.2693459898,
		0.69:         0.2703147851,
		0.70:         0.2712987714,
		0.71:         0.2722992278,
		0.72:         0.2733175492,
		0.73:         0.2743552630,
		0.74:         0.2754140490,
		0.75:         0.2764957635,
		0.76:         0.2776024679,
		0.77:         0.2787364629,
		0.78:         0.2799003304,
		0.79:         0.2810969848,
		0.80:         0.2823297365,
		0.81:         0.2836023711,
		0.82:         0.2849192498,
		0.83:         0.2862854369,
		0.84:         0.2877068665,
		0.85:         0.2891905598,
		0.86:         0.2907449174,
		0.87:         0.2923801158,
		0.88:         0.2941086598,
		0.89:         0.2959461695,
		0.90:         0.2979125336,
		0.91:         0.3000336568,
		0.92:         0.3023442147,
		0.93:         0.3048922092,
		0.94:         0.3077469660,
		0.95:         0.3110142919,
		0.96:         0.3148682820,
		0.97:         0.3196283966,
		0.98:         0.3259927614,
		0.99:         0.3361048240,
		1.00:         1.,
		1.000000001:  math.NaN(),
		1000.:        math.NaN(),
	}

	estimator := New(42, 167)

	for confidence, value := range expectedQuantiles {
		estimatedQuantile := estimator.Quantile(confidence)
		if !almostEqual(estimatedQuantile, value) {
			t.Fatalf("Unexpected %.2f quantile, want %.10f, have %.10f", confidence, value, estimatedQuantile)
		}
	}
}

func BenchmarkLookup(b *testing.B) {
	for i := 0; i < b.N; i++ {
		estimator := New(37, 163)
		lowInterval, highInterval := estimator.IntervalRelative()

		_, _ = lowInterval, highInterval
	}
}

func BenchmarkEstimate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		estimator := New(123, 1635)
		lowInterval, highInterval := estimator.IntervalRelative()

		_, _ = lowInterval, highInterval
	}
}

func TestMain(m *testing.M) {

	// Parse flags
	flag.BoolVar(&printDetails, "printDetails", false, "Print all results while running tests")
	flag.Parse()

	os.Exit(m.Run())
}

////////////////////////////////////////////////////////////////////////////////

func estimate(cs testCase) error {
	estimator := New(cs.nSuccess, cs.nTrial)
	lowInterval, highInterval := estimator.IntervalRelative()

	if !almostEqual(lowInterval, cs.lowInterval) {
		return fmt.Errorf("Unexpected low interval boundary, want %.10f, have %.10f", cs.lowInterval, lowInterval)
	}
	if !almostEqual(highInterval, cs.highInterval) {
		return fmt.Errorf("Unexpected high interval boundary, want %.10f, have %.10f", cs.highInterval, highInterval)
	}
	_, _ = lowInterval, highInterval
	if printDetails {
		fmt.Printf("%s\n", estimator)
	}

	return nil
}

func almostEqual(a, b float64) bool {
	if math.IsNaN(a) && math.IsNaN(b) {
		return true
	}

	return math.Abs(a-b) <= epsilon
}
