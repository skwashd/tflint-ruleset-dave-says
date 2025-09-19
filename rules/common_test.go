package rules

import (
	"testing"
)

func TestSplitWordsOnUnderscore(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"user_data", []string{"user", "data"}},
		{"my_test_name", []string{"my", "test", "name"}},
		{"single", []string{"single"}},
		{"", []string{""}},
	}

	for _, test := range tests {
		result := SplitWordsOnUnderscore(test.input)
		if len(result) != len(test.expected) {
			t.Errorf("SplitWordsOnUnderscore(%q): expected %d words, got %d", test.input, len(test.expected), len(result))
			continue
		}
		for i, word := range result {
			if word != test.expected[i] {
				t.Errorf("SplitWordsOnUnderscore(%q): expected word %d to be %q, got %q", test.input, i, test.expected[i], word)
			}
		}
	}
}

func TestSplitWordsOnDash(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"user-data", []string{"user", "data"}},
		{"my-test-name", []string{"my", "test", "name"}},
		{"single", []string{"single"}},
		{"", []string{""}},
	}

	for _, test := range tests {
		result := SplitWordsOnDash(test.input)
		if len(result) != len(test.expected) {
			t.Errorf("SplitWordsOnDash(%q): expected %d words, got %d", test.input, len(test.expected), len(result))
			continue
		}
		for i, word := range result {
			if word != test.expected[i] {
				t.Errorf("SplitWordsOnDash(%q): expected word %d to be %q, got %q", test.input, i, test.expected[i], word)
			}
		}
	}
}

func TestSplitWords(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"aws_s3_bucket", []string{"aws", "s3", "bucket"}},
		{"my-test-name", []string{"my", "test", "name"}},
		{"mixed_and-separated", []string{"mixed", "and", "separated"}},
		{"single", []string{"single"}},
		{"", []string{}},
	}

	for _, test := range tests {
		result := SplitWords(test.input)
		if len(result) != len(test.expected) {
			t.Errorf("SplitWords(%q): expected %d words, got %d", test.input, len(test.expected), len(result))
			continue
		}
		for i, word := range result {
			if word != test.expected[i] {
				t.Errorf("SplitWords(%q): expected word %d to be %q, got %q", test.input, i, test.expected[i], word)
			}
		}
	}
}

func TestContainsAnyWord(t *testing.T) {
	tests := []struct {
		haystack     []string
		needle       []string
		expectFound  bool
		expectedWord string
	}{
		{[]string{"user", "bucket"}, []string{"aws", "s3", "bucket"}, true, "bucket"},
		{[]string{"admin", "role"}, []string{"aws", "iam", "role"}, true, "role"},
		{[]string{"user", "data"}, []string{"aws", "s3", "bucket"}, false, ""},
		{[]string{"S3", "storage"}, []string{"aws", "s3", "bucket"}, true, "s3"}, // case insensitive
		{[]string{}, []string{"aws", "s3"}, false, ""},
		{[]string{"test"}, []string{}, false, ""},
	}

	for _, test := range tests {
		found, word := ContainsAnyWord(test.haystack, test.needle)
		if found != test.expectFound {
			t.Errorf("ContainsAnyWord(%v, %v): expected found=%v, got %v", test.haystack, test.needle, test.expectFound, found)
		}
		if word != test.expectedWord {
			t.Errorf("ContainsAnyWord(%v, %v): expected word=%q, got %q", test.haystack, test.needle, test.expectedWord, word)
		}
	}
}

func TestSnakeRegex(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"my_bucket", true},
		{"bucket123", true},
		{"test_123_name", true},
		{"myBucket", false},
		{"my-bucket", false},
		{"My_Bucket", false},
		{"bucket!", false},
		{"", false}, // empty string doesn't match
	}

	for _, test := range tests {
		result := SnakeRegex.MatchString(test.input)
		if result != test.expected {
			t.Errorf("SnakeRegex.MatchString(%q): expected %v, got %v", test.input, test.expected, result)
		}
	}
}

func TestKebabRegex(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"my-bucket", true},
		{"bucket123", true},
		{"test-123-name", true},
		{"myBucket", false},
		{"my_bucket", false},
		{"My-Bucket", false},
		{"bucket!", false},
		{"", false}, // empty string doesn't match
	}

	for _, test := range tests {
		result := KebabRegex.MatchString(test.input)
		if result != test.expected {
			t.Errorf("KebabRegex.MatchString(%q): expected %v, got %v", test.input, test.expected, result)
		}
	}
}
