package cmd

import (
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/twpayne/chezmoi/next/internal/chezmoi"
)

var executeTemplateCmd = &cobra.Command{
	Use:     "execute-template [templates...]",
	Short:   "Execute the given template(s)",
	Long:    mustGetLongHelp("execute-template"),
	Example: getExample("execute-template"),
	RunE:    config.makeRunEWithSourceState(config.runExecuteTemplateCmd),
}

type executeTemplateCmdConfig struct {
	init         bool
	promptString map[string]string
}

func init() {
	rootCmd.AddCommand(executeTemplateCmd)

	persistentFlags := executeTemplateCmd.PersistentFlags()
	persistentFlags.BoolVarP(&config.executeTemplate.init, "init", "i", config.executeTemplate.init, "simulate chezmoi init")
	persistentFlags.StringToStringVarP(&config.executeTemplate.promptString, "promptString", "p", config.executeTemplate.promptString, "simulate promptString")
}

func (c *Config) runExecuteTemplateCmd(cmd *cobra.Command, args []string, sourceState *chezmoi.SourceState) error {
	if c.executeTemplate.init {
		c.templateFuncs["promptString"] = func(prompt string) string {
			if value, ok := c.executeTemplate.promptString[prompt]; ok {
				return value
			}
			return prompt
		}
	}

	output := &strings.Builder{}
	switch len(args) {
	case 0:
		data, err := ioutil.ReadAll(c.stdin)
		if err != nil {
			return err
		}
		result, err := sourceState.ExecuteTemplateData("stdin", data)
		if err != nil {
			return err
		}
		if _, err = output.Write(result); err != nil {
			return err
		}
	default:
		for i, arg := range args {
			result, err := sourceState.ExecuteTemplateData("arg"+strconv.Itoa(i+1), []byte(arg))
			if err != nil {
				return err
			}
			if _, err := output.Write(result); err != nil {
				return err
			}
		}
	}

	return c.writeOutput([]byte(output.String()))
}
