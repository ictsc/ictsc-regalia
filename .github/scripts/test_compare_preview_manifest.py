import tempfile
from pathlib import Path
import unittest
from compare_preview_manifest import same_manifest


class ManifestTests(unittest.TestCase):
    def test_same_raw_bytes_are_idempotent(self):
        with tempfile.TemporaryDirectory() as tmp:
            source, existing = Path(tmp) / "source", Path(tmp) / "existing"
            source.write_bytes(b'{"schemaVersion":2}\n')
            existing.write_bytes(source.read_bytes())
            self.assertTrue(same_manifest(source, existing))

    def test_different_bytes_are_rejected(self):
        with tempfile.TemporaryDirectory() as tmp:
            source, existing = Path(tmp) / "source", Path(tmp) / "existing"
            source.write_bytes(b'{"schemaVersion":2}\n')
            existing.write_bytes(b'{"schemaVersion":2}')
            self.assertFalse(same_manifest(source, existing))


if __name__ == "__main__":
    unittest.main()
