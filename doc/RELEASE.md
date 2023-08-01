# How to release

1. Release to GitHub Repo

When you push the tag, the action is automatically executed on GitHub and the release is registered.

2. Registering a new version in homebrew-tools

   a. Compress the executable macos binary with tar.

   b. Generate a SHA256 fingerprint.

   ```sh
   shasum -a 256 gqlmerge-0.2.6.tar.gz
   ```

   c. Update the brew formula.

   ```ruby
   class Gqlmerge < Formula
     desc "A merge and stitch tool for GraphQL schema"
     homepage "https://github.com/mununki/gqlmerge"
     # Path to the tar file you added to your GitHub Release
     url "https://github.com/mununki/gqlmerge/releases/download/v0.2.6/gqlmerge-0.2.6.tar.gz"
     # sha256 fingerprint
     sha256 "a6124d4804204f94c598d3d9e542e5c93d410be47e9df05aefc5ed1d2b3ae159"

     def install
       # Rename the binary
       bin.install "gqlmerge-macos" => "gqlmerge"
     end

     test do
       system "#{bin}/gqlmerge", "-v"
     end
   end
   ```

   d. Push to the homebrew-tools GitHub repo.
