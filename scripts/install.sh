#!/usr/bin/env bash
# This script is used to download and verify the binary.
# It expects the defaults based on the .goreleaser.yml

set -e

# Modify these two values based on org repo
OWNER=catpaladin
REPO=net-tools

if [[ -z "$VERSION" ]]; then
  VERSION="latest"
else
  echo "version set to: $VERSION"
  VERSION="tags/v$VERSION"
fi

OS_TYPE=$(uname | tr 'A-Z' 'a-z')
case ${OS_TYPE} in
  darwin)
    # arch set to all for universal binary
    ARCH="all"
    ;;
  linux)
    ARCH_CHECK=$(uname -m)
    if [[ "$ARCH_CHECK" == "x86_64" ]]; then
      ARCH="amd64"
    elif [[ "$ARCH_CHECK" == "aarch64" ]]; then
      ARCH="arm64"
    else
      printf >&2 '*** error: ARCH must be one of amd64|arm64'; exit 1
    fi
    ;;
  *) printf >&2 '*** error: OS_TYPE must be one of darwin|linux'; exit 1;;
esac

if ! hash jq 2> /dev/null; then
  echo "error: you do not have 'jq' installed which is required for this script."
  exit 1
fi

if ! hash curl 2> /dev/null; then
  echo "error: you do not have 'curl' installed which is required for this script."
  exit 1
fi

INSTALL_PATH="$HOME/.local/bin"
DOWNLOAD_PATH=$(mktemp -d -t $REPO.XXXXXXXX)
HASH_ARTIFACT="$DOWNLOAD_PATH/checksums.txt"
ARTIFACT="$DOWNLOAD_PATH/$REPO.tar.gz"

# trap just in case and cleanup tmp
cleanup() {
  local code=$?
  set +e
  trap - EXIT
  rm -rf "$DOWNLOAD_PATH"
  exit $code
}
trap cleanup INT EXIT

# curl and store metadata associated with the version release
METADATA=$(\
  curl -s "https://api.github.com/repos/$OWNER/$REPO/releases/$VERSION" \
  | jq -r ".assets[]")

# get urls for checksums and artifact
CHECKSUMS_DOWNLOAD_URL=$(echo $METADATA | jq -r ". | select(.name | contains(\"checksums.txt\")) | {browser_download_url} | .browser_download_url")
ARTIFACT_DOWNLOAD_URL=$(echo $METADATA | jq -r ". | select(.name | contains(\"${OS_TYPE}_${ARCH}\")) | {browser_download_url} | .browser_download_url")

# download files
curl -SfL --progress-bar "$CHECKSUMS_DOWNLOAD_URL" \
  -o "$HASH_ARTIFACT"

curl -SfL --progress-bar "$ARTIFACT_DOWNLOAD_URL" \
  -o "$ARTIFACT"

# determine if sha356sum or shasum on host
# to be used in verify_binary
compute_sha256sum() {
  cmd=$(which sha256sum shasum | head -n 1)
  case $(basename "$cmd") in
    sha256sum)
      sha256sum "$1" | cut -f 1 -d ' '
      ;;
    shasum)
      shasum -a 256 "$1" | cut -f 1 -d ' '
      ;;
    *) printf >&2 '*** error: Can not find sha256sum or shasum to compute checksum'; exit 1;;
  esac
}

# Verify downloaded binary hash
verify_binary() {
  local binary_path=$1
  local expected_hash=$2

  HASH_BIN=$(compute_sha256sum $binary_path)
  HASH_BIN=${HASH_BIN%%[[:blank:]]*}
  if [[ "$expected_hash" != "${HASH_BIN}" ]]; then
      printf >&2 "*** error: Download sha256 does not match $expected_hash, got ${HASH_BIN}"
      exit 1
  fi
  echo "binary verified successfully"
}

EXPECTED_HASH=$(grep -i "${OS_TYPE}_${ARCH}.tar.gz" "$HASH_ARTIFACT")
EXPECTED_HASH=${EXPECTED_HASH%%[[:blank:]]*}
verify_binary "$ARTIFACT" "$EXPECTED_HASH"

# Install
mkdir -p "$INSTALL_PATH"
chmod 755 "$ARTIFACT"
tar -xzof "$ARTIFACT" -C "$DOWNLOAD_PATH"
mv "$DOWNLOAD_PATH/$REPO" "$INSTALL_PATH"

echo "$REPO installed into \$HOME/.local/bin"
echo "please add \$HOME/.local/bin to your \$PATH or run \"export PATH=\$HOME/.local/bin:\$PATH\""
