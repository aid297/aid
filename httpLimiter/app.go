package httpLimiter

var APP struct {
	Visit         Visit
	IPLimiter     IPLimiter
	RouterLimiter RouteLimiter
}
