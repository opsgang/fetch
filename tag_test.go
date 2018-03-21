package main

import (
	"testing"
)

func TestTagToGet(t *testing.T) {
	t.Parallel()
	// specific - should work regardless of meeting semver format
	// specific - includes when val prefixed with =
	// if no tag provided - should provide latest
	// if constraint can not be met, err

	tags := []string{"foo", "v1.0.0", "10.0.20-some-txt", "bar", "10.0.10", "10.1.9", "100.10.1"}

	o := fetchOpts{}

	cases := []struct {
		constraint string
		tag        string
		err        bool
	}{
		{"~>10.0", "10.1.9", false},
		{"bar", "bar", false},
		{"=bar", "bar", false},
		{"", "100.10.1", false},
		{"<1.0.0", "", true},
	}

	for _, tc := range cases {
		o.tagConstraint = tc.constraint

		if tag, err := o.tagToGet(tags); err != nil {
			if tc.err != true {
				t.Fatalf("Did not expect the error we received: %s\nconstraint:[%s]", err, tc.constraint)
			}
		} else {
			if tag != tc.tag {
				t.Fatalf("Did not expect the tag we received: %s\nconstraint:[%s]", tag, tc.constraint)
			}
		}

	}
}

func TestBestFitTag(t *testing.T) {
	t.Parallel()

	cases := []struct {
		tagConstraint string
		tags          []string
		expectedTag   string
	}{
		{"1.0.7", []string{"1.0.7"}, "1.0.7"},

		{"~> 1.0.0", []string{"1.0.5", "1.0.6", "1.0.7", "1.0.8", "1.0.9", "1.1.0", "1.2.3"}, "1.0.9"},
		{"~> 1.0.7", []string{"1.0.5", "1.0.6", "1.0.7", "1.0.8", "1.0.9", "1.1.0", "1.2.3"}, "1.0.9"},
		{"~> 1.1.0", []string{"1.0.5", "1.0.6", "1.0.7", "1.0.8", "1.0.9", "1.1.0", "1.2.3"}, "1.1.0"},
		{"~> 1.1.1", []string{"1.0.5", "1.0.6", "1.0.7", "1.0.8", "1.0.9", "1.1.0", "1.1.1", "1.1.2", "1.1.3", "1.2.3", "1.4.0", "2.0.0", "2.1.0"}, "1.1.3"},
		{"~> 1.2.1", []string{"1.0.5", "1.0.6", "1.0.7", "1.0.8", "1.0.9", "1.1.0", "1.1.1", "1.1.2", "1.1.3", "1.2.3", "1.4.0", "2.0.0", "2.1.0"}, "1.2.3"},
		{"~> 1.1", []string{"1.0.5", "1.0.6", "1.0.7", "1.0.8", "1.0.9", "1.1.0", "1.1.1", "1.1.2", "1.1.3", "1.2.3", "1.4.0", "2.0.0", "2.1.0"}, "1.4.0"},
		{"~> 1.2", []string{"1.0.5", "1.0.6", "1.0.7", "1.0.8", "1.0.9", "1.1.0", "1.1.1", "1.1.2", "1.1.3", "1.2.3", "1.4.0", "2.0.0", "2.1.0"}, "1.4.0"},
		{"~> 1.3", []string{"1.0.5", "1.0.6", "1.0.7", "1.0.8", "1.0.9", "1.1.0", "1.1.1", "1.1.2", "1.1.3", "1.2.3", "1.4.0", "2.0.0", "2.1.0"}, "1.4.0"},

		{">= 1.3", []string{"1.0.5", "1.0.6", "1.0.7", "1.0.8", "1.0.9", "1.1.0", "1.1.1", "1.1.2", "1.1.3", "1.2.3", "1.4.0", "2.0.0", "2.1.0"}, "2.1.0"},

		{"v1.0.7", []string{"v1.0.7"}, "v1.0.7"},
		{"v1.0.7", []string{}, ""},
	}

	for _, tc := range cases {
		tag, err := bestFitTag(tc.tagConstraint, tc.tags)
		if err != nil {
			t.Fatalf("Failed on call to bestFitTag: %s", err)
		}

		if tag != tc.expectedTag {
			t.Fatalf("Given constraint %s and tag list %v, expected %s, but received: %s", tc.tagConstraint, tc.tags, tc.expectedTag, tag)
		}
	}
}

func TestIsTagConstraintOrExactTag(t *testing.T) {
	t.Parallel()

	cases := []struct {
		tagConstraint string
		desiredTag    string
		exact         bool
	}{
		{"1.0.7", "1.0.7", true},
		{" 1.0.7	 ", "1.0.7", true},
		{"=1.0.7", "1.0.7", true},
		{"= 1.0.7", "1.0.7", true},

		{"~> 1.0.0", "~> 1.0.0", false},
		{"> 1.3", "> 1.3", false},
		{">= 1.3", ">= 1.3", false},
		{"<= 1.3", "<= 1.3", false},
		{"< 1.3", "< 1.3", false},

		{"v1.0.7", "v1.0.7", true},
		{" v1.0.7	 ", "v1.0.7", true},
		{"=v1.0.7", "v1.0.7", true},
		{"= v1.0.7", "v1.0.7", true},
	}

	for _, tc := range cases {
		exact, desiredTag := isTagConstraintOrExactTag(tc.tagConstraint)
		if exact != tc.exact {
			t.Fatalf("Given constraint: \"%s\", expected %t, but received %t", tc.tagConstraint, tc.exact, exact)
		}
		if desiredTag != tc.desiredTag {
			t.Fatalf("Given constraint: \"%s\", expected result tag: \"%s\", but received \"%s\"", tc.tagConstraint, tc.desiredTag, desiredTag)
		}
	}
}

func TestBestFitTagOnEmptyConstraint(t *testing.T) {
	t.Parallel()

	cases := []struct {
		tagConstraint string
		tags          []string
		expectedTag   string
	}{
		{"", []string{"v0.0.1", "v0.0.2", "v0.0.3"}, "v0.0.3"},
		{"", []string{"1.0.5", "1.0.6", "1.0.7", "1.0.8", "1.0.9", "1.1.0", "1.2.3"}, "1.2.3"},
		{"", []string{"1.0.5", "1.0.6", "1.0.7", "1.0.8", "1.0.9", "1.1.0", "1.1.1", "1.1.2", "1.1.3", "1.2.3", "1.4.0", "2.0.0", "2.1.0"}, "2.1.0"},
		{"", []string{}, ""},
	}

	for _, tc := range cases {
		tag, err := bestFitTag(tc.tagConstraint, tc.tags)
		if err != nil {
			t.Fatalf("Failed on call to bestFitTag: %s", err)
		}

		if tag != tc.expectedTag {
			t.Fatalf("Given constraint %s and tag list %v, expected %s, but received: %s", tc.tagConstraint, tc.tags, tc.expectedTag, tag)
		}
	}
}

func TestBestFitTagOnMalformedConstraint(t *testing.T) {
	t.Parallel()

	cases := []struct {
		tagConstraint string
	}{
		{"josh"},
		{"plump elephants dancing in the night"},
	}

	for _, tc := range cases {
		_, err := bestFitTag(tc.tagConstraint, []string{"v0.0.1"})
		if err == nil {
			t.Fatalf("Expected malformed constraint error, but received nothing.")
		}
	}
}

func TestBestFitTagNoSuchTag(t *testing.T) {
	t.Parallel()
	cases := []struct {
		tagConstraint string
		tags          []string
	}{
		{"~> 0.0.4", []string{"0.0.1", "0.0.2", "0.0.3"}},
		{"> 0.0.4", []string{"0.0.1", "0.0.2", "0.0.3"}},
	}

	for _, tc := range cases {
		_, err := bestFitTag(tc.tagConstraint, tc.tags)
		if err == nil {
			t.Fatalf("Expected 'Tag does not exist' but received nothing")
		}
	}
}
