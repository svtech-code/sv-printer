<?php

/**
 * SV Printer PHP client.
 *
 * Usage:
 *   $sv = new SVPrint("your-token");
 *   $printers = $sv->printers();
 *   $sv->printReceipt($printers[0]["id"], [
 *       ["text" => "HELLO", "style" => ["bold" => true, "align" => "center"]],
 *       ["text" => "World"],
 *   ]);
 */

class SVPrint
{
    private string $base;
    private string $token;

    public function __construct(string $token, string $host = "127.0.0.1", int $port = 9876)
    {
        $this->token = $token;
        $this->base = "http://{$host}:{$port}";
    }

    public function health(): array
    {
        return $this->get("/health");
    }

    public function info(): array
    {
        return $this->get("/api/v1/info");
    }

    public function printers(): array
    {
        return $this->get("/api/v1/printers")["printers"];
    }

    public function discover(): array
    {
        return $this->post("/api/v1/printers/discover", [])["printers"];
    }

    public function getPrinter(string $id): array
    {
        return $this->get("/api/v1/printers/" . rawurlencode($id));
    }

    public function testReceipt(string $printerId): array
    {
        return $this->post("/api/v1/printers/" . rawurlencode($printerId) . "/test", []);
    }

    public function printReceipt(string $printerId, array $lines, bool $cut = true): array
    {
        return $this->post("/api/v1/print/receipt", [
            "printer_id" => $printerId,
            "cut" => $cut,
            "lines" => $lines,
        ]);
    }

    public function rawPrint(string $printerId, string $escposBytes): array
    {
        return $this->post("/api/v1/print", [
            "printer_id" => $printerId,
            "payload" => base64_encode($escposBytes),
            "format" => "escpos",
        ]);
    }

    public function getJob(string $jobId): array
    {
        return $this->get("/api/v1/jobs/" . rawurlencode($jobId));
    }

    private function get(string $path): array
    {
        $ch = curl_init($this->base . $path);
        curl_setopt_array($ch, [
            CURLOPT_RETURNTRANSFER => true,
            CURLOPT_HTTPHEADER => ["Authorization: Bearer {$this->token}"],
        ]);
        $body = curl_exec($ch);
        $code = curl_getinfo($ch, CURLINFO_HTTP_CODE);
        curl_close($ch);

        if ($code >= 400) {
            $data = json_decode($body, true);
            $msg = $data["error"]["message"] ?? "HTTP {$code}";
            throw new RuntimeException($msg, $code);
        }
        return json_decode($body, true);
    }

    private function post(string $path, array $data): array
    {
        $ch = curl_init($this->base . $path);
        curl_setopt_array($ch, [
            CURLOPT_RETURNTRANSFER => true,
            CURLOPT_POST => true,
            CURLOPT_POSTFIELDS => json_encode($data),
            CURLOPT_HTTPHEADER => [
                "Authorization: Bearer {$this->token}",
                "Content-Type: application/json",
            ],
        ]);
        $body = curl_exec($ch);
        $code = curl_getinfo($ch, CURLINFO_HTTP_CODE);
        curl_close($ch);

        if ($code >= 400) {
            $resp = json_decode($body, true);
            $msg = $resp["error"]["message"] ?? "HTTP {$code}";
            throw new RuntimeException($msg, $code);
        }
        return json_decode($body, true);
    }
}
