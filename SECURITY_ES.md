> **Language:** Español | [English](SECURITY.md)

# Política de Seguridad

Nos tomamos en serio la seguridad de **SV Printer**. Dado que este agente maneja acceso a red local, comunicación con impresoras y verificación de licencias, proteger su despliegue es nuestra máxima prioridad.

---

## 1. Versiones Soportadas

Soportamos y parchamos activamente problemas de seguridad en las siguientes versiones:

| Versión | Soportada         |
| ------- | ----------------- |
| 0.1.x   | :white_check_mark:|
| < 0.1.0 | :x:               |

> Política: **solo la última versión menor.** Los fixes de seguridad llegan a la versión menor actual. Se espera que versiones menores anteriores actualicen.

---

## 2. Reportar una Vulnerabilidad

**No abras un issue público en GitHub para vulnerabilidades de seguridad.**

Si descubres una vulnerabilidad de seguridad, repórtala de forma privada:

1. **GitHub Security Advisory:** ve a la pestaña **Security** del repositorio y selecciona **Advisories → New draft advisory**. Esto nos permite discutir y parchear el issue en privado.
2. **Email:** alternativamente, contacta al equipo SVTech en `security@svtech.software`.

Responderemos a tu reporte dentro de **48 horas** y proporcionaremos una línea de tiempo detallada para un parche.

---

## 3. Modelo de Seguridad

SV Printer sigue estos principios de seguridad:

- **Vinculación local únicamente:** el agente se vincula a `127.0.0.1` por defecto y nunca a `0.0.0.0`, previniendo acceso remoto a menos que se configure explícitamente.
- **Autenticación por token:** todos los endpoints `/api/*` requieren `Authorization: Bearer <token>`. El token se genera en la primera ejecución y se almacena en el archivo de configuración.
- **Restricción CORS:** solo los orígenes listados explícitamente en `-origin` / `allowed_origins` pueden llamar a la API desde un navegador.
- **Verificación de licencia Ed25519:** las licencias están firmadas criptográficamente y se verifican offline. La clave privada de firma nunca se commitea al repositorio.
- **Vinculación a dispositivo:** las licencias pueden vincularse a una máquina específica mediante verificación de huella.

---

## 4. Higiene de Secretos

- **Claves privadas:** la clave privada de firma Ed25519 (`*.priv`, `.sv-license.priv`) está **gitignoreada** y nunca debe ser commiteada o compartida.
- **Archivo de configuración:** el token de autenticación se almacena en el archivo JSON de configuración. Asegúrate de que el directorio de configuración tenga permisos de archivo apropiados.
- **Sin telemetría:** SV Printer no hace phone home ni transmite ningún dato externamente. Toda la verificación de licencias es offline.

---

## 5. Alcance

Esta política de seguridad cubre el binario del agente SV Printer, su API HTTP y el sistema de verificación de licencias. No cubre:

- Dependencias de terceros (reportarlas upstream).
- Seguridad física de impresoras.
- Configuraciones de red más allá del vinculado del agente.

---

## 6. Actualizaciones

Los fixes de seguridad se publican como versiones menores (ej. `0.1.1`) y se anuncian en el [CHANGELOG](./CHANGELOG.md). Recomendamos siempre ejecutar la última versión.
