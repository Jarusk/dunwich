package dunwich

import (
	"testing"
)

func BenchmarkBuildCorpus(b *testing.B) {
	for i := 0; i < b.N; i++ {
		target := []string{}
		buildCorpus(&target)
	}
}

func TestIsChapterMarker(t *testing.T) {
	cases := map[string]struct {
		input    string
		expected bool
	}{
		"shouldReturnTrueWithSingleDigit": {
			input:    "1",
			expected: true,
		},
		"shouldReturnTrueWithMultiDigit": {
			input:    "11",
			expected: true,
		},
		"shouldReturnTrueWithSingleDigitAndNewline": {
			input:    "5\n",
			expected: true,
		},
		"shouldReturnFalseWithEmptyLine": {
			input:    "",
			expected: false,
		},
		"shouldReturnFalseWithEmptyLineAndNewline": {
			input:    "\n",
			expected: false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			res := isChapterMarker(tc.input)
			if res != tc.expected {
				t.Errorf("exepcted %v, got %v", tc.expected, res)
			}
		})
	}
}

func TestGetNumSegements(t *testing.T) {
	cases := map[string]struct {
		expected int
	}{
		"shouldReturnConfirmedSegmentCount": {
			expected: 149,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			res := GetNumSegments()
			if res != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, res)
			}
		})
	}
}

func TestGetSegment(t *testing.T) {
	cases := map[string]struct {
		id          int
		expected    string
		errExpected bool
	}{
		"shouldMatchSegment0": {
			id:          0,
			errExpected: false,
			expected: `When a traveler in north central Massachusetts takes the wrong fork
at the junction of the Aylesbury pike just beyond Dean's Corners he
comes upon a lonely and curious country. The ground gets higher, and
the brier-bordered stone walls press closer and closer against the ruts
of the dusty, curving road. The trees of the frequent forest belts
seem too large, and the wild weeds, brambles, and grasses attain a
luxuriance not often found in settled regions. At the same time the
planted fields appear singularly few and barren; while the sparsely
scattered houses wear a surprizing uniform aspect of age, squalor, and
dilapidation. Without knowing why, one hesitates to ask directions
from the gnarled, solitary figures spied now and then on crumbling
doorsteps or in the sloping, rock-strewn meadows. Those figures are
so silent and furtive that one feels somehow confronted by forbidden
things, with which it would be better to have nothing to do. When a
rise in the road brings the mountains in view above the deep woods,
the feeling of strange uneasiness is increased. The summits are too
rounded and symmetrical to give a sense of comfort and naturalness, and
sometimes the sky silhouettes with especial clearness the queer circles
of tall stone pillars with which most of them are crowned.
`,
		},
		"shouldMatchSegment148": {
			id:          148,
			errExpected: false,
			expected: `"But as to this thing we've just sent back--the Whateleys raised it for
a terrible part in the doings that were to come. It grew fast and big
from the same reason that Wilbur grew fast and big--but it beat him
because it had a greater share of the _outsideness_ in it. You needn't
ask how Wilbur called it out of the air. He didn't call it out. _It was
his twin brother, but it looked more like the father than he did._"
`,
		},
		"shouldReturnErrOnOutOfIndex": {
			id:          123456,
			errExpected: true,
			expected:    "",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			res, err := GetSegment(tc.id)
			if err != nil {
				if tc.errExpected {
					// Err was expected, quit test case
					return
				}
				t.Fatalf("received error fetching segment %d", tc.id)
			}
			if *res != tc.expected {
				t.Errorf("expected: \n%s\nreceived: %s\n", tc.expected, *res)
			}
		})
	}
}
