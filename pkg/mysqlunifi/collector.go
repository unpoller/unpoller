package main

/* Everything in this file runs after the config is unmarshalled and we've
   verified the configuration for the poller. */

func (p *plugin) runCollector() error {
	p.Logf("mysql plugin is not finished")

	return nil
}
