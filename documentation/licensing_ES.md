> **Language:** Español | [English](licensing.md)

# Licencias

SV Printer usa un modelo **freemium** con licencias firmadas offline usando Ed25519.

## Modos

| Modo | Activación | Marca de agua | Cuota diaria | Funciones pro |
|---|---|---|---|---|
| **Trial** (por defecto) | ninguno | sí, al inicio y final de cada recibo | 50 impresiones/día | no |
| **Beta** | licencia firmada con `expiry` | no | ilimitado | según `features` |
| **Full** | licencia firmada | no | ilimitado | según `features` |

La marca de agua de trial es `*** SV PRINT — LICENCIA DE PRUEBA ***`.

## Funciones protegidas

| Función | Desbloquea |
|---|---|
| `raw_print` | `POST /api/v1/print` (payloads ESC/POS raw) |
| `websocket` | `GET /api/v1/events` (flujo de eventos WebSocket) |

Una licencia debe listar una función en su array `features` para activarla. Sin
una licencia (modo trial), solo el endpoint de recibo estructurado
(`POST /api/v1/print/receipt`) y los endpoints de impresora/prueba están disponibles.

## Cómo funcionan las licencias

1. Se genera una **parella de claves de firma** una sola vez. La **clave pública**
   está embebida en el binario del agente (`internal/license/public.key`); la
   **clave privada** se mantiene en secreto y nunca se commitea al repositorio.
2. Cada licencia es un envoltorio JSON:
   ```json
   {
     "payload": "<claims codificados en base64>",
     "signature": "<firma Ed25519 en base64 sobre el payload>"
   }
   ```
   Los claims son `license_id`, `product`, `customer`, `tier`, `features`,
   `expiry` (RFC3339 opcional) y `fingerprint` (opcional).
3. El agente verifica la firma contra la clave pública embebida al iniciar.
   Una licencia faltante, manipulada, expirada o de producto incorrecto
   regresa al modo trial.

## Emitir una licencia

```bash
# 1. Generar la pareja de claves (una sola vez).
go run ./cmd/sv-license keygen
#   -> clave privada (hex) impresa en stderr: guárdala en 1Password / un gestor de secretos
#   -> clave pública (hex) impresa en stdout: colócala en internal/license/public.key

# 2. Firmar una licencia.
SV_LICENSE_KEY=<clave-privada-hex> go run ./cmd/sv-license sign \
  -customer "ACME" \
  -tier full \
  -features raw_print,websocket \
  > license.key

# Con expiración (ej. una beta):
SV_LICENSE_KEY=<clave-privada-hex> go run ./cmd/sv-license sign \
  -customer "ACME" \
  -tier beta \
  -expiry 2027-01-01T00:00:00Z \
  -features raw_print,websocket \
  > license.key
```

## Instalar una licencia en el agente

El agente carga su licencia desde, en orden:

1. el flag `-license` (`sv-printer -license ./license.key`),
2. la variable de entorno `SV_PRINT_LICENSE`,
3. `license.key` junto al archivo de configuración (ubicación por defecto).

El estado de la licencia actual se expone por `GET /api/v1/info` a través de los
campos `tier` y `licensed`.

## Vinculación a dispositivo

Una licencia puede vincularse a una sola máquina configurando el claim `fingerprint`.
Cuando hay un fingerprint presente, el agente lo verifica contra el ID del dispositivo
(vía [`machineid`](https://github.com/denisbrodbeck/machineid)) y regresa al
modo trial si no coincide.

Flujo:

1. El cliente ejecuta `sv-printer device-id` y te envía la huella.
2. Tú firmas la licencia con `-fingerprint <valor>`.
3. La licencia solo se activa en esa máquina.

```bash
# Cliente:
sv-printer device-id
# -> ej. 3f8c2a...

# Vendedor:
go run ./cmd/sv-license sign \
  -customer "ACME" \
  -tier full \
  -features raw_print,websocket \
  -fingerprint "3f8c2a..." \
  > license.key
```

## Notas de seguridad

- La **clave privada nunca debe** ser commiteada o compartida. Cualquiera que la
  tenga puede firmar licencias arbitrarias.
- Para un repositorio público, este es el control central: incluso con el código
  completo, una licencia no puede ser falsificada sin la clave privada.
- La **vinculación a dispositivo** previene compartir una licencia entre máquinas,
  pero es impuesta por software y puede ser superada por un atacante que parchee
  el binario.
- Las licencias solo offline no pueden revocar una licencia ni prevenir el
  retroceso del reloj. Si la revocación o el medición estricta se convierten
  en un requisito, agregar un paso de validación phone-home (ver la hoja de ruta).
