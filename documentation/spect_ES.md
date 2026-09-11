# SV Print Agent

> **Versión canónica:** la versión canónica de esta especificación es
> [`spect.md`](./spect.md) (en inglés). Este documento es un espejo en español
> para referencia únicamente.

## 1. Descripción general

**SV Print Agent** es un agente local multiplataforma desarrollado en **Go**, cuyo objetivo es permitir que aplicaciones web se comuniquen con impresoras térmicas conectadas al equipo del usuario.

El agente actúa como puente entre:

```text
Aplicación Web
      │
      │ HTTP / WebSocket
      ▼
SV Print Agent
      │
      ├── USB
      ├── Serial
      ├── TCP/IP
      └── Impresora del sistema
      │
      ▼
Impresora térmica
```

El objetivo principal es permitir la impresión térmica desde aplicaciones web sin depender del diálogo de impresión del navegador.

El agente debe funcionar en:

- Windows
- macOS
- Linux

El protocolo principal para el MVP será **ESC/POS**.

---

# 2. Objetivos

## 2.1 Objetivos principales

El sistema debe:

- Detectar impresoras disponibles localmente.
- Identificar el tipo de conexión de cada impresora.
- Permitir configurar impresoras manualmente.
- Exponer una API local para aplicaciones web.
- Imprimir documentos mediante ESC/POS.
- Soportar impresoras USB.
- Soportar impresoras seriales.
- Soportar impresoras de red.
- Consultar el estado de las impresoras cuando el hardware/protocolo lo permita.
- Funcionar en Windows, macOS y Linux.
- Ejecutarse como proceso o servicio en segundo plano.
- Proporcionar una CLI para diagnóstico y administración.
- Proteger la API mediante autenticación.
- Mantener una arquitectura independiente de cualquier framework web.

---

## 2.2 Objetivos secundarios

El proyecto debe estar preparado para incorporar posteriormente:

- Descubrimiento automático avanzado.
- Cola persistente de impresión.
- Reintentos automáticos.
- Monitoreo de impresoras.
- Eventos mediante WebSocket.
- Múltiples impresoras.
- Alias para impresoras.
- Configuración persistente.
- Historial de trabajos.
- Integración con CUPS.
- Integración con Windows Print Spooler.
- Interfaz gráfica.
- System Tray.
- Administración remota.

---

# 3. Alcance del MVP

El MVP debe concentrarse en las funcionalidades esenciales.

### Sistemas operativos

```text
Windows
macOS
Linux
```

### Conexiones

```text
USB
Serial
TCP/IP
```

### Funcionalidades

```text
✓ Detección de impresoras
✓ Configuración manual
✓ Listado de impresoras
✓ Prueba de impresión
✓ Impresión ESC/POS
✓ Cola de impresión
✓ Estado del trabajo
✓ API HTTP local
✓ Autenticación
✓ Restricción CORS
✓ CLI
✓ Logs estructurados
✓ Health check
```

---

# 4. Tecnologías

## 4.1 Lenguaje

El proyecto debe desarrollarse en:

```text
Go
```

Utilizar la versión estable de Go disponible al iniciar el desarrollo.

---

## 4.2 Principios técnicos

El proyecto debe priorizar:

- Simplicidad.
- Portabilidad.
- Bajo consumo de recursos.
- Seguridad.
- Testabilidad.
- Separación de responsabilidades.
- Interfaces pequeñas.
- Código idiomático de Go.
- Dependencias mínimas.
- Evitar CGO cuando sea técnicamente posible.

---

# 5. Arquitectura

Se utilizará una arquitectura modular inspirada en Clean Architecture.

Estructura propuesta:

```text
sv-print/
├── cmd/
│   └── sv-print/
│       └── main.go
│
├── internal/
│   ├── application/
│   │   ├── discovery/
│   │   ├── printing/
│   │   └── health/
│   │
│   ├── domain/
│   │   ├── printer/
│   │   ├── job/
│   │   └── errors/
│   │
│   ├── infrastructure/
│   │   ├── usb/
│   │   ├── serial/
│   │   ├── network/
│   │   ├── system/
│   │   ├── storage/
│   │   └── logging/
│   │
│   └── interfaces/
│       ├── http/
│       ├── websocket/
│       └── cli/
│
├── pkg/
│   └── escpos/
│
├── docs/
├── scripts/
├── test/
│
├── go.mod
└── README.md
```

La estructura podrá modificarse durante el desarrollo si existe una razón técnica, pero se debe mantener la separación entre dominio, aplicación, infraestructura e interfaces.

---

# 6. Arquitectura de alto nivel

```text
                         APLICACIÓN WEB
                               │
                               │ HTTP
                               ▼
                       API LOCAL DE SV PRINT
                               │
                     ┌─────────┴─────────┐
                     │                   │
                    HTTP             WebSocket
                     │                   │
                     └─────────┬─────────┘
                               │
                      CAPA DE APLICACIÓN
                               │
                 ┌─────────────┴─────────────┐
                 │                           │
          Descubrimiento                Cola de impresión
          de impresoras                       │
                 │                       Print Worker
                 │                           │
                 └─────────────┬─────────────┘
                               │
                    Abstracción de impresora
                               │
             ┌─────────────────┼─────────────────┐
             │                 │                 │
            USB             Serial              TCP
             │                 │                 │
             └─────────────────┼─────────────────┘
                               │
                       IMPRESORA TÉRMICA
```

---

# 7. Soporte multiplataforma

## 7.1 Windows

Versiones objetivo:

```text
Windows 10+
Windows 11
```

Debe soportar:

- Detección USB.
- Detección serial.
- Impresoras de red.
- Integración futura con Windows Print Spooler.
- Ejecución como Windows Service.

Binario:

```text
sv-print.exe
```

Instalador futuro:

```text
sv-print-installer.exe
```

---

## 7.2 macOS

Versión mínima inicial:

```text
macOS 12+
```

Debe soportar:

- USB.
- Serial.
- TCP/IP.
- CUPS.
- Ejecución en segundo plano.

Arquitecturas:

```text
amd64
arm64
```

Distribución futura:

```text
SV Print.app
.dmg
.pkg
```

La distribución de producción deberá considerar:

- Firma de código.
- Notarización.
- Permisos del sistema.

---

## 7.3 Linux

Distribuciones iniciales:

```text
Ubuntu
Debian
```

Arquitecturas:

```text
amd64
arm64
```

Debe soportar:

- USB.
- Serial.
- TCP/IP.
- CUPS.
- systemd.

Paquete futuro:

```text
.deb
```

Posibles formatos posteriores:

```text
.rpm
AppImage
```

---

# 8. Detección de impresoras

La detección de impresoras es una funcionalidad fundamental.

El agente debe disponer de un servicio de descubrimiento:

```text
Printer Discovery
       │
       ├── USB
       ├── Serial
       ├── Network
       └── System
```

Cada mecanismo debe implementar una interfaz común.

Ejemplo:

```go
type PrinterDiscovery interface {
    Discover(ctx context.Context) ([]Printer, error)
}
```

La capa de aplicación no debe depender directamente de APIs específicas de Windows, macOS o Linux.

---

# 9. Detección USB

El agente debe detectar impresoras conectadas mediante USB.

Cuando sea posible debe obtener:

- Vendor ID.
- Product ID.
- Fabricante.
- Modelo.
- Número de serie.
- Ruta del dispositivo.
- Información de interfaz.

Ejemplo:

```json
{
  "id": "usb-1234-5678",
  "name": "XPrinter XP-Q200",
  "connection": "usb",
  "manufacturer": "XPrinter",
  "model": "XP-Q200"
}
```

La implementación debe utilizar una biblioteca Go mantenida y compatible con los tres sistemas operativos.

La biblioteca definitiva debe seleccionarse considerando:

- Compatibilidad multiplataforma.
- Mantenimiento.
- Licencia.
- Necesidad de CGO.
- Compatibilidad con dispositivos térmicos.

---

# 10. Detección serial

Debe detectar dispositivos seriales.

### Windows

```text
COM3
COM4
```

### Linux

```text
/dev/ttyUSB0
/dev/ttyACM0
```

### macOS

```text
/dev/cu.usbserial-XXXX
/dev/cu.usbmodemXXXX
```

Configuración:

```text
Baud rate
Data bits
Stop bits
Parity
Flow control
```

Configuración inicial recomendada:

```text
9600
8
N
1
```

Debe ser configurable.

---

# 11. Impresoras de red

Debe soportarse impresión RAW mediante TCP.

Puerto inicial:

```text
9100
```

Ejemplo:

```text
192.168.1.100:9100
```

Configuración:

```json
{
  "name": "Caja 1",
  "connection": "network",
  "address": "192.168.1.100:9100",
  "protocol": "escpos"
}
```

El descubrimiento automático avanzado mediante:

```text
mDNS
DNS-SD
SNMP
```

queda inicialmente como funcionalidad futura.

No se debe realizar escaneo indiscriminado de redes por defecto.

---

# 12. Modelo de impresora

Todas las impresoras deben representarse mediante un modelo común.

Ejemplo:

```go
type Printer struct {
    ID           string
    Name         string
    Manufacturer string
    Model        string
    Connection   ConnectionType
    Address      string
    Status       PrinterStatus
    Protocol     PrinterProtocol
    IsDefault    bool
}
```

Tipos de conexión:

```text
usb
serial
network
system
unknown
```

Protocolos:

```text
escpos
cups
windows
raw_tcp
unknown
```

Estados:

```text
unknown
ready
printing
offline
paper_out
error
busy
```

El sistema no debe asumir que todas las impresoras proporcionan información detallada de estado.

---

# 13. Identificación de impresoras

Los IDs deben ser lo más estables posible.

Prioridad:

1. Número de serie del dispositivo.
2. Vendor ID + Product ID + ruta.
3. Dirección de red.
4. Identificador del sistema operativo.

No se deben generar IDs aleatorios cada vez que se ejecuta el descubrimiento si existe una identificación estable disponible.

---

# 14. Abstracción de transporte

La aplicación no debe depender directamente de USB, TCP o Serial.

Interfaz propuesta:

```go
type PrinterTransport interface {
    Open(ctx context.Context) error
    Write(ctx context.Context, data []byte) error
    Close() error
}
```

Implementaciones:

```text
USBTransport
SerialTransport
TCPTransport
SystemPrinterTransport
```

Esto permitirá agregar nuevos mecanismos de conexión sin modificar la lógica principal de impresión.

---

# 15. ESC/POS

ESC/POS será el protocolo principal del MVP.

Debe existir un paquete independiente:

```text
pkg/escpos
```

Debe permitir construir trabajos de impresión.

Ejemplo:

```go
receipt := escpos.NewReceipt()

receipt.
    Center().
    Bold().
    Text("SV TECH").
    LineFeed().
    ResetStyle().
    Text("------------------------------").
    LineFeed().
    Text("Servicio soporte     $25.000").
    LineFeed().
    Text("------------------------------").
    Bold().
    Text("TOTAL                $25.000").
    LineFeed().
    Cut()
```

Funciones iniciales:

- Inicializar impresora.
- Texto.
- Negrita.
- Cursiva cuando sea compatible.
- Subrayado.
- Alineación.
- Tamaño de fuente.
- Espaciado.
- Saltos de línea.
- Tablas.
- Imágenes.
- QR.
- Códigos de barras.
- Corte de papel.
- Pulso para cajón de dinero.
- Reset.

---

# 16. Trabajo de impresión

Cada impresión debe representarse como un trabajo.

Ejemplo:

```go
type PrintJob struct {
    ID        string
    PrinterID string
    Payload   []byte
    CreatedAt time.Time
    Status    PrintJobStatus
}
```

Estados:

```text
queued
printing
completed
failed
cancelled
```

Cada trabajo debe recibir un ID único.

Ejemplo:

```text
job_01KXYZ
```

---

# 17. Cola de impresión

El MVP debe utilizar una cola en memoria.

```text
Solicitud HTTP
      │
      ▼
Cola de impresión
      │
      ▼
Print Worker
      │
      ▼
Impresora
```

Para una impresora física debe evitarse enviar simultáneamente varios trabajos.

Versiones posteriores podrán incorporar:

- Cola persistente.
- Reintentos.
- Prioridades.
- Cancelación.
- Historial.
- Recuperación después de reiniciar el agente.

---

# 18. API local

El agente debe exponer una API HTTP local.

Dirección predeterminada:

```text
127.0.0.1:9876
```

El puerto debe poder configurarse.

Por seguridad, el agente debe escuchar únicamente en localhost por defecto.

No debe exponerse a:

```text
0.0.0.0
```

ni a Internet salvo configuración explícita.

---

# 19. Endpoints

## Health

```http
GET /health
```

Respuesta:

```json
{
  "status": "ok",
  "version": "0.1.0"
}
```

---

## Información del agente

```http
GET /api/v1/info
```

Respuesta:

```json
{
  "name": "SV Print Agent",
  "version": "0.1.0",
  "platform": "darwin",
  "architecture": "arm64"
}
```

---

## Listar impresoras

```http
GET /api/v1/printers
```

Respuesta:

```json
{
  "printers": [
    {
      "id": "usb-1234-5678",
      "name": "POS-58",
      "connection": "usb",
      "protocol": "escpos",
      "status": "ready"
    }
  ]
}
```

---

## Detectar impresoras

```http
POST /api/v1/printers/discover
```

Debe ejecutar nuevamente el proceso de descubrimiento.

---

## Obtener impresora

```http
GET /api/v1/printers/{printerID}
```

---

## Imprimir prueba

```http
POST /api/v1/printers/{printerID}/test
```

Debe generar un comprobante estándar de prueba.

---

## Crear trabajo de impresión

```http
POST /api/v1/print
```

Ejemplo:

```json
{
  "printer_id": "usb-1234-5678",
  "format": "escpos",
  "payload": "..."
}
```

Respuesta:

```json
{
  "job_id": "job_01KXYZ",
  "status": "queued"
}
```

---

## Consultar trabajo

```http
GET /api/v1/jobs/{jobID}
```

Respuesta:

```json
{
  "id": "job_01KXYZ",
  "status": "completed"
}
```

---

# 20. WebSocket

Se debe dejar preparado soporte para eventos mediante WebSocket.

Endpoint:

```text
ws://127.0.0.1:9876/api/v1/events
```

Eventos:

```text
printer.discovered
printer.removed
printer.status_changed
print.started
print.completed
print.failed
agent.status
```

Ejemplo:

```json
{
  "event": "printer.status_changed",
  "printer_id": "usb-1234-5678",
  "status": "paper_out"
}
```

Puede considerarse experimental durante el MVP.

---

# 21. Autenticación

La API local debe utilizar autenticación.

Durante la primera inicialización se debe generar un token.

Ejemplo:

```text
SV Print Agent initialized.

Agent ID:
agent_xxxxxxxxx

Token:
xxxxxxxxxxxxxxxxxxxxxxxx
```

Las solicitudes utilizarán:

```http
Authorization: Bearer <token>
```

Los tokens no deben aparecer en logs.

Cuando sea posible, las credenciales deben almacenarse utilizando mecanismos seguros del sistema operativo.

---

# 22. CORS

La API debe implementar restricciones CORS.

No se debe utilizar:

```text
Access-Control-Allow-Origin: *
```

como configuración predeterminada.

Debe existir una lista de orígenes autorizados.

Ejemplo:

```text
https://app.svtech.cl
https://localhost:3000
```

La lista debe ser configurable.

---

# 23. Seguridad

Principios obligatorios:

1. Escuchar en localhost por defecto.
2. Utilizar autenticación.
3. Validar orígenes CORS.
4. Validar todos los payloads.
5. Limitar el tamaño de los trabajos.
6. Limitar el tamaño de la cola.
7. No permitir acceso arbitrario al sistema de archivos.
8. No ejecutar comandos del sistema enviados desde la web.
9. No permitir conexiones de red arbitrarias.
10. No exponer el agente directamente a Internet.
11. No registrar credenciales.
12. No registrar información sensible de clientes.
13. Validar los IDs de impresoras.
14. Evitar SSRF mediante configuraciones de red controladas.

---

# 24. Límites

El agente debe tener límites para evitar abuso o consumo excesivo de recursos.

Valores iniciales sugeridos:

```text
Tamaño máximo de payload: 5 MB
Tamaño máximo de cola: 100 trabajos
Trabajos simultáneos por impresora: 1
```

Estos valores deben poder configurarse.

---

# 25. CLI

El MVP no requiere una interfaz gráfica.

Debe proporcionar una CLI administrativa.

Ejemplos:

```bash
sv-print status

sv-print printers

sv-print discover

sv-print test <printer-id>

sv-print print <printer-id> receipt.json

sv-print config

sv-print logs

sv-print doctor

sv-print version
```

Ejemplo:

```bash
sv-print printers
```

Resultado:

```text
SV Print Agent

Impresoras
────────────────────────────────────────────

ID              NOMBRE          CONEXIÓN

usb-1234        POS-58          USB
net-001         Epson TM20      TCP 192.168.1.50:9100

Estado
────────────────────────────────────────────

POS-58          READY
Epson TM20      OFFLINE
```

---

# 26. ¿UI gráfica o CLI?

La primera versión debe utilizar **CLI y no UI gráfica**.

La web será la interfaz principal para el usuario final.

La CLI estará destinada a:

- Instalación.
- Diagnóstico.
- Configuración.
- Pruebas.
- Soporte técnico.
- Resolución de problemas.

El núcleo del agente debe permanecer completamente independiente de la CLI.

---

# 27. Posible UI futura

Si la experiencia de usuario demuestra que la configuración resulta compleja, se podrá incorporar una interfaz gráfica.

La arquitectura futura podría ser:

```text
              SV Print UI
                  │
               Tauri
                  │
                  ▼
            SV Print Core
                  │
        ┌─────────┼─────────┐
        │         │         │
       USB      Serial      TCP
```

La UI no debe contener la lógica principal de impresión.

El núcleo Go debe seguir funcionando de manera independiente.

---

# 28. Configuración

La configuración debe almacenarse fuera del binario.

Ejemplo:

```yaml
server:
  host: 127.0.0.1
  port: 9876

security:
  allowed_origins:
    - https://app.svtech.cl

printers:
  - id: usb-1234-5678
    name: Caja 1
    protocol: escpos

logging:
  level: info
```

La ubicación del archivo debe seguir las convenciones de cada sistema operativo.

---

# 29. Logs

Se utilizarán logs estructurados.

Se recomienda:

```text
log/slog
```

Niveles:

```text
debug
info
warn
error
```

Ejemplo:

```text
INFO printer discovered
    printer_id=usb-1234
    connection=usb
    model=POS-58
```

Nunca deben registrarse:

- Tokens.
- Contraseñas.
- Credenciales.
- Payloads completos que puedan contener información sensible.
- Datos personales innecesarios.

---

# 30. Diagnóstico

Debe existir:

```bash
sv-print doctor
```

Ejemplo:

```text
SV Print Doctor

✓ Sistema operativo detectado
✓ Configuración válida
✓ API local disponible
✓ Subsistema USB disponible
✓ Red disponible
✓ Impresora detectada

Impresoras:

✓ POS-58
```

Si existe un problema:

```text
✗ POS-58

Motivo:
No fue posible establecer conexión.

Acciones sugeridas:
- Verificar cable USB.
- Verificar alimentación.
- Verificar permisos.
- Ejecutar nuevamente la detección.
```

---

# 31. Instalación como servicio

El agente debe iniciarse automáticamente con el sistema operativo.

## Windows

Utilizar:

```text
Windows Service
```

## Linux

Utilizar:

```text
systemd
```

## macOS

Utilizar:

```text
launchd
```

---

# 32. Ciclo de vida

La CLI debe permitir:

```bash
sv-print service start

sv-print service stop

sv-print service restart

sv-print service status
```

La implementación específica de cada sistema operativo debe estar aislada de la lógica principal.

---

# 33. Monitoreo de impresoras

El agente debe poder monitorear impresoras configuradas.

Ejemplo:

```text
READY
  │
  ▼
PRINTING
  │
  ▼
READY
```

Error:

```text
READY
  │
  ▼
OFFLINE
  │
  ▼
READY
```

El intervalo de monitoreo debe ser configurable.

La disponibilidad de estados como `paper_out` dependerá de las capacidades de cada impresora y protocolo.

---

# 34. Manejo de errores

Los errores deben utilizar códigos estructurados.

Ejemplos:

```text
PRINTER_NOT_FOUND
PRINTER_OFFLINE
PRINTER_BUSY
PRINTER_PAPER_OUT
PRINTER_CONNECTION_FAILED
PRINT_FAILED
INVALID_PAYLOAD
UNAUTHORIZED
ORIGIN_NOT_ALLOWED
PAYLOAD_TOO_LARGE
UNSUPPORTED_PROTOCOL
```

Respuesta:

```json
{
  "error": {
    "code": "PRINTER_OFFLINE",
    "message": "Printer is currently offline."
  }
}
```

La aplicación web será responsable de traducir los mensajes al idioma del usuario mediante i18n.

---

# 35. Pruebas

## Unitarias

Se deben cubrir:

- Entidades.
- Validaciones.
- ESC/POS.
- Cola.
- Jobs.
- Identificación de impresoras.
- Configuración.
- Autenticación.
- Manejo de errores.

## Integración

Se deben probar:

- API HTTP.
- WebSocket.
- TCP.
- Serial cuando sea posible.
- Flujo completo de impresión.

## Hardware

Las pruebas reales con impresoras deben mantenerse separadas de las pruebas unitarias y ejecutarse sobre hardware físico.

---

# 36. CI/CD

El proyecto debe utilizar GitHub Actions.

Build matrix:

```text
windows/amd64
windows/arm64

linux/amd64
linux/arm64

darwin/amd64
darwin/arm64
```

El pipeline debe ejecutar:

```text
go fmt
go vet
go test
go build
```

Se podrán agregar posteriormente:

- golangci-lint.
- análisis de seguridad.
- pruebas de integración.
- generación automática de releases.

---

# 37. Artefactos de lanzamiento

Cada release debe generar:

```text
sv-print-windows-amd64.exe
sv-print-windows-arm64.exe

sv-print-linux-amd64
sv-print-linux-arm64

sv-print-darwin-amd64
sv-print-darwin-arm64
```

Posteriormente:

```text
Windows Installer
.deb
.rpm
AppImage
.dmg
.pkg
```

Cada release debe generar checksums.

---

# 38. Versionado

Se utilizará:

```text
Semantic Versioning
```

Formato:

```text
MAJOR.MINOR.PATCH
```

Ejemplo:

```text
0.1.0
```

Durante el desarrollo inicial el proyecto permanecerá en versión `0.x`.

---

# 39. Versionado de API

Todos los endpoints públicos utilizarán:

```text
/api/v1/
```

Ejemplo:

```text
/api/v1/printers
```

Los cambios incompatibles requerirán:

```text
/api/v2/
```

Los cambios compatibles deben mantenerse dentro de la misma versión de API.

---

# 40. Flujo principal

El flujo esperado para un usuario será:

```text
1. Instalar SV Print Agent
             │
             ▼
2. Agent inicia automáticamente
             │
             ▼
3. Detecta impresoras
             │
             ▼
4. Aplicación web detecta Agent
             │
             ▼
5. Aplicación obtiene impresoras
             │
             ▼
6. Usuario selecciona impresora
             │
             ▼
7. Aplicación genera trabajo
             │
             ▼
8. Aplicación envía trabajo
             │
             ▼
9. SV Print Agent procesa cola
             │
             ▼
10. Agent envía ESC/POS
             │
             ▼
11. Impresora imprime
             │
             ▼
12. Agent informa resultado
```

---

# 41. Ejemplo de integración

La aplicación web consulta:

```http
GET http://127.0.0.1:9876/health
```

Después:

```http
GET http://127.0.0.1:9876/api/v1/printers
```

Obtiene:

```json
{
  "printers": [
    {
      "id": "usb-1234",
      "name": "Caja 1",
      "status": "ready"
    }
  ]
}
```

Finalmente:

```http
POST http://127.0.0.1:9876/api/v1/print
```

El agente:

```text
Validar solicitud
        │
        ▼
Autenticar
        │
        ▼
Validar impresora
        │
        ▼
Crear Print Job
        │
        ▼
Agregar a cola
        │
        ▼
Abrir conexión
        │
        ▼
Enviar ESC/POS
        │
        ▼
Cerrar conexión
        │
        ▼
Actualizar estado
```

---

# 42. Consideraciones de compatibilidad

La compatibilidad con una impresora no debe determinarse únicamente por su marca.

Se debe considerar:

```text
Marca
Modelo
Tipo de conexión
Protocolo
Comandos ESC/POS soportados
Sistema operativo
Driver requerido
```

El proyecto debe mantener una estrategia de compatibilidad progresiva:

```text
ESC/POS estándar
       ↓
Extensiones comunes
       ↓
Comandos específicos de fabricante
```

Los comandos específicos de fabricante deben mantenerse aislados.

---

# 43. Roadmap

## Fase 1 — MVP

```text
Go
CLI
HTTP API
USB
TCP
ESC/POS
Windows
macOS
Linux
```

## Fase 2

```text
Serial
WebSocket
Monitoreo
Configuración persistente
Instaladores
Servicios del sistema
```

## Fase 3

```text
CUPS
Windows Print Spooler
Descubrimiento avanzado
SNMP
mDNS / DNS-SD
```

## Fase 4

```text
UI gráfica
System Tray
Configuración amigable
```

## Fase 5

```text
SV Print Cloud
Administración de dispositivos
Configuración remota
Licenciamiento
Características empresariales
```

---

# 44. Principios de diseño

El proyecto debe respetar:

- Código idiomático de Go.
- Interfaces pequeñas.
- Separación de responsabilidades.
- Inyección de dependencias cuando corresponda.
- Dependencia invertida.
- Código multiplataforma.
- Aislamiento del código específico del sistema operativo.
- Manejo explícito de errores.
- Uso de `context.Context`.
- Operaciones seguras para concurrencia.
- Evitar estado global mutable.
- Código testeable.
- Configuración segura por defecto.
- Dependencias externas mínimas.

---

# 45. Resultado esperado

El resultado final debe ser un agente pequeño y confiable que permita a cualquier aplicación web compatible realizar:

```text
Detectar impresora
       ↓
Seleccionar impresora
       ↓
Enviar trabajo
       ↓
Imprimir
       ↓
Consultar resultado
```

sin que la aplicación tenga que conocer los detalles específicos de:

```text
Windows
macOS
Linux
USB
Serial
Ethernet
Wi-Fi
ESC/POS
```

La API local debe constituir el contrato principal de integración.

La CLI debe constituir la herramienta principal de administración técnica.

La UI gráfica queda como una posible capa futura y no debe formar parte del núcleo del proyecto.
