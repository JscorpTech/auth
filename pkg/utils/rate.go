package utils

import "golang.org/x/time/rate"

var visitors = make(map[string]*rate.Limiter)

func GetVisitor(ip string) *rate.Limiter {
	if limiter, exists := visitors[ip]; exists {
		return limiter
	}

	limiter := rate.NewLimiter(5, 7)
	visitors[ip] = limiter
	return limiter
}
