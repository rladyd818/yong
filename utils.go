package yong

import (
	"fmt"
	"io"
	"os"
	"strings"
)

var DefaultWriter io.Writer = os.Stdout

func debugPrint(format string, values ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Fprintf(DefaultWriter, "[YONG-info] "+format, values...)
}
