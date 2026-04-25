package stats

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

const timeLayout = "2006-01-02 15:04:05 UTC"

// Print writes a human-readable summary table to w.
func Print(w io.Writer, s Summary) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "--- log summary ---")
	fmt.Fprintf(tw, "Total entries:\t%d\n", s.Total)
	fmt.Fprintf(tw, "Matched:\t%d\n", s.Matched)
	fmt.Fprintf(tw, "Skipped:\t%d\n", s.Skipped)
	if s.PatternHit > 0 {
		fmt.Fprintf(tw, "Pattern hits:\t%d\n", s.PatternHit)
	}
	if !s.Earliest.IsZero() {
		fmt.Fprintf(tw, "Earliest:\t%s\n", s.Earliest.UTC().Format(timeLayout))
	}
	if !s.Latest.IsZero() {
		fmt.Fprintf(tw, "Latest:\t%s\n", s.Latest.UTC().Format(timeLayout))
	}
	if !s.Earliest.IsZero() && !s.Latest.IsZero() && s.Latest.After(s.Earliest) {
		dur := s.Latest.Sub(s.Earliest).Round(time.Second)
		fmt.Fprintf(tw, "Span:\t%s\n", dur)
	}
	_ = tw.Flush()
}
