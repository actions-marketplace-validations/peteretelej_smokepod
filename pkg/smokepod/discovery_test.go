package smokepod

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindTestFiles_SingleFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.test")
	if err := os.WriteFile(testFile, []byte("## section\n$ echo hello\nhello\n"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	files, err := FindTestFiles(testFile)
	if err != nil {
		t.Fatalf("FindTestFiles failed: %v", err)
	}

	if len(files) != 1 {
		t.Errorf("len(files) = %d, want 1", len(files))
	}
}

func TestFindTestFiles_Directory(t *testing.T) {
	tmpDir := t.TempDir()

	testFile1 := filepath.Join(tmpDir, "test1.test")
	testFile2 := filepath.Join(tmpDir, "test2.test")
	otherFile := filepath.Join(tmpDir, "other.txt")

	if err := os.WriteFile(testFile1, []byte("## section\n$ echo hello\nhello\n"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	if err := os.WriteFile(testFile2, []byte("## section\n$ echo world\nworld\n"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	if err := os.WriteFile(otherFile, []byte("not a test file"), 0644); err != nil {
		t.Fatalf("Failed to create other file: %v", err)
	}

	files, err := FindTestFiles(tmpDir)
	if err != nil {
		t.Fatalf("FindTestFiles failed: %v", err)
	}

	if len(files) != 2 {
		t.Errorf("len(files) = %d, want 2", len(files))
	}
}

func TestFindTestFiles_NestedDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	nestedDir := filepath.Join(tmpDir, "nested")
	if err := os.MkdirAll(nestedDir, 0755); err != nil {
		t.Fatalf("Failed to create nested directory: %v", err)
	}

	testFile1 := filepath.Join(tmpDir, "test1.test")
	testFile2 := filepath.Join(nestedDir, "test2.test")

	if err := os.WriteFile(testFile1, []byte("## section\n$ echo hello\nhello\n"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	if err := os.WriteFile(testFile2, []byte("## section\n$ echo world\nworld\n"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	files, err := FindTestFiles(tmpDir)
	if err != nil {
		t.Fatalf("FindTestFiles failed: %v", err)
	}

	if len(files) != 2 {
		t.Errorf("len(files) = %d, want 2", len(files))
	}
}

func TestFindTestFiles_NonTestFile(t *testing.T) {
	tmpDir := t.TempDir()
	otherFile := filepath.Join(tmpDir, "other.txt")
	if err := os.WriteFile(otherFile, []byte("not a test file"), 0644); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	_, err := FindTestFiles(otherFile)
	if err == nil {
		t.Error("FindTestFiles should return error for non-.test file")
	}
}

func TestFindTestFiles_NotFound(t *testing.T) {
	_, err := FindTestFiles("/nonexistent/path")
	if err == nil {
		t.Error("FindTestFiles should return error for nonexistent path")
	}
}
