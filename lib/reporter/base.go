package reporter

import (
	"fmt"
	"strings"
)

type TextReport string
type TextReporter func(r ReporterSchema) TextReport

type ReporterSchema struct {
	FileUrl []string
	Error   error
}

func TyporaReporter() TextReporter {
	return func(r ReporterSchema) TextReport {
		res := "Upload Success:\n"
		if r.Error != nil {
			res = res + r.Error.Error() + "\n"
			return TextReport(res)
		}
		res += strings.Join(r.FileUrl, "\n")
		return TextReport(res)
	}
}

func (a TextReport) Print() {
	fmt.Print(a)
}
