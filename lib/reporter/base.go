package reporter

import (
	"fmt"
	"os"
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
		for _, item := range r.FileUrl {
			res = res + item + "\n"
			return TextReport(res)
		}
		return TextReport(res)
	}
}

func (a TextReport) Print(f *os.File) {
	_, _ = fmt.Fprintln(f, a)
}
