"""
SV Print Python client.

Usage:
    from sv_print_client import SVPrint

    sv = SVPrint(token="your-token")
    printers = sv.printers()
    sv.print_receipt(printers[0]["id"], lines=[
        {"text": "HELLO", "style": {"bold": True, "align": "center"}},
        {"text": "World"},
    ])
"""

try:
    import requests
except ImportError:
    raise ImportError("pip install requests")

import json
import base64


class SVPrint:
    def __init__(self, token, host="127.0.0.1", port=9876):
        self.base = f"http://{host}:{port}"
        self.token = token

    def health(self):
        return requests.get(f"{self.base}/health").json()

    def info(self):
        return self._get("/api/v1/info")

    def printers(self):
        return self._get("/api/v1/printers")["printers"]

    def discover(self):
        return self._post("/api/v1/printers/discover", {})["printers"]

    def get_printer(self, printer_id):
        return self._get(f"/api/v1/printers/{printer_id}")

    def test_receipt(self, printer_id):
        return self._post(f"/api/v1/printers/{printer_id}/test", {})

    def print_receipt(self, printer_id, lines, cut=True):
        return self._post("/api/v1/print/receipt", {
            "printer_id": printer_id,
            "cut": cut,
            "lines": lines,
        })

    def raw_print(self, printer_id, escpos_bytes):
        return self._post("/api/v1/print", {
            "printer_id": printer_id,
            "payload": base64.b64encode(escpos_bytes).decode(),
            "format": "escpos",
        })

    def get_job(self, job_id):
        return self._get(f"/api/v1/jobs/{job_id}")

    def _get(self, path):
        r = requests.get(
            f"{self.base}{path}",
            headers={"Authorization": f"Bearer {self.token}"},
        )
        r.raise_for_status()
        return r.json()

    def _post(self, path, body):
        r = requests.post(
            f"{self.base}{path}",
            json=body,
            headers={"Authorization": f"Bearer {self.token}"},
        )
        r.raise_for_status()
        return r.json()
