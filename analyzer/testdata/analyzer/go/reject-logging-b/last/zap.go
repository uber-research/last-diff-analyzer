package zap

import "go.uber.org/zap"

// Int is a custom version of the zap.Int function from the Zap
// package
func Int(key string, value int) zap.Field {
	return zap.Int(key, value)
}
