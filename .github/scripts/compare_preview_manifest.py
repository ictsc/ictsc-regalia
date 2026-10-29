#!/usr/bin/env python3
"""Reject a GHCR tag that points at different raw OCI manifest bytes."""
import hashlib
from pathlib import Path
import sys


def same_manifest(source, destination):
    return hashlib.sha256(Path(source).read_bytes()).digest() == hashlib.sha256(Path(destination).read_bytes()).digest()


if __name__ == "__main__":
    if len(sys.argv) != 3 or not same_manifest(sys.argv[1], sys.argv[2]):
        raise SystemExit("OCI tag digest differs from PR artifact")
