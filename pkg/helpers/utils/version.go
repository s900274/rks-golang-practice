package utils

import (
	"fmt"
	"runtime"
)

const Binary = "1.0-alpha"

func Version(app string) string {
	return fmt.Sprintf("%s v%s (built w/%s)", app, Binary, runtime.Version())
}
