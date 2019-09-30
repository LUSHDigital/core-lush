package log

import (
	"log"
)

const (
	flags = log.Lshortfile | log.LstdFlags
)

func init() {
	// Setup logs as part of import side-effect.
	log.SetFlags(flags)
}
