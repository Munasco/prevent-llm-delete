# Installation Guide

Complete installation instructions for all platforms.

## Table of Contents

- [macOS](#macos)
- [Linux](#linux)
- [Windows](#windows)
- [From Source](#from-source)
- [Troubleshooting](#troubleshooting)

---

## macOS

### Method 1: Homebrew (Recommended)

```bash
# Add the tap (one-time)
brew tap yourusername/tap

# Install
brew install prevent-llm-delete

# Activate
prevent-llm-delete install
source ~/.zshrc
```

### Method 2: Curl Script

```bash
curl -fsSL https://raw.githubusercontent.com/yourusername/prevent-llm-delete/main/install.sh | bash
prevent-llm-delete install
source ~/.zshrc
```

### Method 3: Manual Download

```bash
# Intel Macs
curl -L https://github.com/yourusername/prevent-llm-delete/releases/download/v1.0.0/prevent-llm-delete-darwin-amd64.tar.gz -o prevent-llm-delete.tar.gz

# Apple Silicon (M1/M2/M3)
curl -L https://github.com/yourusername/prevent-llm-delete/releases/download/v1.0.0/prevent-llm-delete-darwin-arm64.tar.gz -o prevent-llm-delete.tar.gz

# Extract and install
tar -xzf prevent-llm-delete.tar.gz
sudo mv prevent-llm-delete-* /usr/local/bin/prevent-llm-delete
chmod +x /usr/local/bin/prevent-llm-delete

# Activate
prevent-llm-delete install
source ~/.zshrc
```

---

## Linux

### Method 1: Curl Script

```bash
curl -fsSL https://raw.githubusercontent.com/yourusername/prevent-llm-delete/main/install.sh | bash
prevent-llm-delete install
source ~/.bashrc  # or ~/.zshrc if using zsh
```

### Method 2: Manual Download

```bash
# x86_64
curl -L https://github.com/yourusername/prevent-llm-delete/releases/download/v1.0.0/prevent-llm-delete-linux-amd64.tar.gz -o prevent-llm-delete.tar.gz

# ARM64
curl -L https://github.com/yourusername/prevent-llm-delete/releases/download/v1.0.0/prevent-llm-delete-linux-arm64.tar.gz -o prevent-llm-delete.tar.gz

# Extract and install
tar -xzf prevent-llm-delete.tar.gz
sudo mv prevent-llm-delete-* /usr/local/bin/prevent-llm-delete
chmod +x /usr/local/bin/prevent-llm-delete

# Activate
prevent-llm-delete install
source ~/.bashrc
```

### Method 3: From Package (Future)

```bash
# Debian/Ubuntu (coming soon)
sudo dpkg -i prevent-llm-delete_1.0.0_amd64.deb

# Fedora/RHEL (coming soon)
sudo rpm -i prevent-llm-delete-1.0.0.x86_64.rpm
```

---

## Windows

### Method 1: Download Binary

1. Download [prevent-llm-delete-windows-amd64.zip](https://github.com/yourusername/prevent-llm-delete/releases/download/v1.0.0/prevent-llm-delete-windows-amd64.zip)

2. Extract the ZIP file

3. **Run PowerShell as Administrator**

4. Install:
   ```powershell
   Move-Item prevent-llm-delete.exe C:\Windows\System32\
   ```

5. Activate:
   ```powershell
   prevent-llm-delete install
   ```

6. Restart PowerShell

### Method 2: PowerShell Script (Future)

```powershell
# Coming soon
iwr -useb https://raw.githubusercontent.com/yourusername/prevent-llm-delete/main/install.ps1 | iex
```

### Method 3: Chocolatey (Future)

```powershell
# Coming soon
choco install prevent-llm-delete
```

---

## From Source

### Prerequisites

- Go 1.21 or later
- Make (optional, but recommended)

### Build and Install

```bash
# Clone the repository
git clone https://github.com/yourusername/prevent-llm-delete.git
cd prevent-llm-delete

# Build
make build

# Install
make install

# Or manually
sudo mv prevent-llm-delete /usr/local/bin/

# Activate
prevent-llm-delete install
source ~/.zshrc  # or appropriate shell config
```

---

## Verification

After installation, verify it works:

```bash
# Check version
prevent-llm-delete version

# Check status
prevent-llm-delete status

# Should show:
# ✅ Installed
# ✅ trash-cli Available (Unix)
```

## Test It

```bash
# Create a test file
echo "test" > /tmp/prevent-llm-delete-test.txt

# Try to delete with dangerous flags
rm -rf /tmp/prevent-llm-delete-test.txt

# You should see a warning:
# ⚠️  Dangerous flag '-rf' stripped by prevent-llm-delete

# Verify file is in trash (Unix)
trash-list | grep prevent-llm-delete-test

# Restore if needed
trash-restore
```

---

## Troubleshooting

### "Command not found"

**Problem:** `prevent-llm-delete: command not found`

**Solution:**
```bash
# Check if binary exists
ls -la /usr/local/bin/prevent-llm-delete

# If missing, reinstall
# If exists, check PATH
echo $PATH | grep /usr/local/bin

# Add to PATH if needed
export PATH="/usr/local/bin:$PATH"
```

### "trash: command not found" (Unix)

**Problem:** trash-cli is not installed

**Solution:**
```bash
# macOS
brew install trash

# Linux (Debian/Ubuntu)
sudo apt install trash-cli

# Linux (Fedora/RHEL)
sudo dnf install trash-cli

# Or via npm
npm install -g trash-cli
```

### "Permission denied"

**Problem:** Cannot write to `/usr/local/bin`

**Solution:**
```bash
# Install with sudo
sudo mv prevent-llm-delete /usr/local/bin/
sudo chmod +x /usr/local/bin/prevent-llm-delete
```

### "rm still deletes permanently"

**Problem:** The wrapper isn't active

**Solution:**
```bash
# Check if installed
prevent-llm-delete status

# Reload shell config
source ~/.zshrc  # bash/zsh
source ~/.config/fish/config.fish  # fish

# Or restart your terminal
```

### Windows: "Script execution is disabled"

**Problem:** PowerShell execution policy blocks scripts

**Solution:**
```powershell
# Run as Administrator
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser

# Then try again
prevent-llm-delete install
```

### "Existing rm override detected"

**Problem:** You already have a custom rm function

**Solution:**
This is intentional! The tool won't overwrite your existing protection.

Options:
1. Keep your existing override (if it already protects you)
2. Manually remove it from your shell config and reinstall:
   ```bash
   # Edit your config
   nano ~/.zshrc  # or appropriate config

   # Remove the existing rm function
   # Then reinstall
   prevent-llm-delete install
   ```

---

## Next Steps

- Read the [Quick Start Guide](QUICKSTART.md)
- Check the [Full Documentation](README.md)
- Report issues on [GitHub](https://github.com/yourusername/prevent-llm-delete/issues)

## Uninstall

See [Uninstallation Guide](README.md#uninstallation)

---

Need more help? [Open an issue](https://github.com/yourusername/prevent-llm-delete/issues/new)
