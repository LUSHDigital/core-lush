package nullable

import (
	"time"
)

// StdTime provides default SQL TIME format
const StdTime = "15:04:05"

// emptyTime allows default times to be considered
// null for insertion into the database.
var emptyTime = time.Time{}

// nullLiteral is helpful for checking
// for nulls, as they won't cause errors,
// yet we need the content of the file to change anyway
var nullLiteral = []byte("null")
