#!/bin/bash

# Asegurarse de estar en la raíz del proyecto
cd "$(dirname "$0")/.." || exit

# 1. Verificar que exista el archivo de secretos
if [ ! -f ".env.secret" ]; then
    echo "❌ Error: No se encontró el archivo .env.secret"
    echo "Por favor crea un archivo llamado .env.secret en la raíz del proyecto con el siguiente contenido:"
    echo "SV_LICENSE_KEY=tu_clave_privada_aqui"
    exit 1
fi

# Cargar la clave privada
source .env.secret

if [ -z "$SV_LICENSE_KEY" ]; then
    echo "❌ Error: SV_LICENSE_KEY no está definida en .env.secret"
    exit 1
fi

# 2. Solicitar datos al usuario
echo "====================================="
echo "🖨️  Generador de Licencias SV PRINTER"
echo "====================================="
read -p "👤 Nombre del Cliente o Empresa: " CUSTOMER_NAME
read -p "💻 ID del Equipo (Device ID): " DEVICE_ID

if [ -z "$CUSTOMER_NAME" ] || [ -z "$DEVICE_ID" ]; then
    echo "❌ Error: El nombre y el ID son obligatorios."
    exit 1
fi

# 3. Crear carpeta de salida si no existe
mkdir -p licenses

# 4. Formatear el nombre del archivo (quitar espacios)
SAFE_NAME=$(echo "$CUSTOMER_NAME" | tr ' ' '_')
OUTPUT_FILE="licenses/${SAFE_NAME}_license.key"

# 5. Generar la licencia
echo "Generando licencia PRO para $CUSTOMER_NAME..."

go run ./cmd/sv-license sign \
  -customer "$CUSTOMER_NAME" \
  -tier full \
  -features raw_print,websocket \
  -fingerprint "$DEVICE_ID" \
  > "$OUTPUT_FILE"

if [ $? -eq 0 ]; then
    echo "✅ ¡Licencia generada con éxito!"
    echo "📂 Archivo guardado en: $OUTPUT_FILE"
    echo "Envíale este archivo al cliente."
else
    echo "❌ Error al generar la licencia."
fi
