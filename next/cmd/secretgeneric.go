package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"github.com/twpayne/chezmoi/next/internal/chezmoi"
)

var genericSecretCmd = &cobra.Command{
	Use:   "generic [args...]",
	Short: "Execute a generic secret command",
	RunE:  config.runGenericSecretCmd,
}

type genericSecretCmdConfig struct {
	Command string
}

var (
	secretCache     = make(map[string]string)
	secretJSONCache = make(map[string]interface{})
)

func init() {
	config.addTemplateFunc("secret", config.secretFunc)
	config.addTemplateFunc("secretJSON", config.secretJSONFunc)

	secretCmd.AddCommand(genericSecretCmd)
}

func (c *Config) runGenericSecretCmd(cmd *cobra.Command, args []string) error {
	return c.run("", c.GenericSecret.Command, args)
}

func (c *Config) secretFunc(args ...string) string {
	key := strings.Join(args, "\x00")
	if value, ok := secretCache[key]; ok {
		return value
	}
	name := c.GenericSecret.Command
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	output, err := c.system.IdempotentCmdOutput(cmd)
	if err != nil {
		panic(fmt.Errorf("secret: %s %s: %w\n%s", name, chezmoi.ShellQuoteArgs(args), err, output))
	}
	value := string(bytes.TrimSpace(output))
	secretCache[key] = value
	return value
}

func (c *Config) secretJSONFunc(args ...string) interface{} {
	key := strings.Join(args, "\x00")
	if value, ok := secretJSONCache[key]; ok {
		return value
	}
	name := c.GenericSecret.Command
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	output, err := c.system.IdempotentCmdOutput(cmd)
	if err != nil {
		panic(fmt.Errorf("secretJSON: %s %s: %w\n%s", name, chezmoi.ShellQuoteArgs(args), err, output))
	}
	var value interface{}
	if err := json.Unmarshal(output, &value); err != nil {
		panic(fmt.Errorf("secretJSON: %s %s: %w\n%s", name, chezmoi.ShellQuoteArgs(args), err, output))
	}
	secretJSONCache[key] = value
	return value
}
