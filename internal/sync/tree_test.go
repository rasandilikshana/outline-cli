package sync

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTitleFromPath(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"getting-started", "Getting Started"},
		{"api_reference", "Api Reference"},
		{"hello-world-doc", "Hello World Doc"},
		{"simple", "Simple"},
		{"UPPERCASE", "UPPERCASE"},
		{"mixed-Case_test", "Mixed Case Test"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := titleFromPath(tt.input)
			if result != tt.expected {
				t.Errorf("titleFromPath(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExtractTitle(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		filename string
		expected string
	}{
		{
			name:     "h1 heading",
			content:  "# My Document\n\nSome content here.",
			filename: "doc.md",
			expected: "My Document",
		},
		{
			name:     "h1 with leading whitespace",
			content:  "\n\n# Title After Blanks\n\nContent",
			filename: "doc.md",
			expected: "Title After Blanks",
		},
		{
			name:     "no heading uses filename",
			content:  "Just some text without a heading.",
			filename: "my-document.md",
			expected: "My Document",
		},
		{
			name:     "empty content uses filename",
			content:  "",
			filename: "getting-started.md",
			expected: "Getting Started",
		},
		{
			name:     "h2 not used as title",
			content:  "## Sub Heading\n\nContent",
			filename: "fallback-name.md",
			expected: "Fallback Name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractTitle(tt.content, tt.filename)
			if result != tt.expected {
				t.Errorf("extractTitle() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"My Document", "my-document"},
		{"Hello World!", "hello-world"},
		{"API Reference", "api-reference"},
		{"test@#$%file", "testfile"},
		{"", "untitled"},
		{"already-clean", "already-clean"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := SanitizeFilename(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeFilename(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestBuildLocalTree(t *testing.T) {
	// Create temp directory structure
	tmpDir := t.TempDir()

	// Create subdirectories
	subDir := filepath.Join(tmpDir, "getting-started")
	os.MkdirAll(subDir, 0755)

	// Create markdown files
	os.WriteFile(filepath.Join(tmpDir, "overview.md"), []byte("# Overview\n\nMain overview."), 0644)
	os.WriteFile(filepath.Join(subDir, "intro.md"), []byte("# Introduction\n\nWelcome."), 0644)
	os.WriteFile(filepath.Join(subDir, "setup.md"), []byte("# Setup Guide\n\nHow to set up."), 0644)

	// Create a non-md file (should be ignored)
	os.WriteFile(filepath.Join(tmpDir, "notes.txt"), []byte("ignored"), 0644)

	// Create a dotfile directory (should be ignored)
	os.MkdirAll(filepath.Join(tmpDir, ".hidden"), 0755)
	os.WriteFile(filepath.Join(tmpDir, ".hidden", "secret.md"), []byte("# Secret"), 0644)

	tree, err := BuildLocalTree(tmpDir)
	if err != nil {
		t.Fatalf("BuildLocalTree() error: %v", err)
	}

	if !tree.IsDir {
		t.Error("Root should be a directory")
	}

	// Should have 2 children: overview.md and getting-started/
	if len(tree.Children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(tree.Children))
	}

	// Find the subdirectory
	var dirNode, fileNode *LocalNode
	for _, child := range tree.Children {
		if child.IsDir {
			dirNode = child
		} else {
			fileNode = child
		}
	}

	if dirNode == nil {
		t.Fatal("Expected a directory child")
	}
	if dirNode.Title != "Getting Started" {
		t.Errorf("Dir title = %q, want %q", dirNode.Title, "Getting Started")
	}
	if len(dirNode.Children) != 2 {
		t.Errorf("Expected 2 children in subdir, got %d", len(dirNode.Children))
	}

	if fileNode == nil {
		t.Fatal("Expected a file child")
	}
	if fileNode.Title != "Overview" {
		t.Errorf("File title = %q, want %q", fileNode.Title, "Overview")
	}
	if fileNode.Content == "" {
		t.Error("File content should not be empty")
	}
}

func TestBuildLocalTree_NonExistent(t *testing.T) {
	_, err := BuildLocalTree("/nonexistent/path")
	if err == nil {
		t.Error("Expected error for non-existent path")
	}
}

func TestBuildLocalTree_FileNotDir(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "file.md")
	os.WriteFile(tmpFile, []byte("content"), 0644)

	_, err := BuildLocalTree(tmpFile)
	if err == nil {
		t.Error("Expected error when passing a file instead of directory")
	}
}

func TestExtractTitleFromContent(t *testing.T) {
	result := ExtractTitleFromContent("# Public Title\n\nContent", "file.md")
	if result != "Public Title" {
		t.Errorf("ExtractTitleFromContent() = %q, want %q", result, "Public Title")
	}
}
