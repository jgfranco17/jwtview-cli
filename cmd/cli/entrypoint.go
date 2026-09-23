package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jgfranco17/jwtview-cli/pkg/logging"
	"github.com/jgfranco17/jwtview-cli/pkg/parsing"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	flagVerbose string = "verbose"
)

func init() {
	cobra.EnableTraverseRunHooks = true
}

type CommandRoot struct {
	rootCmd   *cobra.Command
	verbosity int
}

type Options struct {
	Version string `json:"version"`
	Author  string `json:"author"`
}

func (o *Options) SetDefaults() {
	if o.Version == "" {
		o.Version = "0.0.0"
	}
}

// NewCommandRoot creates a new instance of CommandRegistry
func NewCommandRoot(options Options) *CommandRoot {
	var verbosity int
	options.SetDefaults()

	root := &cobra.Command{
		Use:     "jwtview",
		Version: options.Version,
		Short:   "Display JWT tokens in a human-readable format.",
		Long: `JwtView is a command-line tool for inspecting JSON Web Tokens (JWTs).

It decodes a token and displays its contents in a human-readable format, making
it easier to examine token headers and claims while developing or troubleshooting
authentication flows.

JwtView helps you understand a token's structure and contents; it does not
replace signature verification or other token validation.
`,
		Args:    cobra.ExactArgs(1),
		Example: "jwtview [flags] <raw-token>",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			verbosity, _ := cmd.Flags().GetCount(flagVerbose)
			level := convertCountToLogLevel(verbosity)

			logger := logging.New(cmd.ErrOrStderr(), level)
			ctx := logging.AddToContext(cmd.Context(), logger)

			ctx, cancel := context.WithCancel(ctx)
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
			go func() {
				select {
				case <-c:
					logger.Warn("Received termination signal, shutting down.")
					cancel()
				case <-ctx.Done():
				}
			}()

			cmd.SetContext(ctx)
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			logger := logging.FromContext(ctx)

			rawToken := args[0]
			result, err := parsing.Read(rawToken)
			if err != nil {
				return err
			}
			logger.WithFields(logrus.Fields{
				"headers":   len(result.Header),
				"payload":   len(result.Payload),
				"signature": result.Signature != "",
			}).Info("Decoded JWT successfully.")
			jsonOutput, err := result.ToJSON()
			if err != nil {
				return err
			}
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), jsonOutput)
			return nil
		},
	}

	root.PersistentFlags().CountVarP(&verbosity, flagVerbose, "v", "Increase log verbosity, repeat for more detail (up to -vvv)")
	return &CommandRoot{
		rootCmd:   root,
		verbosity: verbosity,
	}
}

// Execute executes the root command
func (cr *CommandRoot) Execute() error {
	return cr.rootCmd.ExecuteContext(context.Background())
}

func convertCountToLogLevel(count int) logrus.Level {
	switch count {
	case 0:
		return logrus.WarnLevel
	case 1:
		return logrus.InfoLevel
	case 2:
		return logrus.DebugLevel
	case 3:
		return logrus.TraceLevel
	default:
		return logrus.InfoLevel
	}
}
