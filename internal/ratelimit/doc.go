// Package ratelimit implements a token-bucket rate limiter for log entry streams.
//
// # Overview
//
// When processing high-volume log files or tailing active streams it can be
// useful to cap the throughput so that downstream consumers are not
// overwhelmed. The Limiter type uses a simple token-bucket algorithm:
//
//   - A bucket is initialised with `rate` tokens (where rate is entries/second).
//   - Each accepted entry consumes one token.
//   - Tokens are replenished proportionally to the real time elapsed between
//     calls, up to a maximum of `rate` tokens.
//
// # Usage
//
//	l, err := ratelimit.New(100) // allow up to 100 entries/s
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, entry := range entries {
//		if l.Allow(entry) {
//			// forward entry
//		}
//	}
package ratelimit
