package httpLimiters

var APP struct {
	Visit         Visit
	IPLimiter     IPLimiter
	RouterLimiter RouteLimiter
}
