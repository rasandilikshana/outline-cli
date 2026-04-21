package mcp

import "fmt"

func errMissing(field string) error {
	return fmt.Errorf("missing required field: %s", field)
}

func errNotFound(kind, ident string) error {
	return fmt.Errorf("%s not found: %s", kind, ident)
}
