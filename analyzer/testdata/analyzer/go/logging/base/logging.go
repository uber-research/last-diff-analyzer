package logging

import "go.uber.org/zap"

var logger = zap.NewExample()

func getLogger() *zap.Logger {
	return logger
}

// rename tests renaming of the logging method
func rename() int {
	logger.Debug("message")
	return 42
}

// modStr tests changing the string argument to the logging method
func modStr() {
	logger.Error("error")
}

// modFunAccess tests changing the call with access path involving a function.
func modFunAccess() {
	getLogger().Error("error")
}

// modAddr tests changing an argument to the logging method to take its address
func modAddr(s string) {
	logger.Error("error", zap.Any("key", s))
}

// genMod tests a more general modification of the argument to the logging method
func genMod(s1 string, s2 string, a []string) {
	logger.Error("error", zap.Int("key", len(s1)))
}

// argRemove tests removing an argument to the logging method
func argRemove(s string) {
	logger.Error("error", zap.String("key", s))
}

// argRemoveCall tests removing a function call argument to the logging method
func argRemoveCall(s string) {
	logger.Error("error", zap.Int("key", len(s)))
}

// argAdd tests adding an argument to the logging method
func argAdd(s string) {
	logger.Error("error")
}

// argAddCall tests adding a function call argument to the logging method
func argAddCall(s string) {
	logger.Error("error")
}

// callRemove tests removing the logging call
func callRemove() string {
	s := "foo"
	logger.Error("error", zap.Any("key", &s))
	return s
}

// callRemoveLast tests removing the logging call being last statement in the function
func callRemoveLast() {
	const s = "foo"
	logger.Error("error", zap.String("key", s))
}

// callAdd tests adding the logging call
func callAdd(s string) string {
	return s
}

var lastLogged = 0

// customInt is a custom version of zap.Int for testing purposes
func customInt(key string, val int) zap.Field {
	lastLogged = val
	return zap.Int(key, val)
}

// argSideEffect tests approval of a logging call with arguments that
// both have side-effects (and do not change) and are side-effect-free
// (and do change)
func argSideEffect() int {
	logger.Error("error", customInt("key1", 7), zap.String("key2", "error"), customInt("key3", 42))
	return lastLogged
}
