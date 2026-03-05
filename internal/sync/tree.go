package sync

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

type LocalNode struct {
	Path     string
	Title    string
	Content  string
	IsDir    bool
	Children []*LocalNode
}

func BuildLocalTree(rootPath string) (*LocalNode, error) {
	info, err := os.Stat(rootPath)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, os.ErrInvalid
	}

	root := &LocalNode{
		Path:  rootPath,
		Title: titleFromPath(filepath.Base(rootPath)),
		IsDir: true,
	}

	if err := walkDir(rootPath, root); err != nil {
		return nil, err
	}

	return root, nil
}

func walkDir(dirPath string, parent *LocalNode) error {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		fullPath := filepath.Join(dirPath, entry.Name())

		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		if entry.IsDir() {
			child := &LocalNode{
				Path:  fullPath,
				Title: titleFromPath(entry.Name()),
				IsDir: true,
			}
			if err := walkDir(fullPath, child); err != nil {
				return err
			}
			if len(child.Children) > 0 {
				parent.Children = append(parent.Children, child)
			}
		} else if strings.HasSuffix(strings.ToLower(entry.Name()), ".md") {
			content, err := os.ReadFile(fullPath)
			if err != nil {
				return err
			}
			title := extractTitle(string(content), entry.Name())
			child := &LocalNode{
				Path:    fullPath,
				Title:   title,
				Content: string(content),
				IsDir:   false,
			}
			parent.Children = append(parent.Children, child)
		}
	}

	return nil
}

var headingRe = regexp.MustCompile(`^#\s+(.+)$`)

func ExtractTitleFromContent(content, filename string) string {
	return extractTitle(content, filename)
}

func extractTitle(content, filename string) string {
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		matches := headingRe.FindStringSubmatch(line)
		if len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		}
		break
	}
	return titleFromPath(strings.TrimSuffix(filename, filepath.Ext(filename)))
}

func titleFromPath(name string) string {
	name = strings.ReplaceAll(name, "-", " ")
	name = strings.ReplaceAll(name, "_", " ")
	words := strings.Fields(name)
	for i, w := range words {
		runes := []rune(w)
		if len(runes) > 0 {
			runes[0] = unicode.ToUpper(runes[0])
		}
		words[i] = string(runes)
	}
	return strings.Join(words, " ")
}

func SanitizeFilename(title string) string {
	title = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' || r == '-' || r == '_' {
			return r
		}
		return -1
	}, title)
	title = strings.ReplaceAll(title, " ", "-")
	title = strings.ToLower(title)
	if title == "" {
		title = "untitled"
	}
	return title
}
