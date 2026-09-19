> **Language:** Español | [English](integration.md)

# Guía de Integración

Cómo ejecutar SV Printer, configurarlo e integrarlo con tu aplicación web
usando la API HTTP local.

## Requisitos

- Go (ver [`go.mod`](../go.mod))
- Una impresora térmica conectada vía USB, serial o red (TCP)

## Ejecutar el agente

```bash
# Compilar
go build -o sv-printer ./cmd/sv-printer

# Ejecutar (auto-detecta impresoras, genera un token en la primera ejecución)
./sv-printer

# O configurar manualmente
./sv-printer -token mysecret -port 9876 -printer "Caja 1@192.168.1.100:9100"
```

En la primera ejecución sin `-token`, SV Printer genera uno y lo imprime en stdout.
Guárdalo — lo necesitarás para cada petición a la API.

### Configuración

| Flag | Variable de entorno | Valor por defecto | Descripción |
|---|---|---|---|
| `-host` | `SV_PRINT_HOST` | `127.0.0.1` | Dirección de vinculación |
| `-port` | `SV_PRINT_PORT` | `9876` | Puerto de vinculación |
| `-token` | `SV_PRINT_TOKEN` | (generado) | Token de autenticación |
| `-max-payload` | `SV_PRINT_MAX_PAYLOAD` | `5242880` (5 MB) | Tamaño máximo del body |
| `-origin` | `SV_PRINT_ALLOWED_ORIGINS` | (ninguno) | Origen CORS permitido (repetible) |
| `-printer` | — | (ninguno) | Impresora manual como `nombre@dirección` (repetible) |
| `-serial` | — | (ninguno) | Impresora serial manual como `nombre@puerto[@baud]` (repetible) |
| `-config` | `SV_PRINT_CONFIG` | ruta de plataforma | Ruta del archivo de configuración |
| `-license` | `SV_PRINT_LICENSE` | `license.key` junto a config | Ruta del archivo de licencia |

El archivo de configuración es JSON (`sv-printer config` muestra su ruta).

### Subcomandos

| Comando | Descripción |
|---|---|
| `sv-printer version` | Mostrar versión |
| `sv-printer status` | Verificar estado del agente |
| `sv-printer printers` | Listar impresoras configuradas |
| `sv-printer config` | Mostrar configuración actual |
| `sv-printer discover` | Redescubrir impresoras |
| `sv-printer test <id-impresora>` | Enviar un recibo de prueba |
| `sv-printer print <id-impresora> <archivo.json>` | Imprimir un recibo estructurado desde un archivo JSON |
| `sv-printer logs` | Seguir archivo de log |
| `sv-printer doctor` | Ejecutar diagnósticos |
| `sv-printer device-id` | Mostrar huella del dispositivo (para vinculación de licencia) |

## Modelo de seguridad

- El agente se vincula **solo a `127.0.0.1`** (nunca a `0.0.0.0` por defecto).
- Todos los endpoints `/api/*` requieren `Authorization: Bearer <token>`.
- CORS: solo los orígenes listados en `-origin` / `allowed_origins` pueden llamar
  a la API desde un navegador. Configura esto al origen de tu app web
  (ej. `https://svtech.cl`).
- Autenticación WebSocket: la API WebSocket nativa del navegador no puede setear
  headers, así que `?token=<token>` se acepta como alternativa para
  `GET /api/v1/events`.

## Inicio rápido (curl)

```bash
TOKEN="tu-token-aquí"
BASE="http://127.0.0.1:9876"

# Salud (sin auth)
curl http://127.0.0.1:9876/health

# Info
curl -H "Authorization: Bearer $TOKEN" $BASE/api/v1/info

# Listar impresoras
curl -H "Authorization: Bearer $TOKEN" $BASE/api/v1/printers

# Recibo de prueba (trial: agrega marca de agua, cuenta para la cuota)
curl -X POST -H "Authorization: Bearer $TOKEN" \
  $BASE/api/v1/printers/<id-impresora>/test
```

## Imprimir un recibo estructurado

`POST /api/v1/print/receipt` acepta un body JSON que describe el diseño del recibo.

### Esquema de la petición

```json
{
  "printer_id": "net-192.168.1.100:9100",
  "cut": true,
  "lines": [
    { "text": "MI TIENDA", "style": { "bold": true, "align": "center" } },
    { "text": "Calle Principal 123" },
    { "text": "-------------------" },
    { "text": "Artículo 1", "style": { "align": "left" } },
    { "text": "$10.00", "style": { "align": "right" } },
    { "text": "TOTAL", "style": { "bold": true, "align": "center", "size": [2, 2] } }
  ]
}
```

### Campos

| Campo | Tipo | Descripción |
|---|---|---|
| `printer_id` | string | **Requerido.** ID de la impresora (de `/api/v1/printers`) |
| `cut` | bool | Avanzar y cortar el papel después de imprimir |
| `lines` | array | **Requerido.** Arreglo de objetos de línea |

### Objeto de línea

| Campo | Tipo | Descripción |
|---|---|---|
| `text` | string | **Requerido.** Texto a imprimir |
| `style.bold` | bool | Texto en negrita |
| `style.underline` | bool | Texto subrayado |
| `style.italic` | bool | Texto en cursiva |
| `style.reverse` | bool | Invertido (blanco sobre negro) |
| `style.align` | string | `left` (por defecto), `center` o `right` |
| `style.size` | `[w, h]` | Multiplicador de tamaño de carácter `[ancho, alto]` (1–8 cada uno) |

### Ejemplo (curl)

```bash
curl -X POST -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "printer_id": "net-192.168.1.100:9100",
    "cut": true,
    "lines": [
      { "text": "ORDEN #42", "style": { "bold": true, "align": "center" } },
      { "text": "Café x2", "style": { "align": "left" } },
      { "text": "$8.00", "style": { "align": "right" } },
      { "text": "GRACIAS", "style": { "align": "center" } }
    ]
  }' \
  $BASE/api/v1/print/receipt
```

Respuesta (201 Created):

```json
{ "job_id": "job_f3a9c1d2e4b5061789a0b2c3d4e5f607", "status": "queued" }
```

## Impresión ESC/POS raw (pro)

`POST /api/v1/print` envía bytes ESC/POS raw. Requiere una licencia (función pro).

```json
{
  "printer_id": "net-192.168.1.100:9100",
  "payload": "AQIDBA==",
  "format": "escpos"
}
```

`payload` es un string **codificado en base64** de comandos ESC/POS raw.

## Estado del trabajo

`GET /api/v1/jobs/{id}` retorna:

```json
{
  "id": "job_f3a9c1d2e4b5061789a0b2c3d4e5f607",
  "printer_id": "net-192.168.1.100:9100",
  "created_at": "2026-09-11T12:00:00Z",
  "status": "completed"
}
```

Valores de estado: `queued` → `printing` → `completed` | `failed` | `cancelled`.

## Eventos WebSocket (pro)

`GET /api/v1/events` abre una conexión WebSocket. Requiere una licencia (función pro).

### Autenticación

Los navegadores no pueden setear headers personalizados en WebSocket, así que usa
un parámetro de query:

```
ws://127.0.0.1:9876/api/v1/events?token=tu-token
```

Los clientes Node.js / que no son navegadores pueden usar el header `Authorization`:

```
ws://127.0.0.1:9876/api/v1/events
Authorization: Bearer tu-token
```

### Tipos de eventos

| Evento | Datos |
|---|---|
| `agent.status` | `{ "status": "started" }` / `{ "status": "stopped" }` |
| `print.started` | `{ "job_id": "...", "printer_id": "..." }` |
| `print.completed` | `{ "job_id": "...", "printer_id": "..." }` |
| `print.failed` | `{ "job_id": "...", "printer_id": "...", "error": "..." }` |
| `printer.discovered` | `{ "id": "...", "name": "..." }` |
| `printer.removed` | `{ "id": "..." }` |
| `printer.status_changed` | `{ "id": "...", "status": "..." }` |

### Ejemplo (JavaScript)

```js
const ws = new WebSocket(
  "ws://127.0.0.1:9876/api/v1/events?token=tu-token"
);

ws.onmessage = (e) => {
  const event = JSON.parse(e.data);
  console.log("Evento:", event.event, event);
};
```

## Respuestas de error

Todos los errores siguen un envoltorio JSON consistente:

```json
{
  "error": {
    "code": "PRINTER_NOT_FOUND",
    "message": "Impresora no encontrada."
  }
}
```

### Códigos de error

| HTTP | Código | Significado |
|---|---|---|
| 400 | `INVALID_PAYLOAD` | Body de petición malformado |
| 400 | `UNSUPPORTED_PROTOCOL` | `format` desconocido en impresión raw |
| 401 | `UNAUTHORIZED` | Token faltante o inválido |
| 403 | `LICENSE_REQUIRED` | Función pro sin licencia |
| 403 | `ORIGIN_NOT_ALLOWED` | Origen CORS no está en la lista permitida |
| 404 | `PRINTER_NOT_FOUND` | ID de impresora desconocido |
| 404 | `JOB_NOT_FOUND` | ID de trabajo desconocido |
| 409 | `PRINTER_BUSY` | La impresora está ocupada |
| 409 | `PRINTER_PAPER_OUT` | La impresora no tiene papel |
| 409 | `JOB_ALREADY_EXISTS` | ID de trabajo duplicado |
| 413 | `PAYLOAD_TOO_LARGE` | El body excede el límite |
| 429 | `LICENSE_QUOTA_EXCEEDED` | Límite diario de trial alcanzado (50/día) |
| 502 | `PRINTER_CONNECTION_FAILED` | No se puede conectar a la impresora |
| 503 | `PRINTER_OFFLINE` | La impresora está offline |
| 503 | `EVENTS_UNAVAILABLE` | El bus de eventos no está configurado |
| 500 | `PRINT_FAILED` | El trabajo de impresión falló |
| 500 | `INTERNAL_ERROR` | Error interno del servidor |

## Notas de licencia / trial

Sin una licencia (modo trial):

- Los recibos se prefijan y sufijan con
  `*** SV PRINT — LICENCIA DE PRUEBA ***`.
- Puedes imprimir hasta **50 recibos por día** (la cuota se reinicia a medianoche).
- `POST /api/v1/print` (ESC/POS raw) y `GET /api/v1/events` (WebSocket)
  están **bloqueados** (403 `LICENSE_REQUIRED`).

Para instalar una licencia:

```bash
./sv-printer -license ./license.key
```

## Solución de problemas

```bash
# Ejecutar diagnósticos
sv-printer doctor

# Revisar logs
sv-printer logs

# Mostrar huella del dispositivo (para vinculación de licencia)
sv-printer device-id
```
