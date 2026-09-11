/**
 * SV Print client for JavaScript/TypeScript.
 *
 * Usage (browser or Node.js):
 *   import { SVPrint } from "./sv-print-client.js";
 *
 *   const sv = new SVPrint({ token: "your-token", port: 9876 });
 *   const info = await sv.info();
 *   const printers = await sv.printers();
 *   const result = await sv.printReceipt(printers[0].id, {
 *     cut: true,
 *     lines: [
 *       { text: "HELLO", style: { bold: true, align: "center" } },
 *     ],
 *   });
 */

export class SVPrint {
  /** @type {string} */
  #base;
  /** @type {string} */
  #token;

  /**
   * @param {object} opts
   * @param {string} opts.token  - Auth token
   * @param {string} [opts.host="127.0.0.1"] - Agent host
   * @param {number} [opts.port=9876] - Agent port
   */
  constructor({ token, host = "127.0.0.1", port = 9876 }) {
    if (!token) throw new Error("token is required");
    this.#token = token;
    this.#base = `http://${host}:${port}`;
  }

  /** Health check (no auth). */
  async health() {
    const r = await fetch(`${this.#base}/health`);
    return r.json();
  }

  /** Agent info (name, version, tier, etc.). */
  async info() {
    return this.#get("/api/v1/info");
  }

  /** List detected printers. */
  async printers() {
    const r = await this.#get("/api/v1/printers");
    return r.printers;
  }

  /** Re-run printer discovery. */
  async discover() {
    const r = await this.#post("/api/v1/printers/discover", {});
    return r.printers;
  }

  /** Get a single printer by ID. */
  async getPrinter(id) {
    return this.#get(`/api/v1/printers/${encodeURIComponent(id)}`);
  }

  /** Send a test receipt to a printer. */
  async testReceipt(printerId) {
    return this.#post(
      `/api/v1/printers/${encodeURIComponent(printerId)}/test`,
      {}
    );
  }

  /**
   * Print a structured receipt.
   * @param {string} printerId
   * @param {{ cut?: boolean, lines: Array<{ text: string, style?: object }> }} doc
   * @returns {Promise<{ job_id: string, status: string }>}
   */
  async printReceipt(printerId, doc) {
    return this.#post("/api/v1/print/receipt", {
      printer_id: printerId,
      ...doc,
    });
  }

  /**
   * Print raw ESC/POS bytes (requires license).
   * @param {string} printerId
   * @param {string} payload - Base64-encoded ESC/POS bytes
   * @returns {Promise<{ job_id: string, status: string }>}
   */
  async rawPrint(printerId, payload) {
    return this.#post("/api/v1/print", {
      printer_id: printerId,
      payload,
      format: "escpos",
    });
  }

  /** Get job status. */
  async getJob(id) {
    return this.#get(`/api/v1/jobs/${encodeURIComponent(id)}`);
  }

  /**
   * Subscribe to events via WebSocket.
   * Works in browsers (uses ?token= query param).
   * @param {(event: object) => void} onEvent
   * @param {string} [host="127.0.0.1"]
   * @param {number} [port=9876]
   * @returns {WebSocket}
   */
  onEvent(onEvent, host = "127.0.0.1", port = 9876) {
    const ws = new WebSocket(
      `ws://${host}:${port}/api/v1/events?token=${this.#token}`
    );
    ws.onmessage = (e) => onEvent(JSON.parse(e.data));
    return ws;
  }

  async #get(path) {
    const r = await fetch(`${this.#base}${path}`, {
      headers: { Authorization: `Bearer ${this.#token}` },
    });
    return this.#handle(r);
  }

  async #post(path, body) {
    const r = await fetch(`${this.#base}${path}`, {
      method: "POST",
      headers: {
        Authorization: `Bearer ${this.#token}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    });
    return this.#handle(r);
  }

  async #handle(r) {
    const data = await r.json();
    if (!r.ok) {
      const msg = data?.error?.message || `HTTP ${r.status}`;
      throw new Error(msg);
    }
    return data;
  }
}
