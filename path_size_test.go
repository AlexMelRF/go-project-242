package code

import (
	"testing"
)

type testCase struct {
	path		string
	output		int64
	recursive	bool
	all			bool
}

var testDataArray = []testCase {
	{"./testdata/dir1", 1599, true, false},
	{"./testdata/dir1/dir2", 1599, true, false},
	{"./testdata/dir1", 3997, true, true},
	{"./testdata/dir1/dir2", 1599, false, true},
	{"./testdata", 2399, false, true},
}

func TestGetSize(t *testing.T) {
	for _, td := range testDataArray {
		t.Run(td.path, func(t *testing.T) {
			result, err := GetSize(td.path, td.recursive, td.all)
			if err != nil {
				t.Errorf("GetSize(%s, %v, %v) returned error: %v", 
					td.path, td.recursive, td.all, err)
				return
			}

			if result != td.output {
				t.Errorf("GetSize(%s, %v, %v) = %d, expected %d", 
					td.path, td.recursive, td.all, result, td.output)
			}
		})
	}	
}

func TestFormatSize(t *testing.T) {
	testCases := []struct {
		input    int64
		expected string
	}{
		{0, "0B"},
		{1023, "1023B"},
		{1024, "1.0KB"},
		{1234567, "1.2MB"},
		{1024 * 1024 * 1024, "1.0GB"},
		{1024 * 1024 * 1024 * 1024, "1.0TB"},
		{1024 * 1024 * 1024 * 1024 * 1024, "1.0PB"},
		{1024 * 1024 * 1024 * 1024 * 1024 * 1024, "1.0EB"},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			result := FormatSize(tc.input, true)
			if result != tc.expected {
				t.Errorf("FormatSize(%d) = %s, expected %s", tc.input, result, tc.expected)
			}
		})
	}
}
