# Guía de Protección del Repositorio (GitHub)

Antes de hacer este repositorio público, es estrictamente necesario habilitar reglas de protección en GitHub para evitar escrituras accidentales en la rama principal y garantizar la calidad del código mediante integración continua.

## Pasos para proteger la rama `main`

Dirígete a tu repositorio en GitHub y haz clic en la pestaña **Settings** (Configuración).
En el menú lateral izquierdo, bajo la sección "Code and automation", haz clic en **Branches** (Ramas).
En la sección "Branch protection rules", haz clic en **Add branch protection rule**.

### Configuración de la regla

1. **Branch name pattern**: Escribe `main`.

2. **Protect matching branches** (Habilita las siguientes opciones):
   - [x] **Require a pull request before merging**
     - [x] **Require approvals**: Configúralo en al menos `1`.
     - [x] **Require review from Code Owners**: Esto es crítico. Hará que GitHub lea el archivo `.github/CODEOWNERS` y exija la aprobación explícita tuya (o de tu equipo) antes de que el código pueda integrarse.
   
   - [x] **Require status checks to pass before merging**
     - [x] **Require branches to be up to date before merging**
     - Busca y selecciona los flujos de CI críticos en la barra de búsqueda (ej. `build`, `test`, `lint`, dependiendo de los nombres de tus GitHub Actions). Esto asegura que no se pueda mezclar código si las pruebas fallan.

   - [x] **Do not allow bypassing the above settings**: (Opcional, pero recomendado). Si lo habilitas, ni siquiera tú como administrador podrás saltarte los tests o la revisión. Te protege de tus propios errores.

3. **Rules applied to everyone including administrators**:
   - [x] **Restrict who can push to matching branches**: Asegúrate de que no haya nadie seleccionado aquí, o restablécelo solo a administradores muy específicos. En la práctica, *nadie* debería hacer push directo a `main`.
   - [x] **Allow force pushes**: Asegúrate de que esté **DESMARCADO**. Nunca se debe reescribir la historia pública de `main`.
   - [x] **Allow deletions**: Asegúrate de que esté **DESMARCADO**.

4. Finalmente, haz clic en **Create** (Crear) o **Save changes** (Guardar cambios).

## Impacto de esta configuración

Al hacer público el repositorio con estas reglas:
* Cualquier persona externa tendrá que hacer un **Fork** y crear un **Pull Request**.
* Tú (y tu equipo) tendrán que revisar el código y aprobarlo.
* Las GitHub Actions ejecutarán las pruebas automáticamente.
* Nadie, ni siquiera accidentalmente mediante scripts locales o comandos mal escritos (`git push origin main --force`), podrá destruir el código de producción.

## Transición a Equipos (Opción A futura)

Actualmente, el archivo `.github/CODEOWNERS` te asigna como individuo (`* @msandoval-93`). Cuando el equipo de SV Tech SpA crezca, debes:
1. Crear un equipo en la organización de GitHub (ej. `@svtech-code/core`).
2. Cambiar el archivo a `* @svtech-code/core`.
3. Todos los miembros agregados a ese equipo en GitHub heredarán permisos de Code Owner automáticamente sin tocar los archivos del repositorio.
