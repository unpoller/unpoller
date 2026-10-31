package inputunifi

// collectConsoleSystem is best-effort: an unsupported or inaccessible console
// endpoint must neither fail the Network poll nor manufacture a no-update value.
func (u *InputUnifi) collectConsoleSystem(c *Controller, m *Metrics) {
	system, err := c.Unifi.GetConsoleSystem()
	if err != nil {
		u.LogDebugf("unifi.GetConsoleSystem(%s): %v (continuing)", c.URL, err)

		return
	}

	if c.ID != "" {
		system.SourceName = c.ID
	}

	m.ConsoleSystem = system
}
