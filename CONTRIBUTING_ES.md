# Contribuir a SV Printer

Gracias por tu interés en contribuir a **SV Printer**. Esta guía describe el flujo de trabajo y los estándares.

> **¿Solo quieres instalar SV Printer?** Consulta la [Guía de Integración](./documentation/integration_ES.md) para instrucciones de compilación.

---

## 1. Estrategia de Ramas y Flujo de Trabajo

Seguimos un flujo de trabajo simple con ramas de features:

1. **Nomenclatura de ramas:** crea una rama desde `main` con prefijos descriptivos:
   - `feat/nombre-feature` para nuevas features.
   - `fix/nombre-bug` para corrección de bugs.
   - `documentation/nombre-doc` para actualizaciones de documentación.
   - `refactor/nombre-refactor` para limpieza de código.
2. **Pull Requests:** cuando tus cambios estén listos y probados:
   - Empuja tu rama a GitHub.
   - Abre un Pull Request (PR) contra `main`.
   - Asegúrate de que todos los tests pasen antes de solicitar revisión.

> **Nota para colaboradores externos:** todos los cambios deben llegar a través de un pull request. La Integración Continua (CI) ejecuta `gofmt`, `go vet`, `go test` y una verificación de compilación en Linux, macOS y Windows, y **debe pasar en verde** antes de que tu PR pueda ser mergeado.

---

## 2. Estándares de Mensajes de Commit

Este repositorio aplica estrictamente la especificación **[Conventional Commits](https://www.conventionalcommits.org/)**. **Todos los mensajes de commit deben estar escritos en inglés.**

Formato: `<type>(<scope>): <description>`

- **Tipos permitidos:** `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `build`, `ci`, `chore`.
- **Scope:** opcional pero recomendado (ej. `api`, `escpos`, `config`, `licensing`).
- **Ejemplos:**
  - `feat(api): add receipt structured endpoint`
  - `fix(ws): subscribe before Accept to eliminate EventsHandler race`
  - `docs(readme): redesign with hero, badges and Mermaid`

---

## 3. Entorno de Desarrollo

Asegúrate de tener Go 1.26+ instalado (ver [`go.mod`](./go.mod)).

### Comandos Clave

```bash
# Compilar
go build ./cmd/sv-printer

# Ejecutar tests unitarios
go test ./...

# Formatear código
gofmt -l .

# Análisis estático
go vet ./...
```

---

## 4. Puerta de Pre-Commit (Requerida)

Antes de presentar cualquier mensaje de commit, **debes** ejecutar la barrera completa y corregir cualquier fallo:

```bash
gofmt -l .
go vet ./...
go test ./...
```

Si alguno falla, corrige los issues primero, re-ejecuta, y solo entonces presenta el mensaje de commit. Esto previene ciclos de CI desperdiciados y ramas principales rotas.

---

## 5. Convenciones del Proyecto

- **TDD:** escribir tests después de cada tarea y mantenerlos pasando.
- **Biblioteca estándar primero:** evitar paquetes externos cuando la stdlib de Go alcanza.
- **Spec-driven:** los cambios de comportamiento y arquitectura pasan por el flujo de spec de `sv-memory` antes de la implementación.
- **Sin automatización de git:** nunca ejecutar `git add`, `git commit` o `git push` de forma autónoma — el desarrollador maneja estos comandos manualmente.
- **Protocolo de memoria:** usar las herramientas de `sv-memory` (buscar, guardar, grafo) durante el flujo de trabajo.

---

## 6. Releases

Los releases se automatizan al empujar un tag (`vX.Y.Z`) y son atómicos y reproducibles:

1. **`goreleaser`** compila los archivos CLI multiplataforma y `checksums.txt`, y abre el release como **draft**.
2. **macOS** compila los bundles `.app` por arquitectura y universal; **Windows** compila el instalador NSIS. Ambos suben sus artefactos al draft.
3. Un job final etiqueta los assets (app GUI vs CLI headless) y publica el draft como último release.

Pautas:

- **Nunca re-ejecutes el workflow sobre un tag existente.** El job `guard` aborta si ya existe un release — sube la versión y empuja un tag nuevo.
- Las builds usan `-trimpath` y un `mod_timestamp` determinista, así que el mismo commit produce checksums idénticos.
- Si un release falla a mitad de camino, elimina el draft release antes de volver a etiquetar.

---

## 7. Licencia

SV Printer se distribuye bajo la [Licencia Comercial 1.1](./LICENSE) (BSL 1.1). Al contribuir, aceptas que tus contribuciones se licenciarán bajo los mismos términos.
