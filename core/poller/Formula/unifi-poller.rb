# Homebrew Formula, still under development - June 2019
require "language/go"

# Classname should match the name of the installed package.
class UnifiPoller < Formula
  version "1.2.3"
  desc "This daemon polls a Unifi controller at a short interval and stores the collected metric data in an Influx Database."
  homepage "https://github.com/davidnewhall/unifi-poller"

  # Source code archive. Each tagged release will have one
  url "https://github.com/davidnewhall/unifi-poller/archive/v#{version}.tar.gz"
  sha256 "d536cb767b663a1c24410b27bd7e51a2b9a78820ed1ceeb2bb61e30e27235890"
  head "https://github.com/davidnewhall/unifi-poller"

  depends_on "go" => :build
  depends_on "dep"

  def install
    ENV["GOPATH"] = buildpath

    bin_path = buildpath/"src/github.com/davidnewhall/unifi-poller"
    # Copy all files from their current location (GOPATH root)
    # to $GOPATH/src/github.com/davidnewhall/unifi-poller
    bin_path.install Dir["*"]
    cd bin_path do
      # Install the compiled binary into Homebrew's `bin` - a pre-existing
      # global variable
      system "dep", "ensure"
      system "make", "install", "VERSION=#{version}", "PREFIX=#{prefix}"
    end
  end

  test do
    assert_match "unifi-poller v#{version}", shell_output("#{bin}/unifi-poller -v 2>&1", 2)
  end
end
