package main

import (
	"fmt"

	"github.com/timohahaa/report-maker/internal/report"
)

func main() {
	err := report.CreateReport()
	if err != nil {
		fmt.Println("Error when creating report: ", err)
	}
}
