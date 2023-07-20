package poller

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

// PrintRawMetrics prints raw json from the UniFi Controller. This is currently
// tied into the -j CLI arg, and is probably not very useful outside that context.
func (u *UnifiPoller) PrintRawMetrics() (err error) {
	split := strings.SplitN(u.Flags.DumpJSON, " ", 2)
	filter := &Filter{Kind: split[0]}

	// Allows you to grab a controller other than 0 from config.
	if split2 := strings.Split(filter.Kind, ":"); len(split2) > 1 {
		filter.Kind = split2[0]
		filter.Unit, _ = strconv.Atoi(split2[1])
	}

	// Used with "other"
	if len(split) > 1 {
		filter.Path = split[1]
	}

	// As of now we only have one input plugin, so target that [0].
	m, err := inputs[0].RawMetrics(filter)
	fmt.Println(string(m))

	return err
}

// PrintPasswordHash prints a bcrypt'd password. Useful for the web server.
func (u *UnifiPoller) PrintPasswordHash() (err error) {
	pwd := []byte(u.Flags.HashPW)

	if u.Flags.HashPW == "-" {
		fmt.Print("Enter Password: ")

		pwd, err = term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return fmt.Errorf("reading stdin: %w", err)
		}

		fmt.Println() // print a newline.
	}

	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	fmt.Println(string(hash))

	return err //nolint:wrapcheck
}

func (u *UnifiPoller) DebugIO() error {
	inputSync.RLock()
	defer inputSync.RUnlock()
	outputSync.RLock()
	defer outputSync.RUnlock()

	allOK := true

	var allErr error

	u.Logf("Checking inputs...")

	totalInputs := len(inputs)

	for i, input := range inputs {
		u.Logf("\t(%d/%d) Checking input %s...", i+1, totalInputs, input.Name)

		ok, err := input.DebugInput()
		if !ok {
			u.LogErrorf("\t\t %s Failed: %v", input.Name, err)

			allOK = false
		} else {
			u.Logf("\t\t %s is OK", input.Name)
		}

		if err != nil {
			if allErr == nil {
				allErr = err
			} else {
				allErr = fmt.Errorf("%v: %w", err, allErr)
			}
		}
	}

	u.Logf("Checking outputs...")

	totalOutputs := len(outputs)

	for i, output := range outputs {
		u.Logf("\t(%d/%d) Checking output %s...", i+1, totalOutputs, output.Name)

		ok, err := output.DebugOutput()
		if !ok {
			u.LogErrorf("\t\t %s Failed: %v", output.Name, err)

			allOK = false
		} else {
			u.Logf("\t\t %s is OK", output.Name)
		}

		if err != nil {
			if allErr == nil {
				allErr = err
			} else {
				allErr = fmt.Errorf("%v: %w", err, allErr)
			}
		}
	}

	if !allOK {
		u.LogErrorf("No all checks passed, please fix the logged issues.")
	}

	return allErr
}
