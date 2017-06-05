package casper

import "testing"

func TestGenerateCookie(t *testing.T) {
	cases := []struct {
		assets      []string
		P           int
		cookieValue string
	}{
		{
			[]string{
				"/static/example.js",
			},
			1 << 6,
			"JA",
		},

		{
			[]string{
				"/js/jquery-1.9.1.min.js",
				"/assets/style.css",
			},
			1 << 6,
			"gU4",
		},

		{
			[]string{
				"/js/jquery-1.9.1.min.js",
				"/assets/style.css",
				"/static/logo.jpg",
				"/static/cover.jpg",
			},
			1 << 6,
			"gU54MA",
		},

		{
			[]string{
				"/js/jquery-1.9.1.min.js",
				"/assets/style.css",
				"/static/logo.jpg",
				"/static/cover.jpg",
			},
			1 << 10,
			"MMOJEkWo",
		},

		// See how long cookie is when push many files.
		// Minimum number of bits is N*log(P) = 20 * log(1<<6) = 120 bits = 15bytes
		{
			[]string{
				"/static/example1.jpg",
				"/static/example2.jpg",
				"/static/example3.jpg",
				"/static/example4.jpg",
				"/static/example5.jpg",
				"/static/example6.jpg",
				"/static/example7.jpg",
				"/static/example8.jpg",
				"/static/example9.jpg",
				"/static/example10.jpg",
				"/static/example11.jpg",
				"/static/example12.jpg",
				"/static/example13.jpg",
				"/static/example14.jpg",
				"/static/example15.jpg",
				"/static/example16.jpg",
				"/static/example17.jpg",
				"/static/example18.jpg",
				"/static/example19.jpg",
				"/static/example20.jpg",
			},
			1 << 6,
			"FmDhUxQHeuwQYINoQrxmr1g_iw", // 26bytes
		},
	}

	for _, tc := range cases {
		casper := &Casper{
			p: uint(tc.P),
			n: uint(len(tc.assets)),
		}

		hashValues := make([]uint, 0, len(tc.assets))
		for _, content := range tc.assets {
			hashValues = append(hashValues, casper.hash([]byte(content)))
		}

		cookie, err := casper.generateCookie(hashValues)
		if err != nil {
			t.Fatalf("generateCookie should not fail")
		}

		if got, want := cookie.Value, tc.cookieValue; got != want {
			t.Fatalf("generateCookie=%q, want=%q", got, want)
		}
	}
}

var benchCase = struct {
	assets      []string
	cookieValue string
}{
	// See how long cookie is when push many files.
	// Minimum number of bits is N*log(P) = 20 * log(1<<6) = 120 bits = 15bytes
	[]string{
		"/static/example1.jpg",
		"/static/example2.jpg",
		"/static/example3.jpg",
		"/static/example4.jpg",
		"/static/example5.jpg",
		"/static/example6.jpg",
		"/static/example7.jpg",
		"/static/example8.jpg",
		"/static/example9.jpg",
		"/static/example10.jpg",
		"/static/example11.jpg",
		"/static/example12.jpg",
		"/static/example13.jpg",
		"/static/example14.jpg",
		"/static/example15.jpg",
		"/static/example16.jpg",
		"/static/example17.jpg",
		"/static/example18.jpg",
		"/static/example19.jpg",
		"/static/example20.jpg",
	},

	"yKhHjfQdD63uyqmI4_ducgNojOGO_8QiuzPZxkHzPQqLsR82H_h7wA",
}

var benchCasper = &Casper{
	p: uint(64 * 64),
	n: uint(64),
}

func BenchmarkGenerateCookie(b *testing.B) {

	hashValues := make([]uint, 0, len(benchCase.assets))

	for n := 0; n < b.N; n++ {
		for _, content := range benchCase.assets {
			hashValues = append(hashValues, benchCasper.hash([]byte(content)))
		}

		cookie, err := benchCasper.generateCookie(hashValues)
		if err != nil {
			b.Fatal("generateCookie should not fail")
		}

		if got, want := cookie.Value, benchCase.cookieValue; got != want {
			b.Fatalf("generateCookie=%q, want=%q", got, want)
		}

		hashValues = hashValues[:0]
	}
}
