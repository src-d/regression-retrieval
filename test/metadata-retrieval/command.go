package metadataretrieval

import (
	"os"
	"os/exec"
	"syscall"
)

// Command wraps a metadata-retrieval server instance.
type Command struct {
	cmd    *exec.Cmd
	binary string
	orgs   string
	dir    string
}

// NewCommand creates a new metadata-retrieval command struct.
func NewCommand(binary, orgs string) *Command {
	return &Command{
		cmd:    new(exec.Cmd),
		binary: binary,
		orgs:   orgs,
	}
}

// Run runs metadata-retrieval util to discover and download organizations
func (c *Command) Run(envs map[string]string) error {
	c.cmd = exec.Command(
		c.binary,
		"ghsync",
		"--version", "0",
		"--orgs", c.orgs,
		"--no-forks",
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
