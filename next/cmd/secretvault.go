package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/twpayne/chezmoi/next/internal/chezmoi"
)

var vaultCmd = &cobra.Command{
	Use:   "vault [args...]",
	Short: "Execute the Hashicorp Vault CLI (vault)",
	RunE:  config.runVaultCmd,
}

type vaultCmdConfig struct {
	Command string
}

var vaultCache = make(map[string]interface{})

func init() {
	config.addTemplateFunc("vault", config.vaultFunc)

	secretCmd.AddCommand(vaultCmd)
}

func (c *Config) runVaultCmd(cmd *cobra.Command, args []string) error {
	return c.run("", c.Vault.Command, args)
}

func (c *Config) vaultFunc(key string) interface{} {
	if data, ok := vaultCache[key]; ok {
		return data
	}
	name := c.Vault.Command
	args := []string{"kv", "get", "-format=json", key}
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	output, err := c.system.IdempotentCmdOutput(cmd)
	if err != nil {
		panic(fmt.Errorf("vault: %s %s: %w\n%s", name, chezmoi.ShellQuoteArgs(args), err, output))
	}
	var data interface{}
	if err := json.Unmarshal(output, &data); err != nil {
		panic(fmt.Errorf("vault: %s %s: %w\n%s", name, chezmoi.ShellQuoteArgs(args), err, output))
	}
	vaultCache[key] = data
	return data
}
