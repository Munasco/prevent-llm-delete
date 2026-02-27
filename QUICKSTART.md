# 🚀 Quick Start - prevent-llm-delete

Get protected from accidental deletions in 60 seconds!

## Installation

### Option 1: Homebrew (Recommended for macOS/Linux)

```bash
brew install yourusername/tap/prevent-llm-delete
```

### Option 2: Curl Script

```bash
curl -fsSL https://raw.githubusercontent.com/yourusername/prevent-llm-delete/main/install.sh | bash
```

### Option 3: Download Binary

1. Go to [Releases](https://github.com/yourusername/prevent-llm-delete/releases)
2. Download for your platform
3. Extract and move to `/usr/local/bin` (or `C:\Windows\System32` on Windows)

## Activate Protection

```bash
# Run the installer
prevent-llm-delete install

# Reload your shell
source ~/.zshrc  # bash/zsh
source ~/.config/fish/config.fish  # fish
# Or restart PowerShell on Windows
```

## Verify It Works

```bash
# Check status
prevent-llm-delete status

# Should show:
# ✅ Installed
# ✅ trash-cli Available (Unix only)
```

## Test It

```bash
# Create a test file
echo "test" > /tmp/test.txt

# Try to delete with dangerous flags
rm -rf /tmp/test.txt

# You should see:
# ⚠️  Dangerous flag '-rf' stripped by prevent-llm-delete

# Check trash
trash-list | grep test.txt  # Unix
# Check Recycle Bin on Windows
```

## You're Protected! 🛡️

Now when you (or an LLM) runs:
- `rm -rf important-files/` → Safely moved to trash
- `Remove-Item -Force file.txt` → Safely moved to Recycle Bin
- Dangerous flags are automatically stripped
- All deletions are recoverable

## Commands

```bash
prevent-llm-delete status      # Check if installed
prevent-llm-delete uninstall   # Remove protection
prevent-llm-delete help        # Show help
```

## Bypass When Needed

```bash
# Unix: use real rm for one command
command rm file.txt

# Windows: use full cmdlet name
Microsoft.PowerShell.Management\Remove-Item file.txt
```

## Uninstall

```bash
# Remove the wrapper (keeps binary)
prevent-llm-delete uninstall

# Remove the binary
brew uninstall prevent-llm-delete  # Homebrew
sudo rm /usr/local/bin/prevent-llm-delete  # Manual
```

## Need Help?

- Run: `prevent-llm-delete help`
- Read: [Full README](README.md)
- Report issues: [GitHub Issues](https://github.com/yourusername/prevent-llm-delete/issues)

---

**That's it!** You're now protected from accidental deletions. 🎉
