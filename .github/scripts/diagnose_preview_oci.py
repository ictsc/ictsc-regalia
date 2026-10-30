"""Print only differing OCI layer paths; never emit image contents."""

import hashlib
import io
import json
import re
import sys
import tarfile


def read_oci(path):
    with tarfile.open(path) as archive:
        def blob(digest):
            return archive.extractfile("blobs/sha256/" + digest.removeprefix("sha256:")).read()

        index = json.load(archive.extractfile("index.json"))
        manifest = json.loads(blob(index["manifests"][0]["digest"]))
        config = json.loads(blob(manifest["config"]["digest"]))
        layers = []
        for layer in manifest["layers"]:
            files = {}
            with tarfile.open(fileobj=io.BytesIO(blob(layer["digest"])), mode="r:*") as content:
                for member in content:
                    value = (member.type, member.mode, member.linkname)
                    if member.isfile():
                        data = content.extractfile(member).read()
                        value += (hashlib.sha256(data).hexdigest(),
                                  data if member.name.endswith(".html") else b"")
                    files[member.name] = value
            layers.append(files)
        return config, layers


def main(first, second):
    config_a, layers_a = read_oci(first)
    config_b, layers_b = read_oci(second)
    print("Config fields that differ:", sorted(
        key for key in config_a.keys() | config_b.keys()
        if config_a.get(key) != config_b.get(key)))
    for index in range(max(len(layers_a), len(layers_b))):
        a = layers_a[index] if index < len(layers_a) else {}
        b = layers_b[index] if index < len(layers_b) else {}
        different = sorted(path for path in a.keys() | b.keys() if a.get(path) != b.get(path))
        if different:
            print(f"Layer {index}: {len(different)} differing paths")
            for path in different[:30]:
                print("  ", path)
                first, second = a.get(path), b.get(path)
                if first and second and path.endswith(".html"):
                    left, right = first[-1], second[-1]
                    for label, pattern in (
                        ("epoch milliseconds", rb"\b\d{13}\b"),
                        ("ISO datetime", rb"20\d\d-\d\d-\d\dT[\d:.]+Z"),
                        ("UUID", rb"\b[0-9a-f]{8}-[0-9a-f-]{27,}\b"),
                    ):
                        if re.sub(pattern, b"<normalized>", left) == re.sub(
                                pattern, b"<normalized>", right):
                            print("    difference class:", label)
                            break
                    else:
                        offset = next((i for i, (x, y) in enumerate(zip(left, right))
                                       if x != y), min(len(left), len(right)))
                        print("    first differing byte offset:", offset,
                              "lengths:", len(left), len(right))


if __name__ == "__main__":
    main(*sys.argv[1:])
