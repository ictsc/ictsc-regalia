"""Print only differing OCI layer paths; never emit image contents."""

import hashlib
import io
import json
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
                        value += (hashlib.sha256(content.extractfile(member).read()).hexdigest(),)
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


if __name__ == "__main__":
    main(*sys.argv[1:])
