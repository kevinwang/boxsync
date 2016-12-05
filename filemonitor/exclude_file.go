package filemonitor

import (
	"path/filepath"
	"strings"
)

type ExcludePattern struct {
	patterns []string
}

type ExcludeFiles struct {
	files []string
}
