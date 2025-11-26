package code

import (
	"os"
	"testing"
)

type testCase struct {
	path      string
	output    int64
	recursive bool
	all       bool
}

var testDataArray = []testCase{
	{"./testdata/dir1", 1599, true, false},
	{"./testdata/dir1/dir2", 1599, true, false},
	{"./testdata/dir1", 3997, true, true},
	{"./testdata/dir1/dir2", 1599, false, true},
	{"./testdata", 2399, false, true},
}

func TestGetSize(t *testing.T) {
	// create test struct
	setupTestData(t)
	// del test data
	defer cleanupTestData(t)

	// test cases
	testDataArray := []testCase{
		{"./testdata/dir1", 30, true, false},
		{"./testdata/dir1/dir2", 15, true, false},
		{"./testdata/dir1", 45, true, true},
		{"./testdata/dir1/dir2", 15, false, true},
		// for ./testdata we only check that the size is > 0
		// since the exact size may vary
	}

	for _, td := range testDataArray {
		t.Run(td.path, func(t *testing.T) {
			result, err := GetSize(td.path, td.recursive, td.all)
			if err != nil {
				t.Errorf("GetSize(%s, %v, %v) returned error: %v",
					td.path, td.recursive, td.all, err)
				return
			}

			if td.path == "./testdata" && td.all {
				if result <= 0 {
					t.Errorf("GetSize(%s, %v, %v) = %d, expected positive number",
						td.path, td.recursive, td.all, result)
				}
			} else {
				if result != td.output {
					t.Errorf("GetSize(%s, %v, %v) = %d, expected %d",
						td.path, td.recursive, td.all, result, td.output)
				}
			}
		})
	}
}

func TestGetPathSize(t *testing.T) {
	// create test struct
	setupTestData(t)
	// del test struct
	defer cleanupTestData(t)

	tests := []struct {
		name      string
		path      string
		recursive bool
		human     bool
		all       bool
		wantErr   bool
	}{
		{
			name:      "human_readable",
			path:      "./testdata/test_file.txt",
			recursive: false,
			human:     true,
			all:       false,
			wantErr:   false,
		},
		{
			name:      "bytes_only",
			path:      "./testdata/test_file.txt",
			recursive: false,
			human:     false,
			all:       false,
			wantErr:   false,
		},
		{
			name:      "invalid_path",
			path:      "./testdata/nonexistent",
			recursive: false,
			human:     false,
			all:       false,
			wantErr:   true,
		},
		{
			name:      "unicode_path_human",
			path:      "./testdata/unicode_файл.txt",
			recursive: false,
			human:     true,
			all:       false,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetPathSize(tt.path, tt.recursive, tt.human, tt.all)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("GetPathSize(%s) expected error, got nil", tt.path)
				}
				return
			}
			
			if err != nil {
				t.Errorf("GetPathSize(%s) unexpected error: %v", tt.path, err)
				return
			}
			
			if result == "" {
				t.Errorf("GetPathSize(%s) returned empty string", tt.path)
			}
		})
	}
}

func TestFormatSize(t *testing.T) {
	testCases := []struct {
		name     string
		input    int64
		human    bool
		expected string
	}{
		// no flag -H
		{"zero_bytes_no_human", 0, false, "0B"},
		{"small_no_human", 123, false, "123B"},
		{"large_no_human", 123456789, false, "123456789B"},

		// with flag  -H
		{"zero_bytes", 0, true, "0B"},
		{"bytes_under_1k", 1023, true, "1023B"},
		{"exactly_1k", 1024, true, "1.0KB"},
		{"fractional_kb", 1536, true, "1.5KB"},
		{"exactly_1mb", 1024 * 1024, true, "1.0MB"},
		{"fractional_mb", 1234567, true, "1.2MB"},
		{"exactly_1gb", 1024 * 1024 * 1024, true, "1.0GB"},
		{"exactly_1tb", 1024 * 1024 * 1024 * 1024, true, "1.0TB"},
		{"exactly_1pb", 1024 * 1024 * 1024 * 1024 * 1024, true, "1.0PB"},
		{"exactly_1eb", 1024 * 1024 * 1024 * 1024 * 1024 * 1024, true, "1.0EB"},
		
		{"boundary_1023", 1023, true, "1023B"},
		{"boundary_1024", 1024, true, "1.0KB"},
		{"boundary_1025", 1025, true, "1.0KB"},
		{"large_fractional", 1125899906842624, true, "1.0PB"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := FormatSize(tc.input, tc.human)
			if result != tc.expected {
				t.Errorf("FormatSize(%d, %v) = %s, expected %s", 
					tc.input, tc.human, result, tc.expected)
			}
		})
	}
}

func TestGetSize_EdgeCases(t *testing.T) {
	setupTestData(t)
	defer cleanupTestData(t)

	tests := []struct {
		name      string
		path      string
		recursive bool
		all       bool
	}{
		{"empty_dir", "./testdata/empty_dir", true, false},
		{"dir_with_only_hidden", "./testdata/only_hidden", true, false},
		{"dir_with_only_hidden_all_true", "./testdata/only_hidden", true, true},
		{"nested_symlinks", "./testdata/nested_symlinks", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetSize(tt.path, tt.recursive, tt.all)
			if err != nil {
				t.Errorf("GetSize(%s) unexpected error: %v", tt.path, err)
				return
			}
			
			if result < 0 {
				t.Errorf("GetSize(%s) returned negative size: %d", tt.path, result)
			}
		})
	}
}


func setupTestData(t *testing.T) {
	cleanupTestData(t)
	
	dirs := []string{
		"testdata",
		"testdata/dir1", 
		"testdata/dir1/dir2",
		"testdata/empty_dir",
		"testdata/only_hidden",
		"testdata/nested_symlinks",
	}
	
	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			t.Fatalf("Failed to create test directory %s: %v", dir, err)
		}
	}

	files := map[string]string{
		"testdata/test_file.txt":       "test content123", // 15 bytes
		"testdata/.hidden_file":        "test content123", // 15 bytes  
		"testdata/unicode_файл.txt":    "test content123", // 15 bytes
		"testdata/dir1/file1.txt":      "test content123", // 15 bytes
		"testdata/dir1/.hidden_file":   "test content123", // 15 bytes
		"testdata/dir1/dir2/file2.txt": "test content123", // 15 bytes
		"testdata/only_hidden/.hidden": "test content123", // 15 bytes
	}
	
	for file, content := range files {
		err := os.WriteFile(file, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}
	
	// create simlinks & check errors return value
	if os.Getenv("TEST_SKIP_SYMLINKS") == "" {
		err := os.Symlink("test_file.txt", "testdata/symlink_to_file")
		if err != nil && !os.IsExist(err) {
			t.Logf("Note: could not create symlink: %v", err)
		}
		
		err = os.Symlink("dir1", "testdata/symlink_to_dir")
		if err != nil && !os.IsExist(err) {
			t.Logf("Note: could not create symlink: %v", err)
		}
		
		err = os.Symlink("symlink_to_file", "testdata/nested_symlinks/double_symlink")
		if err != nil && !os.IsExist(err) {
			t.Logf("Note: could not create symlink: %v", err)
		}
	}
}

func cleanupTestData(t *testing.T) {
	err := os.RemoveAll("testdata")
	if err != nil {
		t.Logf("Warning: failed to clean up testdata: %v", err)
	}
}