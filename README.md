# 🔒 prevent-llm-delete

A cross-platform CLI tool that wraps deletion commands to prevent accidental permanent deletions by LLMs and humans alike.

**No dependencies. Single binary. Works everywhere.**

## Why?

LLMs (Claude, ChatGPT, etc.) sometimes run destructive commands like `rm -rf` when they shouldn't. This tool prevents those accidents by replacing deletion commands with safe wrappers that use trash/recycle bin instead.

## Features

✅ **Zero dependencies** - Single compiled binary, no runtime required
✅ **Cross-platform** - Windows, macOS, Linux
✅ **Multi-shell** - bash, zsh, fish, PowerShell
✅ **Smart detection** - Won't overwrite existing rm overrides
✅ **Auto flag stripping** - Removes `-r`, `-f`, `-rf`, `-Recurse`, `-Force`
✅ **Recoverable** - Uses trash-cli (Unix) or Recycle Bin (Windows)
✅ **Safe installation** - Backs up your config before modifying

## Installation

### 📦 Homebrew (macOS/Linux)

```bash
brew install yourusername/tap/prevent-llm-delete
prevent-llm-delete install
source ~/.zshrc  # or restart terminal
```

### 🌐 Curl Install (macOS/Linux)

```bash
curl -fsSL https://raw.githubusercontent.com/yourusername/prevent-llm-delete/main/install.sh | bash
prevent-llm-delete install
source ~/.zshrc  # or restart terminal
```

### 📥 Download Binary

1. Download for your platform from [Releases](https://github.com/yourusername/prevent-llm-delete/releases):
   - **macOS (Intel):** `prevent-llm-delete-darwin-amd64.tar.gz`
   - **macOS (Apple Silicon):** `prevent-llm-delete-darwin-arm64.tar.gz`
   - **Linux (x64):** `prevent-llm-delete-linux-amd64.tar.gz`
   - **Linux (ARM):** `prevent-llm-delete-linux-arm64.tar.gz`
   - **Windows:** `prevent-llm-delete-windows-amd64.zip`

2. Extract and install:
   ```bash
   # macOS/Linux
   tar -xzf prevent-llm-delete-*.tar.gz
   sudo mv prevent-llm-delete-* /usr/local/bin/prevent-llm-delete
   chmod +x /usr/local/bin/prevent-llm-delete

   # Windows (PowerShell as Administrator)
   Expand-Archive prevent-llm-delete-windows-amd64.zip
   Move-Item prevent-llm-delete.exe C:\Windows\System32\
   ```

3. Run the installer:
   ```bash
   prevent-llm-delete install
   ```

### 🛠️ Build from Source

```bash
# Clone the repo
git clone https://github.com/yourusername/prevent-llm-delete.git
cd prevent-llm-delete

# Build
make build

# Install
make install

# Or build for all platforms
make build-all
```

## Usage

### Quick Start

```bash
# Install the wrapper
prevent-llm-delete install

# Check status
prevent-llm-delete status

# Reload your shell
source ~/.zshrc  # bash/zsh
source ~/.config/fish/config.fish  # fish
# Or restart PowerShell on Windows
```

### How It Works

**Before:**
```bash
rm -rf important-files/
# ❌ Files gone forever!
```

**After:**
```bash
rm -rf important-files/
# ⚠️  Dangerous flag '-rf' stripped by prevent-llm-delete
# ✅ Files safely moved to trash - recoverable!
```

### Platform-Specific Examples

#### Unix/macOS (bash/zsh/fish)

```bash
# Safe deletions (files go to trash)
rm file.txt
rm *.log
rm -rf dangerous-dir  # -rf is stripped, uses trash

# Bypass if you really need real rm
command rm file.txt       # One-time bypass
unset -f rm               # Disable for session
source ~/.zshrc           # Re-enable
```

#### Windows (PowerShell)

```powershell
# Safe deletions (files go to Recycle Bin)
Remove-Item file.txt
rm file.txt
Remove-Item -Recurse -Force dir  # Flags stripped

# Bypass if needed
Microsoft.PowerShell.Management\Remove-Item file.txt
```

## Commands

```bash
prevent-llm-delete install      # Install the wrapper
prevent-llm-delete uninstall    # Remove the wrapper
prevent-llm-delete status       # Check installation status
prevent-llm-delete version      # Show version
prevent-llm-delete help         # Show help
```

## How It Works

The tool adds a function to your shell config that:

1. **Intercepts** `rm` commands (or `Remove-Item` on Windows)
2. **Strips** dangerous flags like `-r`, `-f`, `-rf`
3. **Redirects** to `trash` (or Windows Recycle Bin)
4. **Shows warnings** when flags are removed

### Unix Shells (bash/zsh)

```bash
rm() {
  # Strip dangerous flags
  # Use trash instead of rm
  command trash "${args[@]}"
}
```

### PowerShell (Windows)

```powershell
function Remove-Item {
  # Strip -Recurse and -Force
  # Use Windows Recycle Bin COM API
}
```

## Configuration Files Modified

| Platform | Shell | Config File |
|----------|-------|-------------|
| macOS/Linux | bash | `~/.bashrc` |
| macOS/Linux | zsh | `~/.zshrc` |
| macOS/Linux | fish | `~/.config/fish/config.fish` |
| Windows | PowerShell | `~/Documents/PowerShell/Microsoft.PowerShell_profile.ps1` |

**Note:** Original config is backed up to `<config>.backup` before modification.

## Safety Features

✅ **Detects existing overrides** - Won't clobber your custom rm functions
✅ **Backs up configs** - Creates `.backup` files before changes
✅ **Clear warnings** - Shows when dangerous flags are stripped
✅ **Easy bypass** - Advanced users can access real `rm` when needed
✅ **Recoverable deletions** - Files go to trash/recycle bin
✅ **Cross-platform** - Works on Windows, macOS, Linux
✅ **No runtime** - Single compiled binary with no dependencies

## Status Check

```bash
$ prevent-llm-delete status

📊 prevent-llm-delete status:

   Platform:     darwin
   Shell:        zsh
   Config:       /Users/you/.zshrc

   Installation: ✅ Installed
   trash-cli:    ✅ Available

ℹ️  The rm command is wrapped to use trash
   - Dangerous flags (-r, -f, -rf) are automatically stripped
   - Use `command rm` to access the real rm if needed
```

## Requirements

### macOS/Linux
- **trash-cli** (auto-installed during setup)
  - macOS: `brew install trash`
  - Linux: `sudo apt install trash-cli`

### Windows
- **PowerShell** (built-in)
- Uses Windows Recycle Bin (no extra tools needed)

## Uninstallation

```bash
# Remove the wrapper (keeps the binary)
prevent-llm-delete uninstall

# Remove the binary
# macOS/Linux
sudo rm /usr/local/bin/prevent-llm-delete

# Homebrew
brew uninstall prevent-llm-delete

# Windows
Remove-Item C:\Windows\System32\prevent-llm-delete.exe
```

## Existing rm Overrides

If you already have a custom `rm` function, the tool will detect it and warn you:

```
⚠️  Warning: An existing rm override/alias/function was detected

Your shell config already has a custom rm definition.
Installing prevent-llm-delete will NOT override it.

Options:
  1. Remove your existing rm override manually and run install again
  2. Keep your existing override (it may already provide similar protection)
```

## Why This Matters

When working with AI coding assistants:

1. **LLMs make mistakes** - They might run `rm -rf` unnecessarily
2. **Context is limited** - AI doesn't always know what's important
3. **Humans make mistakes too** - This protects you from accidents
4. **Recovery is key** - Trash is recoverable, `rm -rf` is not
5. **Peace of mind** - Code with confidence knowing you have a safety net

## Development

```bash
# Build
make build

# Build for all platforms
make build-all

# Run tests
make test

# Create release archives
make release

# Clean
make clean
```

## Contributing

Contributions welcome! Ideas:

- [ ] Whitelist/blacklist for specific paths
- [ ] Dry-run mode to see what would be deleted
- [ ] GUI for managing trashed files
- [ ] Config file for custom behavior
- [ ] Integration tests

## License

MIT

## Related Projects

- [trash-cli](https://github.com/andreafrancia/trash-cli) - Cross-platform trash utility
- [safe-rm](https://github.com/lagerspetz/safe-rm) - Wrapper to prevent accidental deletions
- [trashy](https://github.com/oberblastmeister/trashy) - Alternative trash CLI in Rust

---

**Remember:** Even with this protection, always review what LLMs are doing before approving operations!
