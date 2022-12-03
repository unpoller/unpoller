package datadogunifi

// Logf logs a message.
func (u *DatadogUnifi) Logf(msg string, v ...any) {
	if u.Collector != nil {
		u.Collector.Logf(msg, v...)
	}
}

// LogErrorf logs an error message.
func (u *DatadogUnifi) LogErrorf(msg string, v ...any) {
	if u.Collector != nil {
		u.Collector.LogErrorf(msg, v...)
	}
}

// LogDebugf logs a debug message.
func (u *DatadogUnifi) LogDebugf(msg string, v ...any) {
	if u.Collector != nil {
		u.Collector.LogDebugf(msg, v...)
	}
}
