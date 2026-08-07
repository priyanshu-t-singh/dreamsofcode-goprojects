package utils

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func PrintTable(data [][]string) {
	tw := tabwriter.NewWriter(os.Stdout, 8, 1, 3, ' ', 0)
	defer tw.Flush()

	if len(data) > 0 && len(data[0]) > 0 {
		for _, colName := range data[0] {
			fmt.Fprintf(tw, "%s\t", colName)
		}
		fmt.Fprintln(tw)
	}

	if len(data) > 1 {
		for _, row := range data[1:] {
			for _, cell := range row {
				fmt.Fprintf(tw, "%s\t", cell)
			}
			fmt.Fprintln(tw)
		}
	}
}
