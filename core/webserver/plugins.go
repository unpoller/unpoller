package webserver

import "sync"

type webPlugins struct {
	*Config
	inputs       []*Input
	outputs      []*Output
	sync.RWMutex // Locks both of the above slices.
}

// This is global so other plugins can call its methods.
var plugins = &webPlugins{} // nolint: gochecknoglobals

// UpdateInput allows an input plugin to create an entry or update an existing entry.
func UpdateInput(config *Input) {
	if plugins.Enable {
		plugins.updateInput(config)
	}
}

// NewInputEvent adds an event for an input plugin.
func NewInputEvent(name, id string, event *Event) {
	if plugins.Enable {
		plugins.newInputEvent(name, id, event)
	}
}

// NewOutputEvent adds an event for an output plugin.
func NewOutputEvent(name, id string, event *Event) {
	if plugins.Enable {
		plugins.newOutputEvent(name, id, event)
	}
}

// UpdateOutput allows an output plugin to create an entry or update an existing entry.
func UpdateOutput(config *Output) {
	if plugins.Enable {
		plugins.updateOutput(config)
	}
}

// UpdateOutputCounter allows an output plugin to update a counter's value.
func UpdateOutputCounter(plugin, label string, values ...int64) {
	if plugins.Enable {
		plugins.updateOutputCounter(plugin, label, values...)
	}
}

// UpdateInputCounter allows an input plugin to update a counter's value.
// Set any arbitrary counter. These are displayed on the web interface.
func UpdateInputCounter(plugin, label string, values ...int64) {
	if plugins.Enable {
		plugins.updateInputCounter(plugin, label, values...)
	}
}

func (w *webPlugins) updateOutputCounter(plugin, label string, values ...int64) {
	if len(values) == 0 {
		values = []int64{1}
	}

	output := w.getOutput(plugin)
	if output == nil {
		return
	}

	output.Lock()
	defer output.Unlock()

	if output.Counter == nil {
		output.Counter = make(map[string]int64)
	}

	for _, v := range values {
		output.Counter[label] += v
	}
}

func (w *webPlugins) updateInputCounter(plugin, label string, values ...int64) {
	if len(values) == 0 {
		values = []int64{1}
	}

	input := w.getInput(plugin)
	if input == nil {
		return
	}

	input.Lock()
	defer input.Unlock()

	if input.Counter == nil {
		input.Counter = make(map[string]int64)
	}

	for _, v := range values {
		input.Counter[label] += v
	}
}

func (w *webPlugins) updateInput(config *Input) {
	if config == nil {
		return
	}

	input := w.getInput(config.Name)
	if input == nil {
		w.newInput(config)
		return
	}

	config.Lock()
	defer config.Unlock()

	if config.Clients != nil {
		input.Clients = config.Clients
	}

	if config.Sites != nil {
		input.Sites = config.Sites
	}

	if config.Devices != nil {
		input.Devices = config.Devices
	}

	if config.Config != nil {
		input.Config = config.Config
	}

	if config.Counter != nil {
		input.Counter = config.Counter
	}
}

func (w *webPlugins) updateOutput(config *Output) {
	if config == nil || config.Config == nil {
		return
	}

	output := w.getOutput(config.Name)
	if output == nil {
		w.newOutput(config)
		return
	}

	config.Lock()
	defer config.Unlock()

	if config.Config != nil {
		output.Config = config.Config
	}

	if config.Counter != nil {
		output.Counter = config.Counter
	}
}

func (w *webPlugins) newInputEvent(plugin, id string, event *Event) {
	input := w.getInput(plugin)
	if input == nil {
		return
	}

	input.Lock()
	defer input.Unlock()

	if input.Events == nil {
		input.Events = make(map[string]*Events)
	}

	if _, ok := input.Events[id]; !ok {
		input.Events[id] = &Events{}
	}

	input.Events[id].add(event, int(w.Config.MaxEvents))
}

func (w *webPlugins) newOutputEvent(plugin, id string, event *Event) {
	output := w.getOutput(plugin)
	if output == nil {
		return
	}

	output.Lock()
	defer output.Unlock()

	if output.Events == nil {
		output.Events = make(map[string]*Events)
	}

	if _, ok := output.Events[id]; !ok {
		output.Events[id] = &Events{}
	}

	output.Events[id].add(event, int(w.Config.MaxEvents))
}

func (w *webPlugins) newInput(config *Input) {
	w.Lock()
	defer w.Unlock()
	w.inputs = append(w.inputs, config)
}

func (w *webPlugins) newOutput(config *Output) {
	w.Lock()
	defer w.Unlock()
	w.outputs = append(w.outputs, config)
}

func (w *webPlugins) getInput(name string) *Input {
	w.RLock()
	defer w.RUnlock()

	for i := range w.inputs {
		if w.inputs[i].Name == name {
			return w.inputs[i]
		}
	}

	return nil
}

func (w *webPlugins) getOutput(name string) *Output {
	w.RLock()
	defer w.RUnlock()

	for i := range w.outputs {
		if w.outputs[i].Name == name {
			return w.outputs[i]
		}
	}

	return nil
}
