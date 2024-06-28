---
meta:
  - name: description
    content: "Fairwinds Pluto | Installation documentation"
---
# Installation

## asdf

We have an [asdf](https://asdf-vm.com/#/) plugin [here](https://github.com/FairwindsOps/asdf-pluto). You can install with:

```
asdf plugin-add pluto
asdf install pluto latest
asdf local pluto latest
```

## Binary

Install the binary from our [releases](https://github.com/FairwindsOps/pluto/releases) page.

## Homebrew Tap

```
brew install FairwindsOps/tap/pluto
```

## Scoop (Windows)
Note: This is not maintained by Fairwinds, but should stay up to date with future releases.

```
scoop install pluto
```

# Verify Artifacts

Fairwinds signs the Pluto docker image and the checksums file with [cosign](https://github.com/sigstore/cosign). Our public key is available at https://artifacts.fairwinds.com/cosign.pub

You can verify the checksums file from the [releases](https://github.com/FairwindsOps/pluto/releases) page with the following command:

```
cosign verify-blob checksums.txt --signature=checksums.txt.sig  --key https://artifacts.fairwinds.com/cosign.pub
```

Verifying docker images is even easier:

```
cosign verify us-docker.pkg.dev/fairwinds-ops/oss/pluto:v5 --key https://artifacts.fairwinds.com/cosign.pub
```

