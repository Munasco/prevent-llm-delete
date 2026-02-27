package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const VERSION = "1.0.0"

var (
	MARKER_START    = "# === prevent-llm-delete START ==="
	MARKER_END      = "# === prevent-llm-delete END ==="
	PS_MARKER_START = "# === prevent-llm-delete START ==="
	PS_MARKER_END   = "# === prevent-llm-delete END ==="
)

const UNIX_RM_OVERRIDE = `
# prevent-llm-delete: Safe rm wrapper that uses trash instead
rm() {
  local args=()
  local parsing_flags=1

  for arg in "$@"; do
    if [ "$parsing_flags" -eq 1 ]; then
      case "$arg" in
        --)
          parsing_flags=0
          continue
          ;;
        -r|--recursive|-f|--force|-rf|-fr|-Rf|-RF|-rF)
          echo "⚠️  Dangerous flag '$arg' stripped by prevent-llm-delete"
          continue
          ;;
      esac
    fi
    args+=("$arg")
  done

  if [ ${#args[@]} -eq 0 ]; then
    echo "rm: no files specified"
    return 1
  fi

  command trash "${args[@]}"
}
`

const FISH_RM_OVERRIDE = `
# prevent-llm-delete: Safe rm wrapper that uses trash instead
function rm
    set -l args
    set -l parsing_flags 1

    for arg in $argv
        if test $parsing_flags -eq 1
            switch $arg
                case '--'
                    set parsing_flags 0
                    continue
                case '-r' '--recursive' '-f' '--force' '-rf' '-fr' '-Rf' '-RF' '-rF'
                    echo "⚠️  Dangerous flag '$arg' stripped by prevent-llm-delete"
                    continue
            end
        end
        set args $args $arg
    end

    if test (count $args) -eq 0
        echo "rm: no files specified"
        return 1
    end

    command trash $args
end
`

const POWERSHELL_RM_OVERRIDE = `
# prevent-llm-delete: Safe Remove-Item wrapper
function Remove-Item {
    [CmdletBinding()]
    param(
        [Parameter(ValueFromPipeline=$true, ValueFromPipelineByPropertyName=$true, Position=0)]
        [string[]]$Path,

        [switch]$Recurse,
        [switch]$Force,
        [Parameter(ValueFromRemainingArguments=$true)]
        $RemainingArgs
    )

    process {
        if ($Recurse -or $Force) {
            Write-Warning "⚠️  Dangerous flags (-Recurse/-Force) stripped by prevent-llm-delete"
        }

        if (-not $Path) {
            Write-Error "Remove-Item: no files specified"
            return
        }

        foreach ($item in $Path) {
            if (Get-Command trash -ErrorAction SilentlyContinue) {
                trash $item
            } else {
                # Fallback to PowerShell Recycle Bin
                $shell = New-Object -ComObject Shell.Application
                $item_obj = Get-Item $item -ErrorAction SilentlyContinue
                if ($item_obj) {
                    $shell.Namespace(0).ParseName($item_obj.FullName).InvokeVerb("delete")
                }
            }
        }
    }
}

Set-Alias -Name rm -Value Remove-Item -Option AllScope -Force
Set-Alias -Name del -Value Remove-Item -Option AllScope -Force
`

type ShellConfig struct {
	Name       string
	ConfigPath string
	Override   string
}

func getHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("❌ Could not determine home directory")
		os.Exit(1)
	}
	return home
}

func detectShell() ShellConfig {
	isWindows := runtime.GOOS == "windows"
	home := getHomeDir()

	if isWindows {
		psProfile := filepath.Join(home, "Documents", "PowerShell", "Microsoft.PowerShell_profile.ps1")
		return ShellConfig{
			Name:       "powershell",
			ConfigPath: psProfile,
			Override:   POWERSHELL_RM_OVERRIDE,
		}
	}

	// Check SHELL environment variable
	shell := os.Getenv("SHELL")

	if strings.Contains(shell, "fish") {
		return ShellConfig{
			Name:       "fish",
			ConfigPath: filepath.Join(home, ".config", "fish", "config.fish"),
			Override:   FISH_RM_OVERRIDE,
		}
	}

	if strings.Contains(shell, "zsh") {
		return ShellConfig{
			Name:       "zsh",
			ConfigPath: filepath.Join(home, ".zshrc"),
			Override:   UNIX_RM_OVERRIDE,
		}
	}

	// Default to bash
	return ShellConfig{
		Name:       "bash",
		ConfigPath: filepath.Join(home, ".bashrc"),
		Override:   UNIX_RM_OVERRIDE,
	}
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func checkTrashInstalled() bool {
	if runtime.GOOS == "windows" {
		return true // Windows uses built-in Recycle Bin
	}
	return commandExists("trash")
}

func installTrash() {
	fmt.Println("📦 Installing trash utility...")

	if runtime.GOOS == "windows" {
		fmt.Println("ℹ️  Windows uses built-in Recycle Bin, no installation needed")
		return
	}

	if runtime.GOOS == "darwin" {
		// Try Homebrew first
		if commandExists("brew") {
			cmd := exec.Command("brew", "install", "trash")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Println("⚠️  Homebrew install failed, trying npm...")
				installTrashNpm()
			} else {
				fmt.Println("✅ trash installed via Homebrew")
			}
		} else {
			installTrashNpm()
		}
	} else {
		// Linux - try apt first, then npm
		if commandExists("apt-get") {
			cmd := exec.Command("sudo", "apt-get", "update")
			cmd.Run()
			cmd = exec.Command("sudo", "apt-get", "install", "-y", "trash-cli")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Println("⚠️  apt install failed, trying npm...")
				installTrashNpm()
			} else {
				fmt.Println("✅ trash-cli installed via apt")
			}
		} else {
			installTrashNpm()
		}
	}
}

func installTrashNpm() {
	cmd := exec.Command("npm", "install", "-g", "trash-cli")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("❌ Failed to install trash-cli via npm")
		fmt.Println("")
		fmt.Println("Please install manually:")
		if runtime.GOOS == "darwin" {
			fmt.Println("  macOS: brew install trash")
		} else {
			fmt.Println("  Linux: sudo apt install trash-cli")
		}
		os.Exit(1)
	}
	fmt.Println("✅ trash-cli installed via npm")
}

func checkExistingOverride(content, shell string) bool {
	if shell == "powershell" {
		patterns := []string{
			"function Remove-Item",
			"function rm",
			"Set-Alias.*rm",
		}
		for _, pattern := range patterns {
			if strings.Contains(strings.ToLower(content), strings.ToLower(pattern)) {
				return true
			}
		}
	} else if shell == "fish" {
		return strings.Contains(content, "function rm")
	} else {
		patterns := []string{
			"rm()",
			"function rm",
			"alias rm=",
		}
		for _, pattern := range patterns {
			if strings.Contains(content, pattern) {
				return true
			}
		}
	}
	return false
}

func install() {
	fmt.Println("🔒 Installing prevent-llm-delete...")

	// Check trash first
	if !checkTrashInstalled() {
		fmt.Println("⚠️  trash utility not found")
		fmt.Print("Would you like to install it now? (y/n): ")
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))

		if answer == "y" || answer == "yes" {
			installTrash()
		} else {
			fmt.Println("❌ Installation cancelled")
			os.Exit(1)
		}
	}

	shellConfig := detectShell()
	fmt.Printf("📝 Detected shell: %s\n", shellConfig.Name)
	fmt.Printf("📂 Config file: %s\n", shellConfig.ConfigPath)

	// Ensure config directory exists
	configDir := filepath.Dir(shellConfig.ConfigPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Printf("❌ Failed to create config directory: %v\n", err)
		os.Exit(1)
	}

	// Read existing config
	var configContent string
	if data, err := os.ReadFile(shellConfig.ConfigPath); err == nil {
		configContent = string(data)
	}

	// Check if already installed
	markerStart := MARKER_START
	if shellConfig.Name == "powershell" {
		markerStart = PS_MARKER_START
	}

	if strings.Contains(configContent, markerStart) {
		fmt.Println("⚠️  prevent-llm-delete is already installed")
		fmt.Println("   Run `prevent-llm-delete uninstall` to remove it first")
		os.Exit(0)
	}

	// Check for existing rm override
	if checkExistingOverride(configContent, shellConfig.Name) {
		fmt.Println("⚠️  Warning: An existing rm override/alias/function was detected")
		fmt.Println("")
		fmt.Println("Your shell config already has a custom rm definition.")
		fmt.Println("Installing prevent-llm-delete will NOT override it.")
		fmt.Println("")
		fmt.Println("Options:")
		fmt.Println("  1. Remove your existing rm override manually and run install again")
		fmt.Println("  2. Keep your existing override (it may already provide similar protection)")
		fmt.Println("")
		fmt.Print("Continue anyway? This will NOT affect your existing override. (y/n): ")

		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))

		if answer == "y" || answer == "yes" {
			fmt.Println("")
			fmt.Println("ℹ️  Skipping installation since an rm override already exists.")
			fmt.Println("   Your existing override will continue to work.")
			os.Exit(0)
		} else {
			fmt.Println("❌ Installation cancelled")
			os.Exit(1)
		}
	}

	// Backup existing config
	if _, err := os.Stat(shellConfig.ConfigPath); err == nil {
		backupPath := shellConfig.ConfigPath + ".backup"
		if err := os.WriteFile(backupPath, []byte(configContent), 0644); err != nil {
			fmt.Printf("⚠️  Failed to create backup: %v\n", err)
		} else {
			fmt.Printf("📋 Backed up %s to %s\n", filepath.Base(shellConfig.ConfigPath), filepath.Base(backupPath))
		}
	}

	// Add the override
	markerEnd := MARKER_END
	if shellConfig.Name == "powershell" {
		markerEnd = PS_MARKER_END
	}

	updatedContent := configContent + "\n" + markerStart + shellConfig.Override + markerEnd + "\n"

	if err := os.WriteFile(shellConfig.ConfigPath, []byte(updatedContent), 0644); err != nil {
		fmt.Printf("❌ Failed to write config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ prevent-llm-delete installed successfully!")
	fmt.Println("")

	// Platform-specific activation instructions
	if runtime.GOOS == "windows" {
		fmt.Println("🔄 Restart PowerShell to activate")
		fmt.Println("")
		fmt.Println("ℹ️  Usage:")
		fmt.Println("   - rm file.txt              → moves to Recycle Bin (safe)")
		fmt.Println("   - Remove-Item -Recurse dir → flags stripped, uses Recycle Bin")
	} else {
		sourceCmd := fmt.Sprintf("source %s", shellConfig.ConfigPath)
		if shellConfig.Name == "fish" {
			sourceCmd = "source ~/.config/fish/config.fish"
		}
		fmt.Printf("🔄 Run `%s` or restart your terminal to activate\n", sourceCmd)
		fmt.Println("")
		fmt.Println("ℹ️  Usage:")
		fmt.Println("   - rm file.txt  → moves to trash (safe)")
		fmt.Println("   - rm -rf dir   → dangerous flags stripped, uses trash")
		fmt.Println("   - command rm   → bypass wrapper to use real rm")
		if shellConfig.Name != "fish" {
			fmt.Println("   - unset -f rm  → temporarily disable wrapper")
		}
	}
}

func uninstall() {
	fmt.Println("🔓 Uninstalling prevent-llm-delete...")

	shellConfig := detectShell()

	if _, err := os.Stat(shellConfig.ConfigPath); os.IsNotExist(err) {
		fmt.Printf("❌ %s not found\n", filepath.Base(shellConfig.ConfigPath))
		os.Exit(1)
	}

	data, err := os.ReadFile(shellConfig.ConfigPath)
	if err != nil {
		fmt.Printf("❌ Failed to read config: %v\n", err)
		os.Exit(1)
	}
	configContent := string(data)

	// Check if installed
	markerStart := MARKER_START
	markerEnd := MARKER_END
	if shellConfig.Name == "powershell" {
		markerStart = PS_MARKER_START
		markerEnd = PS_MARKER_END
	}

	if !strings.Contains(configContent, markerStart) {
		fmt.Println("⚠️  prevent-llm-delete is not installed")
		os.Exit(0)
	}

	// Remove the function
	startIdx := strings.Index(configContent, markerStart)
	endIdx := strings.Index(configContent, markerEnd)

	if startIdx == -1 || endIdx == -1 {
		fmt.Printf("❌ Error: Could not find markers in %s\n", filepath.Base(shellConfig.ConfigPath))
		os.Exit(1)
	}

	before := configContent[:startIdx]
	after := configContent[endIdx+len(markerEnd):]
	updatedContent := before + after

	// Backup before uninstall
	backupPath := shellConfig.ConfigPath + ".backup"
	os.WriteFile(backupPath, data, 0644)
	fmt.Printf("📋 Backed up %s to %s\n", filepath.Base(shellConfig.ConfigPath), filepath.Base(backupPath))

	// Write updated config
	if err := os.WriteFile(shellConfig.ConfigPath, []byte(updatedContent), 0644); err != nil {
		fmt.Printf("❌ Failed to write config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ prevent-llm-delete uninstalled successfully!")
	fmt.Println("")

	if runtime.GOOS == "windows" {
		fmt.Println("🔄 Restart PowerShell to apply changes")
	} else {
		sourceCmd := fmt.Sprintf("source %s", shellConfig.ConfigPath)
		if shellConfig.Name == "fish" {
			sourceCmd = "source ~/.config/fish/config.fish"
		}
		fmt.Printf("🔄 Run `%s` or restart your terminal to apply changes\n", sourceCmd)
	}
}

func status() {
	shellConfig := detectShell()

	fmt.Println("📊 prevent-llm-delete status:")
	fmt.Println("")
	fmt.Printf("   Platform:     %s\n", runtime.GOOS)
	fmt.Printf("   Shell:        %s\n", shellConfig.Name)
	fmt.Printf("   Config:       %s\n", shellConfig.ConfigPath)
	fmt.Println("")

	if _, err := os.Stat(shellConfig.ConfigPath); os.IsNotExist(err) {
		fmt.Println("   Installation: ❌ Config file not found")
		return
	}

	data, err := os.ReadFile(shellConfig.ConfigPath)
	if err != nil {
		fmt.Println("   Installation: ❌ Could not read config")
		return
	}
	configContent := string(data)

	markerStart := MARKER_START
	if shellConfig.Name == "powershell" {
		markerStart = PS_MARKER_START
	}

	isInstalled := strings.Contains(configContent, markerStart)
	trashInstalled := checkTrashInstalled()

	if isInstalled {
		fmt.Println("   Installation: ✅ Installed")
	} else {
		fmt.Println("   Installation: ❌ Not installed")
	}

	if runtime.GOOS != "windows" {
		if trashInstalled {
			fmt.Println("   trash-cli:    ✅ Available")
		} else {
			fmt.Println("   trash-cli:    ❌ Not found")
		}
	}

	fmt.Println("")

	if isInstalled {
		if runtime.GOOS == "windows" {
			fmt.Println("ℹ️  Remove-Item is wrapped to use Recycle Bin")
			fmt.Println("   - Dangerous flags (-Recurse, -Force) are automatically stripped")
		} else {
			fmt.Println("ℹ️  The rm command is wrapped to use trash")
			fmt.Println("   - Dangerous flags (-r, -f, -rf) are automatically stripped")
			fmt.Println("   - Use `command rm` to access the real rm if needed")
		}
	}
}

func showHelp() {
	platform := runtime.GOOS
	if platform == "darwin" {
		platform = "macOS"
	} else if platform == "windows" {
		platform = "Windows"
	} else {
		platform = "Linux"
	}

	fmt.Printf(`
🔒 prevent-llm-delete v%s - Cross-platform safe deletion wrapper

Usage:
  prevent-llm-delete install    Install the safe deletion wrapper
  prevent-llm-delete uninstall  Remove the wrapper
  prevent-llm-delete status     Check installation status
  prevent-llm-delete version    Show version
  prevent-llm-delete help       Show this help message

What it does:
  - Overrides deletion commands to use trash/recycle bin
  - Strips dangerous flags to prevent accidental permanent deletions
  - Protects you from destructive commands run by LLMs or humans
  - Files can be recovered instead of being permanently deleted

Supported Platforms:
  ✅ Windows   - PowerShell (uses Recycle Bin)
  ✅ macOS     - bash/zsh/fish (uses trash-cli)
  ✅ Linux     - bash/zsh/fish (uses trash-cli)

Current Platform: %s

`, VERSION, platform)

	if runtime.GOOS == "windows" {
		fmt.Println(`Windows Examples:
  Remove-Item file.txt              → Moves to Recycle Bin
  rm file.txt                       → Moves to Recycle Bin
  Remove-Item -Recurse -Force dir   → Strips flags, uses Recycle Bin`)
	} else {
		fmt.Println(`Unix Examples:
  rm file.txt        → Safely moves to trash
  rm -rf directory   → Strips -rf flag, moves to trash
  command rm file    → Bypass wrapper (use real rm)
  unset -f rm        → Temporarily disable wrapper`)
	}

	fmt.Println(`
Requirements:`)
	if runtime.GOOS == "windows" {
		fmt.Println(`  - PowerShell (built-in on Windows)
  - Uses Windows Recycle Bin (no extra tools needed)`)
	} else {
		fmt.Println(`  - trash-cli (auto-installed if missing)
    macOS: brew install trash
    Linux: sudo apt install trash-cli`)
	}

	fmt.Println(`
Why This Exists:
  LLMs (Claude, ChatGPT, etc.) sometimes run destructive commands
  when they shouldn't. This tool prevents those accidents.

Safety Features:
  ✅ Automatic stripping of dangerous flags
  ✅ Recoverable deletions (trash/recycle bin)
  ✅ Clear warnings when flags are removed
  ✅ Easy bypass for advanced users
  ✅ Works across Windows, macOS, and Linux
`)
}

func main() {
	if len(os.Args) < 2 {
		showHelp()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "install":
		install()
	case "uninstall":
		uninstall()
	case "status":
		status()
	case "version", "--version", "-v":
		fmt.Printf("prevent-llm-delete v%s\n", VERSION)
	case "help", "--help", "-h":
		showHelp()
	default:
		fmt.Printf("❌ Unknown command: %s\n\n", command)
		showHelp()
		os.Exit(1)
	}
}
