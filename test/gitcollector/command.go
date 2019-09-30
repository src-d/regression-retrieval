package gitcollector

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/src-d/regression-core"
)

// Command wraps a gitcollector server instance.
type Command struct {
	cmd    *exec.Cmd
	binary string
	orgs   string
	dir    string
}

// NewCommand creates a new gitcollector command struct.
func NewCommand(binary, orgs string) *Command {
	return &Command{
		cmd:    new(exec.Cmd),
		binary: binary,
		orgs:   orgs,
	}
}

// Run runs gitcollector util to discover and download organizations
func (c *Command) Run(envs map[string]string) error {
	tmpDir, err := regression.CreateTempDir()
	if err != nil {
		return err
	}
	c.dir = tmpDir
	defer c.Cleanup()

	c.cmd = exec.Command(
		c.binary,
		"download",
		"--library", c.dir,
		"--orgs", c.orgs,
	)
	c.cmd.Stdout = os.Stdout
	c.cmd.Stderr = os.Stderr
	c.cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	for k, v := range envs {
		c.cmd.Env = append(c.cmd.Env, k+"="+v)
	}

	return c.cmd.Run()
}

// Rusage returns usage counters
func (c *Command) Rusage() *syscall.Rusage {
	rusage, _ := c.cmd.ProcessState.SysUsage().(*syscall.Rusage)
	return rusage
}

// Cleanup removes gitcollector's library directory
func (c *Command) Cleanup() error {
	return os.RemoveAll(c.dir)
}
