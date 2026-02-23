package commands

import "github.com/jongio/azd-core/azdextutil"

// globalRateLimiter uses the shared azdextutil token bucket.
// 10 burst tokens, refills at 1 token/second (≈60/min).
var globalRateLimiter = azdextutil.NewRateLimiter(10, 1.0)
