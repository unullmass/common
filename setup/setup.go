package setup

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// Task defines a Setup Task. Run() executes the setup task, and Validate() checks whether or not the task succeeded.
// Validate() can and should be run as the first statement of Run() so needless work isn't done again.
type Task interface {
	Run(c Context) error
	Validate(c Context) error
}

// Runner is a task runner for generic Task interfaces. It stores a list of Tasks that will be executed in the
// order they are found in the list. The runner can also be configured to opt-in to ask for user input from stdin
type Runner struct {
	Tasks    []Task
	AskInput bool
}

// Context contains contextual setup runner information
// if askInput is false (default value), the setup task should NOT block and wait for user input
type Context struct {
	askInput bool
}

// RunTasks executes the specified set of Tasks against the registered list of tasks. Any tasks registered that arent in the list provided are skipped.
func (r *Runner) RunTasks(tasks ...string) error {
	fmt.Println("Running setup ...")
	ctx := Context{
		askInput: r.AskInput,
	}
	if len(tasks) == 0 {
		// run ALL the setup tasks
		for _, t := range r.Tasks {
			if err := t.Run(ctx); err != nil {
				return err
			}
			if err := t.Validate(ctx); err != nil {
				return err
			}
		}
	} else {
		// map each task ...string into a map[string]bool
		enabledTasks := make(map[string]bool)
		for _, t := range tasks {
			enabledTasks[strings.ToLower(t)] = true
		}
		// iterate through the proper order of tasks, and execute the ones listed in the parameters
		for _, t := range r.Tasks {
			taskName := strings.ToLower(reflect.TypeOf(t).Name())
			if _, ok := enabledTasks[taskName]; ok {
				if err := t.Run(ctx); err != nil {
					return err
				}
				if err := t.Validate(ctx); err != nil {
					return err
				}
			}
		}
	}
	fmt.Println("Setup finished successfully!")
	return nil
}

// GetenvInt retrieves an integer variable from the environment
// this function will optionally read input from stdin if it was not defined in the environment,
// if Context.askInput is set to true
func (c Context) GetenvInt(env string, description string) (int, error) {
	fmt.Printf("%s:\n", description)
	if intStr, ok := os.LookupEnv(env); ok {
		val, err := strconv.ParseInt(intStr, 10, 32)
		if err == nil {
			fmt.Println(intStr)
			return int(val), nil
		}
		return 0, fmt.Errorf("%s is not not an integer", env)
	}
	if c.askInput {
		var intValue int
		if scanned, err := fmt.Scanf("%d", &intValue); scanned != 1 || err != nil {
			return 0, fmt.Errorf("error reading value for %s", env)
		}
		return intValue, nil
	}
	return 0, fmt.Errorf("%s is not defined", env)
}

// GetenvString retrieves a string variable from the environment
// this function will optionally read input from stdin if it was not defined in the environment,
// if Context.askInput is set to true
func (c Context) GetenvString(env string, description string) (string, error) {
	fmt.Printf("%s:\n", description)
	if str, ok := os.LookupEnv(env); ok {
		fmt.Println(str)
		return str, nil
	}
	if c.askInput {
		var str string
		if scanned, err := fmt.Scanln(&str); scanned != 1 || err != nil {
			return "", fmt.Errorf("error reading value for %s", env)
		}
		return str, nil
	}
	return "", fmt.Errorf("%s is not defined", env)
}

// GetenvSecret retrieves a string variable from the envrionment that is secret
// this is functionally equivalent to GetenvString, but does not print the read value to stdout
// this function will optionally read input from stdin if it was not defined in the environment,
// if Context.askInput is set to true
func (c Context) GetenvSecret(env string, description string) (string, error) {
	fmt.Printf("%s:\n", description)
	if str, ok := os.LookupEnv(env); ok {
		fmt.Println("****")
		return str, nil
	}
	if c.askInput {
		var str string
		if scanned, err := fmt.Scanln(&str); scanned != 1 || err != nil {
			return "", fmt.Errorf("error reading value for %s", env)
		}
		return str, nil
	}
	return "", fmt.Errorf("%s is not defined", env)
}
