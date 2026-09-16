# SV Print JavaScript/TypeScript Client

A lightweight fetch-based client for the SV Print local HTTP API.

## Usage (browser)

```html
<script type="module">
  import { SVPrint } from "./sv-printer-client.js";

  const sv = new SVPrint({ token: "your-token" });

  // List printers
  const printers = await sv.printers();
  console.log(printers);

  // Print a receipt
  const { job_id } = await sv.printReceipt(printers[0].id, {
    cut: true,
    lines: [
      { text: "HELLO", style: { bold: true, align: "center" } },
      { text: "World" },
    ],
  });

  // Subscribe to events
  sv.onEvent((ev) => console.log("Event:", ev.event, ev));
</script>
```

## Usage (Node.js)

```js
import { SVPrint } from "./sv-printer-client.js";

const sv = new SVPrint({ token: process.env.SV_PRINT_TOKEN });

const printers = await sv.printers();
await sv.printReceipt(printers[0].id, {
  cut: true,
  lines: [{ text: "Hello from Node!" }],
});
```

## API

| Method | Description |
|---|---|
| `health()` | Health check (no auth) |
| `info()` | Agent info |
| `printers()` | List printers |
| `discover()` | Re-run discovery |
| `getPrinter(id)` | Get a printer by ID |
| `testReceipt(printerId)` | Send test receipt |
| `printReceipt(printerId, doc)` | Print structured receipt |
| `rawPrint(printerId, payload)` | Print raw ESC/POS (requires license) |
| `getJob(id)` | Get job status |
| `onEvent(callback)` | Subscribe to WebSocket events |

## See also

- [Integration Guide](../../documentation/integration.md)
- [API Reference](../../documentation/integration.md#quickstart-curl)
