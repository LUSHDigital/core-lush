package accounting_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/LUSHDigital/core-lush/accounting"
	"github.com/LUSHDigital/core-lush/currency"
	"github.com/google/go-cmp/cmp"
	"github.com/google/gofuzz"
)

func TestExchange(t *testing.T) {
	// All possible minor currency units sampled from the currency
	// package are in use in this suite.
	//
	// factors of 1, 100, 1000 and 10000 are represented.
	//
	// As there is no point in duplicating tests for currencies with
	// the exact same minor unit, the better approach is to instead
	// bulk up the amount of unique values tested.
	//
	// The test is essentially:
	// value / factor / rate = want
	//
	// Each value in the want array has been calculated manually.
	tests := []struct {
		name   string
		money  currency.Currency
		values [4]float64
		rates  [4]float64
		want   [4]float64
	}{
		{
			name:   "Japanese yen, factor 1",
			money:  currency.JPY,
			values: [4]float64{100, 1000, 10000, 1999},
			rates:  [4]float64{1, 1.5, 145.961337, 145.961337},
			want:   [4]float64{100, 666.67, 68.51, 13.70},
		},
		{
			name:   "Zimbabwean dollar, factor 100",
			money:  currency.ZWL,
			values: [4]float64{100, 1000, 10000, 19999.99},
			rates:  [4]float64{1, 1.5, 469.167, 469.167},
			want:   [4]float64{1, 6.67, 0.21, 0.43},
		},
		{
			name:   "Tunisian dinar, factor 1000",
			money:  currency.TND,
			values: [4]float64{100, 1000, 10000, 1999.99},
			rates:  [4]float64{1, 1.5, 3.714896, 3.714896},
			want:   [4]float64{0.1, 0.67, 2.69, 0.54},
		},
		{
			name:   "Uruguayan centésimos, factor 10000",
			money:  currency.UYW,
			values: [4]float64{100, 1000, 10000, 1999999.99},
			rates:  [4]float64{1, 1.5, 42.530333, 42.530333},
			want:   [4]float64{0.01, 0.07, 0.02, 4.70},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < len(tt.values); i++ {
				got := accounting.Exchange(tt.money, tt.values[i], tt.rates[i])
				if diff := cmp.Diff(tt.want[i], got); diff != "" {
					t.Errorf("Exchange() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestRatNetAmount(t *testing.T) {
	br := func(f64 float64) *big.Rat {
		return new(big.Rat).SetFloat64(f64)
	}

	type args struct {
		value *big.Rat
		rate  *big.Rat
	}
	var tests = []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "gross 10 vat 19",
			args: args{
				value: br(10),
				rate:  br(.19),
			},
			want:    "8.403361345",
			wantErr: false,
		},
		{
			name: "gross 19.99 vat 20",
			args: args{
				value: br(19.99),
				rate:  br(.20),
			},
			want:    "16.658333333",
			wantErr: false,
		},
		{
			name: "gross 123 vat 7",
			args: args{
				value: br(123),
				rate:  br(.07),
			},
			want:    "114.953271028",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := accounting.RatNetAmount(tt.args.value, tt.args.rate)
			if (err != nil) != tt.wantErr {
				t.Errorf("RatNetAmount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got.Text('f', 9)); diff != "" {
				t.Errorf("RatNetAmount() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestToMinorUnit(t *testing.T) {
	for _, tt := range minorUnitTestsCases {
		t.Run(tt.name, func(t *testing.T) {
			i := accounting.ToMinorUnit(currency.GBP, tt.f64)
			if diff := cmp.Diff(tt.i64, i); diff != "" {
				t.Errorf("ToMinorUnit() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFromMinorUnit(t *testing.T) {
	for _, tt := range minorUnitTestsCases {
		t.Run(tt.name, func(t *testing.T) {
			f := accounting.FromMinorUnit(currency.GBP, tt.i64)
			if diff := cmp.Diff(tt.f64, f); diff != "" {
				t.Errorf("FromMinorUnit() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestValidateFloatIsPrecise(t *testing.T) {
	tests := []struct {
		name      string
		amount    float64
		wantError bool
	}{
		{
			name:      "exact number",
			amount:    10,
			wantError: false,
		},
		{
			name:      "too long",
			amount:    12.123,
			wantError: true,
		},
		{
			name:      "valid amount",
			amount:    12.12,
			wantError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := accounting.ValidateFloatIsPrecise(tt.amount)
			if (err != nil) != tt.wantError {
				t.Fatalf(
					"ValidateFloatIsPrecise: failed on %v",
					tt.amount,
				)
			}
		})
	}
}

func TestValidateManyFloatsArePrecise(t *testing.T) {
	tests := []struct {
		name      string
		amounts   []float64
		wantError bool
	}{
		{
			name:      "exact numbers",
			amounts:   []float64{10, 11, 12},
			wantError: false,
		},
		{
			name:      "one is too long",
			amounts:   []float64{11, 12.123, 13},
			wantError: true,
		},
		{
			name:      "valid amounts",
			amounts:   []float64{11.11, 12.12, 13.13},
			wantError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := accounting.ValidateManyFloatsArePrecise(tt.amounts...)
			if (err != nil) != tt.wantError {
				t.Fatalf("ValidateFloatIsPrecise: failed with %v", err)
			}
		})
	}
}

func TestNetAmount(t *testing.T) {
	type args struct {
		gross int64
		rate  float64
	}
	tests := []struct {
		args    args
		want    int64
		wantErr bool
	}{
		{
			args: args{
				gross: 1999,
				rate:  0.20,
			},
			want:    1666,
			wantErr: false,
		},
		{
			args: args{
				gross: 1234,
				rate:  0.19,
			},
			want:    1037,
			wantErr: false,
		},
		{
			args: args{
				gross: 9999,
				rate:  0.33,
			},
			want:    7518,
			wantErr: false,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("NetAmount_%d", i), func(t *testing.T) {
			got, err := accounting.NetAmount(tt.args.gross, tt.args.rate)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetAmount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("NetAmount() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFuzzNetAmount(t *testing.T) {
	type args struct {
		Gross int64
		Rate  float64
	}

	fz := fuzz.New()

	// 1M iterations.
	for i := 0; i < 1e6; i++ {
		var a args
		fz.Fuzz(&a)
		// we don't really care about the values here.
		// we just want to prove that the function cannot panic.
		netAmount, err := accounting.NetAmount(a.Gross, a.Rate)
		useInt64(netAmount)
		useError(err)
	}
}

func TestFuzzTaxAmount(t *testing.T) {
	// One might call this test spurious.
	// And you'd be right.
	// But if we ever change the internals, we'll be glad to have it!
	type args struct {
		Gross int64
		Net   int64
	}

	fz := fuzz.New()

	// 1M iterations.
	for i := 0; i < 1e6; i++ {
		var a args
		fz.Fuzz(&a)

		// we don't really care about the values here.
		// we just want to prove that the function cannot panic.
		taxAmount, taxAmountError := accounting.TaxAmount(a.Gross, a.Net)
		useInt64(taxAmount)
		useError(taxAmountError)
	}
}

var minorUnitTestsCases = []struct {
	name string
	f64  float64
	i64  int64
}{
	{
		name: "known dangerous value",
		f64:  9.95,
		i64:  995,
	},
	{
		name: "known dangerous negative value",
		f64:  -9.95,
		i64:  -995,
	},
	{
		name: "f64:1.15",
		f64:  1.15,
		i64:  115,
	},
	{
		name: "f64:2.25",
		f64:  2.25,
		i64:  225,
	},
	{
		name: "f64:3.35",
		f64:  3.35,
		i64:  335,
	},
	{
		name: "f64:4.45",
		f64:  4.45,
		i64:  445,
	},
	{
		name: "f64:5.55",
		f64:  5.55,
		i64:  555,
	},
	{
		name: "f64:6.65",
		f64:  6.65,
		i64:  665,
	},
	{
		name: "f64:7.75",
		f64:  7.75,
		i64:  775,
	},
	{
		name: "f64:8.85",
		f64:  8.85,
		i64:  885,
	},
	{
		name: "f64:9.95",
		f64:  9.95,
		i64:  995,
	},
	{
		name: "f64:10.10",
		f64:  10.10,
		i64:  1010,
	},
	{
		name: "f64:11.11",
		f64:  11.11,
		i64:  1111,
	},
	{
		name: "f64:12.12",
		f64:  12.12,
		i64:  1212,
	},
	{
		name: "f64:13.13",
		f64:  13.13,
		i64:  1313,
	},
	{
		name: "f64:14.14",
		f64:  14.14,
		i64:  1414,
	},
	{
		name: "f64:15.15",
		f64:  15.15,
		i64:  1515,
	},
	{
		name: "f64:16.16",
		f64:  16.16,
		i64:  1616,
	},
	{
		name: "f64:17.17",
		f64:  17.17,
		i64:  1717,
	},
	{
		name: "f64:18.18",
		f64:  18.18,
		i64:  1818,
	},
	{
		name: "f64:19.19",
		f64:  19.19,
		i64:  1919,
	},
	{
		name: "f64:20.20",
		f64:  20.20,
		i64:  2020,
	},
	{
		name: "f64:21.21",
		f64:  21.21,
		i64:  2121,
	},
	{
		name: "f64:22.22",
		f64:  22.22,
		i64:  2222,
	},
	{
		name: "f64:23.23",
		f64:  23.23,
		i64:  2323,
	},
	{
		name: "f64:24.24",
		f64:  24.24,
		i64:  2424,
	},
	{
		name: "f64:25.25",
		f64:  25.25,
		i64:  2525,
	},
	{
		name: "f64:26.26",
		f64:  26.26,
		i64:  2626,
	},
	{
		name: "f64:27.27",
		f64:  27.27,
		i64:  2727,
	},
	{
		name: "f64:28.28",
		f64:  28.28,
		i64:  2828,
	},
	{
		name: "f64:29.29",
		f64:  29.29,
		i64:  2929,
	},
	{
		name: "f64:30.30",
		f64:  30.30,
		i64:  3030,
	},
	{
		name: "f64:31.31",
		f64:  31.31,
		i64:  3131,
	},
	{
		name: "f64:32.32",
		f64:  32.32,
		i64:  3232,
	},
	{
		name: "f64:33.33",
		f64:  33.33,
		i64:  3333,
	},
	{
		name: "f64:34.34",
		f64:  34.34,
		i64:  3434,
	},
	{
		name: "f64:35.35",
		f64:  35.35,
		i64:  3535,
	},
	{
		name: "f64:36.36",
		f64:  36.36,
		i64:  3636,
	},
	{
		name: "f64:37.37",
		f64:  37.37,
		i64:  3737,
	},
	{
		name: "f64:38.38",
		f64:  38.38,
		i64:  3838,
	},
	{
		name: "f64:39.39",
		f64:  39.39,
		i64:  3939,
	},
	{
		name: "f64:40.40",
		f64:  40.40,
		i64:  4040,
	},
	{
		name: "f64:41.41",
		f64:  41.41,
		i64:  4141,
	},
	{
		name: "f64:42.42",
		f64:  42.42,
		i64:  4242,
	},
	{
		name: "f64:43.43",
		f64:  43.43,
		i64:  4343,
	},
	{
		name: "f64:44.44",
		f64:  44.44,
		i64:  4444,
	},
	{
		name: "f64:45.45",
		f64:  45.45,
		i64:  4545,
	},
	{
		name: "f64:46.46",
		f64:  46.46,
		i64:  4646,
	},
	{
		name: "f64:47.47",
		f64:  47.47,
		i64:  4747,
	},
	{
		name: "f64:48.48",
		f64:  48.48,
		i64:  4848,
	},
	{
		name: "f64:49.49",
		f64:  49.49,
		i64:  4949,
	},
	{
		name: "f64:50.50",
		f64:  50.50,
		i64:  5050,
	},
	{
		name: "f64:51.51",
		f64:  51.51,
		i64:  5151,
	},
	{
		name: "f64:52.52",
		f64:  52.52,
		i64:  5252,
	},
	{
		name: "f64:53.53",
		f64:  53.53,
		i64:  5353,
	},
	{
		name: "f64:54.54",
		f64:  54.54,
		i64:  5454,
	},
	{
		name: "f64:55.55",
		f64:  55.55,
		i64:  5555,
	},
	{
		name: "f64:56.56",
		f64:  56.56,
		i64:  5656,
	},
	{
		name: "f64:57.57",
		f64:  57.57,
		i64:  5757,
	},
	{
		name: "f64:58.58",
		f64:  58.58,
		i64:  5858,
	},
	{
		name: "f64:59.59",
		f64:  59.59,
		i64:  5959,
	},
	{
		name: "f64:60.60",
		f64:  60.60,
		i64:  6060,
	},
	{
		name: "f64:61.61",
		f64:  61.61,
		i64:  6161,
	},
	{
		name: "f64:62.62",
		f64:  62.62,
		i64:  6262,
	},
	{
		name: "f64:63.63",
		f64:  63.63,
		i64:  6363,
	},
	{
		name: "f64:64.64",
		f64:  64.64,
		i64:  6464,
	},
	{
		name: "f64:65.65",
		f64:  65.65,
		i64:  6565,
	},
	{
		name: "f64:66.66",
		f64:  66.66,
		i64:  6666,
	},
	{
		name: "f64:67.67",
		f64:  67.67,
		i64:  6767,
	},
	{
		name: "f64:68.68",
		f64:  68.68,
		i64:  6868,
	},
	{
		name: "f64:69.69",
		f64:  69.69,
		i64:  6969,
	},
	{
		name: "f64:70.70",
		f64:  70.70,
		i64:  7070,
	},
	{
		name: "f64:71.71",
		f64:  71.71,
		i64:  7171,
	},
	{
		name: "f64:72.72",
		f64:  72.72,
		i64:  7272,
	},
	{
		name: "f64:73.73",
		f64:  73.73,
		i64:  7373,
	},
	{
		name: "f64:74.74",
		f64:  74.74,
		i64:  7474,
	},
	{
		name: "f64:75.75",
		f64:  75.75,
		i64:  7575,
	},
	{
		name: "f64:76.76",
		f64:  76.76,
		i64:  7676,
	},
	{
		name: "f64:77.77",
		f64:  77.77,
		i64:  7777,
	},
	{
		name: "f64:78.78",
		f64:  78.78,
		i64:  7878,
	},
	{
		name: "f64:79.79",
		f64:  79.79,
		i64:  7979,
	},
	{
		name: "f64:80.80",
		f64:  80.80,
		i64:  8080,
	},
	{
		name: "f64:81.81",
		f64:  81.81,
		i64:  8181,
	},
	{
		name: "f64:82.82",
		f64:  82.82,
		i64:  8282,
	},
	{
		name: "f64:83.83",
		f64:  83.83,
		i64:  8383,
	},
	{
		name: "f64:84.84",
		f64:  84.84,
		i64:  8484,
	},
	{
		name: "f64:85.85",
		f64:  85.85,
		i64:  8585,
	},
	{
		name: "f64:86.86",
		f64:  86.86,
		i64:  8686,
	},
	{
		name: "f64:87.87",
		f64:  87.87,
		i64:  8787,
	},
	{
		name: "f64:88.88",
		f64:  88.88,
		i64:  8888,
	},
	{
		name: "f64:89.89",
		f64:  89.89,
		i64:  8989,
	},
	{
		name: "f64:90.90",
		f64:  90.90,
		i64:  9090,
	},
	{
		name: "f64:91.91",
		f64:  91.91,
		i64:  9191,
	},
	{
		name: "f64:92.92",
		f64:  92.92,
		i64:  9292,
	},
	{
		name: "f64:93.93",
		f64:  93.93,
		i64:  9393,
	},
	{
		name: "f64:94.94",
		f64:  94.94,
		i64:  9494,
	},
	{
		name: "f64:95.95",
		f64:  95.95,
		i64:  9595,
	},
	{
		name: "f64:96.96",
		f64:  96.96,
		i64:  9696,
	},
	{
		name: "f64:97.97",
		f64:  97.97,
		i64:  9797,
	},
	{
		name: "f64:98.98",
		f64:  98.98,
		i64:  9898,
	},
	{
		name: "f64:99.99",
		f64:  99.99,
		i64:  9999,
	},
	{
		name: "f64:100.105",
		f64:  100.10,
		i64:  10010,
	},
	{
		name: "f64:1.15",
		f64:  -1.15,
		i64:  -115,
	},
	{
		name: "f64:2.25",
		f64:  -2.25,
		i64:  -225,
	},
	{
		name: "f64:3.35",
		f64:  -3.35,
		i64:  -335,
	},
	{
		name: "f64:4.45",
		f64:  -4.45,
		i64:  -445,
	},
	{
		name: "f64:5.55",
		f64:  -5.55,
		i64:  -555,
	},
	{
		name: "f64:6.65",
		f64:  -6.65,
		i64:  -665,
	},
	{
		name: "f64:7.75",
		f64:  -7.75,
		i64:  -775,
	},
	{
		name: "f64:8.85",
		f64:  -8.85,
		i64:  -885,
	},
	{
		name: "f64:9.95",
		f64:  -9.95,
		i64:  -995,
	},
	{
		name: "f64:10.10",
		f64:  -10.10,
		i64:  -1010,
	},
	{
		name: "f64:11.11",
		f64:  -11.11,
		i64:  -1111,
	},
	{
		name: "f64:12.12",
		f64:  -12.12,
		i64:  -1212,
	},
	{
		name: "f64:13.13",
		f64:  -13.13,
		i64:  -1313,
	},
	{
		name: "f64:14.14",
		f64:  -14.14,
		i64:  -1414,
	},
	{
		name: "f64:15.15",
		f64:  -15.15,
		i64:  -1515,
	},
	{
		name: "f64:16.16",
		f64:  -16.16,
		i64:  -1616,
	},
	{
		name: "f64:17.17",
		f64:  -17.17,
		i64:  -1717,
	},
	{
		name: "f64:18.18",
		f64:  -18.18,
		i64:  -1818,
	},
	{
		name: "f64:19.19",
		f64:  -19.19,
		i64:  -1919,
	},
	{
		name: "f64:20.20",
		f64:  -20.20,
		i64:  -2020,
	},
	{
		name: "f64:21.21",
		f64:  -21.21,
		i64:  -2121,
	},
	{
		name: "f64:22.22",
		f64:  -22.22,
		i64:  -2222,
	},
	{
		name: "f64:23.23",
		f64:  -23.23,
		i64:  -2323,
	},
	{
		name: "f64:24.24",
		f64:  -24.24,
		i64:  -2424,
	},
	{
		name: "f64:25.25",
		f64:  -25.25,
		i64:  -2525,
	},
	{
		name: "f64:26.26",
		f64:  -26.26,
		i64:  -2626,
	},
	{
		name: "f64:27.27",
		f64:  -27.27,
		i64:  -2727,
	},
	{
		name: "f64:28.28",
		f64:  -28.28,
		i64:  -2828,
	},
	{
		name: "f64:29.29",
		f64:  -29.29,
		i64:  -2929,
	},
	{
		name: "f64:30.30",
		f64:  -30.30,
		i64:  -3030,
	},
	{
		name: "f64:31.31",
		f64:  -31.31,
		i64:  -3131,
	},
	{
		name: "f64:32.32",
		f64:  -32.32,
		i64:  -3232,
	},
	{
		name: "f64:33.33",
		f64:  -33.33,
		i64:  -3333,
	},
	{
		name: "f64:34.34",
		f64:  -34.34,
		i64:  -3434,
	},
	{
		name: "f64:35.35",
		f64:  -35.35,
		i64:  -3535,
	},
	{
		name: "f64:36.36",
		f64:  -36.36,
		i64:  -3636,
	},
	{
		name: "f64:37.37",
		f64:  -37.37,
		i64:  -3737,
	},
	{
		name: "f64:38.38",
		f64:  -38.38,
		i64:  -3838,
	},
	{
		name: "f64:39.39",
		f64:  -39.39,
		i64:  -3939,
	},
	{
		name: "f64:40.40",
		f64:  -40.40,
		i64:  -4040,
	},
	{
		name: "f64:41.41",
		f64:  -41.41,
		i64:  -4141,
	},
	{
		name: "f64:42.42",
		f64:  -42.42,
		i64:  -4242,
	},
	{
		name: "f64:43.43",
		f64:  -43.43,
		i64:  -4343,
	},
	{
		name: "f64:44.44",
		f64:  -44.44,
		i64:  -4444,
	},
	{
		name: "f64:45.45",
		f64:  -45.45,
		i64:  -4545,
	},
	{
		name: "f64:46.46",
		f64:  -46.46,
		i64:  -4646,
	},
	{
		name: "f64:47.47",
		f64:  -47.47,
		i64:  -4747,
	},
	{
		name: "f64:48.48",
		f64:  -48.48,
		i64:  -4848,
	},
	{
		name: "f64:49.49",
		f64:  -49.49,
		i64:  -4949,
	},
	{
		name: "f64:50.50",
		f64:  -50.50,
		i64:  -5050,
	},
	{
		name: "f64:51.51",
		f64:  -51.51,
		i64:  -5151,
	},
	{
		name: "f64:52.52",
		f64:  -52.52,
		i64:  -5252,
	},
	{
		name: "f64:53.53",
		f64:  -53.53,
		i64:  -5353,
	},
	{
		name: "f64:54.54",
		f64:  -54.54,
		i64:  -5454,
	},
	{
		name: "f64:55.55",
		f64:  -55.55,
		i64:  -5555,
	},
	{
		name: "f64:56.56",
		f64:  -56.56,
		i64:  -5656,
	},
	{
		name: "f64:57.57",
		f64:  -57.57,
		i64:  -5757,
	},
	{
		name: "f64:58.58",
		f64:  -58.58,
		i64:  -5858,
	},
	{
		name: "f64:59.59",
		f64:  -59.59,
		i64:  -5959,
	},
	{
		name: "f64:60.60",
		f64:  -60.60,
		i64:  -6060,
	},
	{
		name: "f64:61.61",
		f64:  -61.61,
		i64:  -6161,
	},
	{
		name: "f64:62.62",
		f64:  -62.62,
		i64:  -6262,
	},
	{
		name: "f64:63.63",
		f64:  -63.63,
		i64:  -6363,
	},
	{
		name: "f64:64.64",
		f64:  -64.64,
		i64:  -6464,
	},
	{
		name: "f64:65.65",
		f64:  -65.65,
		i64:  -6565,
	},
	{
		name: "f64:66.66",
		f64:  -66.66,
		i64:  -6666,
	},
	{
		name: "f64:67.67",
		f64:  -67.67,
		i64:  -6767,
	},
	{
		name: "f64:68.68",
		f64:  -68.68,
		i64:  -6868,
	},
	{
		name: "f64:69.69",
		f64:  -69.69,
		i64:  -6969,
	},
	{
		name: "f64:70.70",
		f64:  -70.70,
		i64:  -7070,
	},
	{
		name: "f64:71.71",
		f64:  -71.71,
		i64:  -7171,
	},
	{
		name: "f64:72.72",
		f64:  -72.72,
		i64:  -7272,
	},
	{
		name: "f64:73.73",
		f64:  -73.73,
		i64:  -7373,
	},
	{
		name: "f64:74.74",
		f64:  -74.74,
		i64:  -7474,
	},
	{
		name: "f64:75.75",
		f64:  -75.75,
		i64:  -7575,
	},
	{
		name: "f64:76.76",
		f64:  -76.76,
		i64:  -7676,
	},
	{
		name: "f64:77.77",
		f64:  -77.77,
		i64:  -7777,
	},
	{
		name: "f64:78.78",
		f64:  -78.78,
		i64:  -7878,
	},
	{
		name: "f64:79.79",
		f64:  -79.79,
		i64:  -7979,
	},
	{
		name: "f64:80.80",
		f64:  -80.80,
		i64:  -8080,
	},
	{
		name: "f64:81.81",
		f64:  -81.81,
		i64:  -8181,
	},
	{
		name: "f64:82.82",
		f64:  -82.82,
		i64:  -8282,
	},
	{
		name: "f64:83.83",
		f64:  -83.83,
		i64:  -8383,
	},
	{
		name: "f64:84.84",
		f64:  -84.84,
		i64:  -8484,
	},
	{
		name: "f64:85.85",
		f64:  -85.85,
		i64:  -8585,
	},
	{
		name: "f64:86.86",
		f64:  -86.86,
		i64:  -8686,
	},
	{
		name: "f64:87.87",
		f64:  -87.87,
		i64:  -8787,
	},
	{
		name: "f64:88.88",
		f64:  -88.88,
		i64:  -8888,
	},
	{
		name: "f64:89.89",
		f64:  -89.89,
		i64:  -8989,
	},
	{
		name: "f64:90.90",
		f64:  -90.90,
		i64:  -9090,
	},
	{
		name: "f64:91.91",
		f64:  -91.91,
		i64:  -9191,
	},
	{
		name: "f64:92.92",
		f64:  -92.92,
		i64:  -9292,
	},
	{
		name: "f64:93.93",
		f64:  -93.93,
		i64:  -9393,
	},
	{
		name: "f64:94.94",
		f64:  -94.94,
		i64:  -9494,
	},
	{
		name: "f64:95.95",
		f64:  -95.95,
		i64:  -9595,
	},
	{
		name: "f64:96.96",
		f64:  -96.96,
		i64:  -9696,
	},
	{
		name: "f64:97.97",
		f64:  -97.97,
		i64:  -9797,
	},
	{
		name: "f64:98.98",
		f64:  -98.98,
		i64:  -9898,
	},
	{
		name: "f64:99.99",
		f64:  -99.99,
		i64:  -9999,
	},
	{
		name: "f64:100.105",
		f64:  -100.10,
		i64:  -10010,
	},
}

//go:noinline
func useInt64(i int64) {}

//go:noinline
func useFloat64(f float64) {}

//go:noinline
func useError(e error) {}
