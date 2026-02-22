package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/jongio/azd-core/azdextutil"
	"github.com/jongio/azd-core/security"
	"github.com/jongio/azd-core/shellutil"
	"github.com/jongio/azd-exec/cli/src/internal/version"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
)

const defaultTimeout = 30 * time.Second

// NewMCPCommand creates the mcp parent command.
func NewMCPCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "mcp",
		Short:  "Model Context Protocol server operations",
		Long:   `Manage the Model Context Protocol (MCP) server for the azd exec extension.`,
		Hidden: true,
	}
	cmd.AddCommand(newMCPServeCommand())
	return cmd
}

func newMCPServeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Start the MCP server",
		Long:  `Starts the Model Context Protocol server using stdio transport.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMCPServer(cmd.Context())
		},
	}
}

func runMCPServer(_ context.Context) error {
	instructions := `This MCP server is provided by the azd exec extension for the Azure Developer CLI.

**Extension Role:**
Execute scripts and commands with full Azure Developer CLI context, including environment
variables and Azure Key Vault secret resolution.

**Tool Categories:**
- Execution: exec_script, exec_inline - Run scripts/commands with azd environment context
- Discovery: list_shells - Discover available shells on the system
- Configuration: get_environment - View current azd environment variables

**Best Practices:**
- Always verify script paths before execution
- Use list_shells to discover available shells before specifying one
- Use get_environment to check available environment variables
- Prefer exec_script for file-based scripts, exec_inline for one-liners
- Be cautious with destructive operations; review commands before executing`

	s := server.NewMCPServer(
		"exec-mcp-server",
		version.Version,
		server.WithToolCapabilities(true),
		server.WithInstructions(instructions),
	)

	s.AddTools(
		newExecScriptTool(),
		newExecInlineTool(),
		newListShellsTool(),
		newGetEnvironmentTool(),
	)

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "MCP server error: %v\n", err)
		return err
	}
	return nil
}

// --- exec_script tool ---

func newExecScriptTool() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(
			"exec_script",
			mcp.WithTitleAnnotation("Execute Script File"),
			mcp.WithDescription("Execute a script file with azd environment context and Key Vault integration. "+
				"The script runs with all azd environment variables available, including resolved Key Vault secrets."),
			mcp.WithReadOnlyHintAnnotation(false),
			mcp.WithDestructiveHintAnnotation(true),
			mcp.WithString("script_path",
				mcp.Description("Path to the script file to execute. Must be an existing file within the project directory."),
				mcp.Required(),
			),
			mcp.WithString("shell",
				mcp.Description("Shell to use for execution (bash, sh, zsh, pwsh, powershell, cmd). Auto-detected from file extension if not specified."),
			),
			mcp.WithString("args",
				mcp.Description("Space-separated arguments to pass to the script."),
			),
		),
		Handler: handleExecScript,
	}
}

func handleExecScript(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !globalRateLimiter.Allow() {
		return mcp.NewToolResultError("Rate limit exceeded. Please wait before making more requests."), nil
	}

	args := getArgsMap(request)

	scriptPath, ok := getStringParam(args, "script_path")
	if !ok || scriptPath == "" {
		return mcp.NewToolResultError("script_path is required"), nil
	}

	shell, _ := getStringParam(args, "shell")
	if shell != "" {
		if err := azdextutil.ValidateShellName(shell); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid shell: %v", err)), nil
		}
	}

	// Validate script path for security
	projectDir, err := azdextutil.GetProjectDir("AZD_EXEC_PROJECT_DIR")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to determine project directory: %v", err)), nil
	}

	validPath, err := security.ValidatePathWithinBases(scriptPath, projectDir)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid script path: %v", err)), nil
	}

	info, statErr := os.Stat(validPath)
	if statErr != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Script file not found: %s", scriptPath)), nil
	}
	if info.IsDir() {
		return mcp.NewToolResultError("script_path must be a file, not a directory"), nil
	}

	// Parse extra args
	var scriptArgs []string
	if argsStr, ok := getStringParam(args, "args"); ok && argsStr != "" {
		scriptArgs = strings.Fields(argsStr)
	}

	// Detect shell
	if shell == "" {
		shell = shellutil.DetectShell(validPath)
	}

	// Build and execute command
	timeout := defaultTimeout
	execCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmdArgs := buildShellArgs(shell, validPath, false, scriptArgs)
	cmd := exec.CommandContext(execCtx, cmdArgs[0], cmdArgs[1:]...)
	cmd.Env = os.Environ()

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	runErr := cmd.Run()
	return marshalExecResult(stdout.String(), stderr.String(), cmd.ProcessState, runErr)
}

// --- exec_inline tool ---

func newExecInlineTool() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(
			"exec_inline",
			mcp.WithTitleAnnotation("Execute Inline Command"),
			mcp.WithDescription("Execute an inline command with azd environment context. "+
				"The command runs with all azd environment variables, including resolved Key Vault secrets."),
			mcp.WithReadOnlyHintAnnotation(false),
			mcp.WithDestructiveHintAnnotation(true),
			mcp.WithString("command",
				mcp.Description("The command to execute inline."),
				mcp.Required(),
			),
			mcp.WithString("shell",
				mcp.Description("Shell to use (bash, sh, zsh, pwsh, powershell, cmd). Defaults to bash on Unix, powershell on Windows."),
			),
		),
		Handler: handleExecInline,
	}
}

func handleExecInline(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !globalRateLimiter.Allow() {
		return mcp.NewToolResultError("Rate limit exceeded. Please wait before making more requests."), nil
	}

	args := getArgsMap(request)

	command, ok := getStringParam(args, "command")
	if !ok || strings.TrimSpace(command) == "" {
		return mcp.NewToolResultError("command is required and cannot be empty"), nil
	}

	shell, _ := getStringParam(args, "shell")
	if shell != "" {
		if err := azdextutil.ValidateShellName(shell); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid shell: %v", err)), nil
		}
	}
	if shell == "" {
		if runtime.GOOS == "windows" {
			shell = shellutil.ShellPowerShell
		} else {
			shell = shellutil.ShellBash
		}
	}

	timeout := defaultTimeout
	execCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmdArgs := buildShellArgs(shell, command, true, nil)
	cmd := exec.CommandContext(execCtx, cmdArgs[0], cmdArgs[1:]...)
	cmd.Env = os.Environ()

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	runErr := cmd.Run()
	return marshalExecResult(stdout.String(), stderr.String(), cmd.ProcessState, runErr)
}

// --- list_shells tool ---

type shellInfo struct {
	Name      string `json:"name"`
	Available bool   `json:"available"`
}

func newListShellsTool() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(
			"list_shells",
			mcp.WithTitleAnnotation("List Available Shells"),
			mcp.WithDescription("List shells available on the system for script execution."),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
			mcp.WithDestructiveHintAnnotation(false),
		),
		Handler: handleListShells,
	}
}

func handleListShells(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !globalRateLimiter.Allow() {
		return mcp.NewToolResultError("Rate limit exceeded. Please wait before making more requests."), nil
	}

	shells := []string{
		shellutil.ShellBash,
		shellutil.ShellSh,
		shellutil.ShellZsh,
		shellutil.ShellPwsh,
		shellutil.ShellPowerShell,
		shellutil.ShellCmd,
	}

	var results []shellInfo
	for _, sh := range shells {
		_, err := exec.LookPath(sh)
		results = append(results, shellInfo{
			Name:      sh,
			Available: err == nil,
		})
	}

	return marshalToolResult(results)
}

// --- get_environment tool ---

type envVar struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func newGetEnvironmentTool() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(
			"get_environment",
			mcp.WithTitleAnnotation("Get Environment Variables"),
			mcp.WithDescription("Get current azd environment variables available for script execution."),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
			mcp.WithDestructiveHintAnnotation(false),
		),
		Handler: handleGetEnvironment,
	}
}

func handleGetEnvironment(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !globalRateLimiter.Allow() {
		return mcp.NewToolResultError("Rate limit exceeded. Please wait before making more requests."), nil
	}

	allowedPrefixes := []string{"AZD_", "AZURE_", "ARM_", "DOTNET_", "NODE_", "PYTHON"}

	var vars []envVar
	for _, e := range os.Environ() {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			name := parts[0]
			key := strings.ToUpper(name)
			matched := false
			for _, prefix := range allowedPrefixes {
				if strings.HasPrefix(key, prefix) {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}

			// Exclude known secret-bearing variable names
			secretPatterns := []string{"SECRET", "PASSWORD", "KEY", "TOKEN", "CREDENTIAL", "CERTIFICATE"}
			isSecret := false
			upperName := strings.ToUpper(name)
			for _, pattern := range secretPatterns {
				if strings.Contains(upperName, pattern) {
					isSecret = true
					break
				}
			}
			if isSecret {
				continue
			}

			vars = append(vars, envVar{Key: name, Value: parts[1]})
		}
	}

	return marshalToolResult(vars)
}

// --- Helpers ---

func getArgsMap(request mcp.CallToolRequest) map[string]interface{} {
	if request.Params.Arguments != nil {
		if m, ok := request.Params.Arguments.(map[string]interface{}); ok {
			return m
		}
	}
	return map[string]interface{}{}
}

func getStringParam(args map[string]interface{}, key string) (string, bool) {
	val, ok := args[key]
	if !ok {
		return "", false
	}
	s, ok := val.(string)
	return s, ok
}

type execResult struct {
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int    `json:"exitCode"`
	Error    string `json:"error,omitempty"`
}

func marshalExecResult(stdout, stderr string, ps *os.ProcessState, err error) (*mcp.CallToolResult, error) {
	result := execResult{
		Stdout: stdout,
		Stderr: stderr,
	}
	if ps != nil {
		result.ExitCode = ps.ExitCode()
	}
	if err != nil {
		result.Error = err.Error()
		if result.ExitCode == 0 {
			result.ExitCode = -1
		}
	}
	return marshalToolResult(result)
}

func marshalToolResult(data interface{}) (*mcp.CallToolResult, error) {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal result: %v", err)), nil
	}
	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// buildShellArgs constructs command arguments for the given shell.
func buildShellArgs(shell, scriptOrCmd string, isInline bool, extraArgs []string) []string {
	shellLower := strings.ToLower(shell)
	switch shellLower {
	case "cmd":
		if isInline {
			return []string{"cmd", "/c", scriptOrCmd}
		}
		args := []string{"cmd", "/c", scriptOrCmd}
		args = append(args, extraArgs...)
		return args
	case "powershell", "pwsh":
		if isInline {
			return []string{shellLower, "-NoProfile", "-Command", scriptOrCmd}
		}
		args := []string{shellLower, "-NoProfile", "-File", scriptOrCmd}
		args = append(args, extraArgs...)
		return args
	default:
		// bash, sh, zsh
		if isInline {
			return []string{shellLower, "-c", scriptOrCmd}
		}
		args := []string{shellLower, scriptOrCmd}
		args = append(args, extraArgs...)
		return args
	}
}

 
