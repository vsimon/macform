#!/usr/bin/env bash
set -euo pipefail

ARTIFACT="$1"

echo "-> Signing $ARTIFACT"
codesign \
  --sign "$APPLE_SIGNING_IDENTITY" \
  --options runtime \
  --timestamp \
  --force \
  "$ARTIFACT"

codesign --verify --verbose "$ARTIFACT"

echo "-> Notarizing $ARTIFACT"
TMPZIP="$(mktemp "${TMPDIR:-/tmp}/notarize-XXXXXX").zip"
TMPKEY=$(mktemp "${TMPDIR:-/tmp}/notarize-XXXXXX.p8")
rm -f "$TMPZIP"
trap 'rm -f "$TMPZIP" "$TMPKEY"' EXIT

echo "$APPLE_API_KEY" > "$TMPKEY"

zip -j "$TMPZIP" "$ARTIFACT"

xcrun notarytool submit "$TMPZIP" \
  --key "$TMPKEY" \
  --key-id "$APPLE_API_KEY_ID" \
  --issuer "$APPLE_API_ISSUER" \
  --wait \
  --timeout 20m

echo "-> Done: $ARTIFACT"
