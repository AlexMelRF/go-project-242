package code

import (
	// "require"
	"testing"
)

type testData struct {
	path		string
	output		int64
	recursive	bool
	all			bool
}

var testDataArray = []testData {
	{"./testdata/dir1", 1599, true, false},
	{"./testdata/dir1/dir2", 1599, true, false},
	{"./testdata/dir1", 1998, true, true},
	{"./testdata/dir1/dir2", 0, false, true},
	{"./testdata", 0, false, true},
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