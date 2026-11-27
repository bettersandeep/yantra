package referee

import "strings"

// Change is one changed file between two commits. Status is the first letter
// of the git status: A, M, D, R, C or T. Old and New are equal unless the file
// was renamed or copied.
type Change struct {
	Status   byte
	Old, New string
}

// inScope reports whether a path is covered by the scope paths of a ticket.
// A scope path that ends with "/" covers all files below it. Other scope
// paths match one file exactly.
// ponytail: no globs and no path cleaning. A scope path such as "./cmd/" or
// "cmd/**" matches nothing, so its files read as out of scope.
func inScope(path string, scope []string) bool {
	for _, s := range scope {
		if s == "" {
			continue
		}
		if strings.HasSuffix(s, "/") {
			if strings.HasPrefix(path, s) {
				return true
			}
		} else if path == s {
			return true
		}
	}
	return false
}

// OutOfScope reports each changed file that is outside the scope paths. A
// renamed file is in scope only when its old path and its new path are both
// in scope.
func OutOfScope(changes []Change, scope []string) []Finding {
	fs := []Finding{}
	for _, c := range changes {
		switch {
		case c.Status == 'A' || c.Status == 'C': // a copy leaves the old file as it was
			if !inScope(c.New, scope) {
				fs = append(fs, Finding{File: c.New, Kind: "out_of_scope", Detail: "new file outside the ticket's scope"})
			}
		case c.Status == 'R':
			if !inScope(c.Old, scope) || !inScope(c.New, scope) {
				fs = append(fs, Finding{File: c.New, Kind: "out_of_scope", Detail: "renamed from " + c.Old + ", and one of the two paths is outside the ticket's scope"})
			}
		default: // M, D, T and each status that this code does not know
			if !inScope(c.Old, scope) || !inScope(c.New, scope) {
				fs = append(fs, Finding{File: c.New, Kind: "out_of_scope", Detail: "changed file outside the ticket's scope"})
			}
		}
	}
	return fs
}
