class PreventLlmDelete < Formula
  desc "Cross-platform safe deletion wrapper to prevent accidental deletions"
  homepage "https://github.com/yourusername/prevent-llm-delete"
  version "1.0.0"

  if OS.mac? && Hardware::CPU.arm?
    url "https://github.com/yourusername/prevent-llm-delete/releases/download/v1.0.0/prevent-llm-delete-darwin-arm64.tar.gz"
    sha256 "YOUR_ARM64_SHA256_HERE"  # Update after release
  elsif OS.mac?
    url "https://github.com/yourusername/prevent-llm-delete/releases/download/v1.0.0/prevent-llm-delete-darwin-amd64.tar.gz"
    sha256 "YOUR_AMD64_SHA256_HERE"  # Update after release
  elsif OS.linux? && Hardware::CPU.arm?
    url "https://github.com/yourusername/prevent-llm-delete/releases/download/v1.0.0/prevent-llm-delete-linux-arm64.tar.gz"
    sha256 "YOUR_LINUX_ARM64_SHA256_HERE"  # Update after release
  else
    url "https://github.com/yourusername/prevent-llm-delete/releases/download/v1.0.0/prevent-llm-delete-linux-amd64.tar.gz"
    sha256 "YOUR_LINUX_AMD64_SHA256_HERE"  # Update after release
  end

  depends_on "trash" => :recommended

  def install
    if OS.mac? && Hardware::CPU.arm?
      bin.install "prevent-llm-delete-darwin-arm64" => "prevent-llm-delete"
    elsif OS.mac?
      bin.install "prevent-llm-delete-darwin-amd64" => "prevent-llm-delete"
    elsif OS.linux? && Hardware::CPU.arm?
      bin.install "prevent-llm-delete-linux-arm64" => "prevent-llm-delete"
    else
      bin.install "prevent-llm-delete-linux-amd64" => "prevent-llm-delete"
    end
  end

  def caveats
    <<~EOS
      🔒 prevent-llm-delete is now installed!

      Get started:
        prevent-llm-delete install    # Install the protection
        prevent-llm-delete status     # Check status

      This will wrap your rm command to use trash instead of permanent deletion.
      Dangerous flags like -rf are automatically stripped.

      To uninstall the wrapper (not this formula):
        prevent-llm-delete uninstall

      Documentation: #{homepage}
    EOS
  end

  test do
    system "#{bin}/prevent-llm-delete", "version"
    system "#{bin}/prevent-llm-delete", "status"
  end
end
