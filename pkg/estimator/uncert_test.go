package estimator

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
	testCase{nTrial: 1, nSuccess: 0, mode: 0.0000000000, lowInterval: 0.0000000000, highInterval: 0.4365000000},
	testCase{nTrial: 1000, nSuccess: 0, mode: 0.0000000000, lowInterval: 0.0000000000, highInterval: 0.0011480000},
	testCase{nTrial: 34, nSuccess: 12, mode: 0.3529411765, lowInterval: 0.0764537806, highInterval: 0.0829555041},
	testCase{nTrial: 34, nSuccess: 20, mode: 0.5882352941, lowInterval: 0.0830851394, highInterval: 0.0803362080},
	testCase{nTrial: 34, nSuccess: 23, mode: 0.6764705882, lowInterval: 0.0812066670, highInterval: 0.0746004179},
	testCase{nTrial: 34, nSuccess: 24, mode: 0.7058823529, lowInterval: 0.0794571706, highInterval: 0.0725751108},
	testCase{nTrial: 19, nSuccess: 10, mode: 0.5263157895, lowInterval: 0.1088157895, highInterval: 0.1081842105},
	testCase{nTrial: 19, nSuccess: 15, mode: 0.7894736842, lowInterval: 0.0990583653, highInterval: 0.0806567289},
	testCase{nTrial: 19, nSuccess: 18, mode: 0.9473684211, lowInterval: 0.0685609070, highInterval: 0.0370386604},
	testCase{nTrial: 19, nSuccess: 19, mode: 1.0000000000, lowInterval: 0.0557894737, highInterval: 0.0000000000},
	testCase{nTrial: 40, nSuccess: 17, mode: 0.4250000000, lowInterval: 0.0746451789, highInterval: 0.0777716785},
	testCase{nTrial: 40, nSuccess: 29, mode: 0.7250000000, lowInterval: 0.0725982940, highInterval: 0.0655618229},
	testCase{nTrial: 40, nSuccess: 34, mode: 0.8500000000, lowInterval: 0.0611735403, highInterval: 0.0499249241},
	testCase{nTrial: 40, nSuccess: 34, mode: 0.8500000000, lowInterval: 0.0611735403, highInterval: 0.0499249241},
	testCase{nTrial: 40, nSuccess: 36, mode: 0.9000000000, lowInterval: 0.0535813106, highInterval: 0.0404893495},
	testCase{nTrial: 340, nSuccess: 200, mode: 0.5882352941, lowInterval: 0.0265572879, highInterval: 0.0265572879},
	testCase{nTrial: 340, nSuccess: 230, mode: 0.6764705882, lowInterval: 0.0254981074, highInterval: 0.0249906824},
	testCase{nTrial: 340, nSuccess: 240, mode: 0.7058823529, lowInterval: 0.0248343791, highInterval: 0.0243401626},
	testCase{nTrial: 190, nSuccess: 100, mode: 0.5263157895, lowInterval: 0.0356801839, highInterval: 0.0360424192},
	testCase{nTrial: 190, nSuccess: 150, mode: 0.7894736842, lowInterval: 0.0303158039, highInterval: 0.0285412203},
	testCase{nTrial: 190, nSuccess: 180, mode: 0.9473684211, lowInterval: 0.0177245476, highInterval: 0.0147475108},
	testCase{nTrial: 190, nSuccess: 190, mode: 1.0000000000, lowInterval: 0.0060000000, highInterval: 0.0000000000},
	testCase{nTrial: 400, nSuccess: 170, mode: 0.4250000000, lowInterval: 0.0243463927, highInterval: 0.0248407357},
	testCase{nTrial: 400, nSuccess: 290, mode: 0.7250000000, lowInterval: 0.0224373424, highInterval: 0.0219908282},
	testCase{nTrial: 400, nSuccess: 340, mode: 0.8500000000, lowInterval: 0.0182999103, highInterval: 0.0172286961},
	testCase{nTrial: 400, nSuccess: 340, mode: 0.8500000000, lowInterval: 0.0182999103, highInterval: 0.0172286961},
	testCase{nTrial: 400, nSuccess: 360, mode: 0.9000000000, lowInterval: 0.0155250000, highInterval: 0.0143250000},
	testCase{nTrial: 1, nSuccess: 1, mode: 1.0000000000, lowInterval: 0.4365000000, highInterval: 0.0000000000},
	testCase{nTrial: 1000, nSuccess: 1000, mode: 1.0000000000, lowInterval: 0.0011480000, highInterval: 0.0000000000},
	testCase{nTrial: 1000000000, nSuccess: 0, mode: 0.0000000000, lowInterval: 0.0000000000, highInterval: 0.0000000011},
	testCase{nTrial: 1000000000, nSuccess: 1, mode: 0.0000000010, lowInterval: 0.0000000007, highInterval: 0.0000000015},
	testCase{nTrial: 1000000000, nSuccess: 500000000, mode: 0.5000000000, lowInterval: 0.0000157323, highInterval: 0.0000158904},
	testCase{nTrial: 1000000000, nSuccess: 1000000000, mode: 1.0000000000, lowInterval: 0.0000000011, highInterval: 0.0000000000},
}

var invalidCases = []testCase{
	testCase{nTrial: 10, nSuccess: 11, mode: math.NaN(), lowInterval: math.NaN(), highInterval: math.NaN()},
	testCase{nTrial: 0, nSuccess: 0, mode: math.NaN(), lowInterval: math.NaN(), highInterval: math.NaN()},
	testCase{nTrial: 0, nSuccess: 1e9, mode: math.NaN(), lowInterval: math.NaN(), highInterval: math.NaN()},
}

func TestStringerInterface(t *testing.T) {
	estimator := New(1, 2)
	expectedString := "(0.5000 -0.2475 +0.2475)"

	if fmt.Sprintf("%s", estimator) != expectedString {
		t.Fatalf("Unexpected formatted string, want \"%s\", have \"%s\"", expectedString, fmt.Sprintf("%s", estimator))
	}
}

func TestOptions(t *testing.T) {

	expectedPrecision := 100004
	expectedConfidence := TwoSigma

	estimator := New(1, 42,
		WithPrecision(expectedPrecision),
		WithConfidence(expectedConfidence),
	)

	if estimator.precision != expectedPrecision {
		t.Fatalf("Unexpected estimator precision, want %d, have %d", expectedPrecision, estimator.precision)
	}
	if estimator.confidence != expectedConfidence {
		t.Fatalf("Unexpected estimator confidence, want %.5f, have %.5f", expectedConfidence, estimator.confidence)
	}

	estimator.estimate()
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
			t.Fatal(err)
		}
	}
}

func TestLoop(t *testing.T) {

	maxTrialsFine, maxTrialsCoarse := uint64(250), uint64(1000000)

	for nTrial := uint64(1); nTrial < maxTrialsFine; nTrial++ {
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

	for nTrial := uint64(1); nTrial < maxTrialsCoarse; nTrial += 11111 {
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
		0.01:         0.1814949164,
		0.02:         0.1892169694,
		0.03:         0.1942530910,
		0.04:         0.1982819882,
		0.05:         0.2013036611,
		0.06:         0.2039895926,
		0.07:         0.2063397827,
		0.08:         0.2086899728,
		0.09:         0.2107044214,
		0.10:         0.2123831286,
		0.11:         0.2140618358,
		0.12:         0.2157405429,
		0.13:         0.2170835087,
		0.14:         0.2187622159,
		0.15:         0.2201051816,
		0.16:         0.2214481474,
		0.17:         0.2224553717,
		0.18:         0.2237983374,
		0.19:         0.2251413032,
		0.20:         0.2261485275,
		0.21:         0.2271557518,
		0.22:         0.2284987175,
		0.23:         0.2295059419,
		0.24:         0.2305131662,
		0.25:         0.2315203905,
		0.26:         0.2325276148,
		0.27:         0.2335348391,
		0.28:         0.2345420634,
		0.29:         0.2355492877,
		0.30:         0.2362207706,
		0.31:         0.2372279949,
		0.32:         0.2382352192,
		0.33:         0.2392424435,
		0.34:         0.2399139264,
		0.35:         0.2409211507,
		0.36:         0.2415926336,
		0.37:         0.2425998579,
		0.38:         0.2436070822,
		0.39:         0.2442785651,
		0.40:         0.2452857894,
		0.41:         0.2459572723,
		0.42:         0.2469644966,
		0.43:         0.2476359795,
		0.44:         0.2486432038,
		0.45:         0.2493146866,
		0.46:         0.2503219110,
		0.47:         0.2509933938,
		0.48:         0.2520006181,
		0.49:         0.2526721010,
		0.50:         0.2536793253,
		0.51:         0.2543508082,
		0.52:         0.2553580325,
		0.53:         0.2560295154,
		0.54:         0.2570367397,
		0.55:         0.2577082226,
		0.56:         0.2587154469,
		0.57:         0.2597226712,
		0.58:         0.2603941541,
		0.59:         0.2614013784,
		0.60:         0.2620728613,
		0.61:         0.2630800856,
		0.62:         0.2640873099,
		0.63:         0.2647587928,
		0.64:         0.2657660171,
		0.65:         0.2667732414,
		0.66:         0.2677804657,
		0.67:         0.2684519486,
		0.68:         0.2694591729,
		0.69:         0.2704663972,
		0.70:         0.2714736215,
		0.71:         0.2724808458,
		0.72:         0.2734880701,
		0.73:         0.2744952944,
		0.74:         0.2755025187,
		0.75:         0.2765097430,
		0.76:         0.2778527088,
		0.77:         0.2788599331,
		0.78:         0.2802028989,
		0.79:         0.2812101232,
		0.80:         0.2825530889,
		0.81:         0.2838960547,
		0.82:         0.2852390204,
		0.83:         0.2865819862,
		0.84:         0.2879249519,
		0.85:         0.2892679177,
		0.86:         0.2909466248,
		0.87:         0.2926253320,
		0.88:         0.2943040392,
		0.89:         0.2959827464,
		0.90:         0.2979971950,
		0.91:         0.3003473851,
		0.92:         0.3023618337,
		0.93:         0.3050477652,
		0.94:         0.3080694381,
		0.95:         0.3110911111,
		0.96:         0.3151200083,
		0.97:         0.3198203884,
		0.98:         0.3261994757,
		0.99:         0.3362717188,
		1.00:         1.,
		1.000000001:  math.NaN(),
		1000.:        math.NaN(),
	}

	estimator := New(42, 167)

	for quantile, value := range expectedQuantiles {

		// First attempt to trigger estimation
		estimatedQuantile := estimator.Quantile(quantile)
		if !almostEqual(estimatedQuantile, value) {
			t.Fatalf("Unexpected %.2f quantile on first attempt, want %.10f, have %.10f", quantile, value, estimatedQuantile)
		}

		// Second attempt to validate result
		estimatedQuantile = estimator.Quantile(quantile)
		if !almostEqual(estimatedQuantile, value) {
			t.Fatalf("Unexpected %.2f quantile on second attempt, want %.10f, have %.10f", quantile, value, estimatedQuantile)
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

	if printDetails {
		fmt.Printf("Result for %d/%d: %s\n", cs.nSuccess, cs.nTrial, estimator)
	}

	return nil
}

func almostEqual(a, b float64) bool {
	if math.IsNaN(a) && math.IsNaN(b) {
		return true
	}

	return math.Abs(a-b) <= epsilon
}
