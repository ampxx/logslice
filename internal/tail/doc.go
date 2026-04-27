// Package tail implements a lightweight file-following mechanism for
// logslice's live-streaming mode.
//
// # Overview
//
// A [Tailer] opens a file, seeks to its current end, and then polls for
// newly appended bytes at a configurable interval (default 250 ms). Each
// complete line is sent on a string channel so callers can feed it directly
// into the existing parser / filter / output pipeline.
//
// # Usage
//
//	tr := tail.New("/var/log/app.log")
//	ch, err := tr.Lines(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for line := range ch {
//		fmt.Println(line)
//	}
//
// The channel is closed automatically when the supplied [context.Context] is
// cancelled, making clean shutdown straightforward.
package tail
