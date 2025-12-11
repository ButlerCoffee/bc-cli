class BcCli < Formula
  desc "Butler Coffee CLI - Discover coffee knowledge and manage subscriptions from your terminal"
  homepage "https://github.com/ButlerCoffee/bc-cli"
  version "0.1.0"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/ButlerCoffee/bc-cli/releases/download/v0.1.0/bc-cli-darwin-arm64"
      sha256 "SHA256_ARM64_PLACEHOLDER" # ARM64
    else
      url "https://github.com/ButlerCoffee/bc-cli/releases/download/v0.1.0/bc-cli-darwin-amd64"
      sha256 "SHA256_AMD64_PLACEHOLDER" # AMD64
    end
  end

  def install
    bin.install "bc-cli-darwin-#{Hardware::CPU.arch}" => "bc-cli"
  end

  test do
    assert_match "bc-cli version", shell_output("#{bin}/bc-cli --version")
  end
end
