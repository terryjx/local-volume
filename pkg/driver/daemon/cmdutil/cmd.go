// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package cmdutil

import (
	"bytes"
	"fmt"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

// ExecutableFactory wraps command execution
type ExecutableFactory func(name string, args ...string) Executable

// NewExecutableFactory creates a new ExecutableFactory with an exec.Cmd inside.
func NewExecutableFactory() ExecutableFactory {
	return func(name string, args ...string) Executable {
		c := WrappedCmd{Cmd: exec.Command(name, args...)}
		c.stdErr, c.stdOut = new(bytes.Buffer), new(bytes.Buffer)
		return &c
	}
}

// WrappedCmd wraps an exec.WrappedCmd to match the Executable interface
type WrappedCmd struct {
	*exec.Cmd
	stdOut, stdErr *bytes.Buffer
}

// Command returns the command arguments
func (c *WrappedCmd) Command() []string { return c.Args }

// StdOut returns the stdout
func (c *WrappedCmd) StdOut() []byte { return c.stdOut.Bytes() }

// StdErr returns the stderr
func (c *WrappedCmd) StdErr() []byte { return c.stdErr.Bytes() }

// Run starts the specified command and waits for it to complete.
func (c *WrappedCmd) Run() error {
	c.Cmd.Stderr, c.Cmd.Stdout = c.stdErr, c.stdOut
	return c.Cmd.Run()
}

// Executable defines the common interface that any executable should have.
type Executable interface {
	// CombinedOutput runs the command and returns its combined standard
	// output and standard error.
	CombinedOutput() ([]byte, error)

	// Command returns the command arguments
	Command() []string

	// StdOut returns the stdout
	StdOut() []byte

	// StdErr returns the stderr
	StdErr() []byte

	// Run starts the specified command and waits for it to complete.
	Run() error
}

// RunCmd runs the given command, and returns a combined output
// of the err message, stdout and stderr in case of error
func RunCmd(cmd Executable) error {
	log.Infof("Running command: %v", cmd.Command())
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s. Output: %s", err.Error(), string(output[:]))
	}
	return nil
}
