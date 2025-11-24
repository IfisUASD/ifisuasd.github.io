# Makefile para el Instituto de Física

# Nombres de archivos y carpetas
BINARY_NAME=site_gen
OUTPUT_DIR=./output
CSS_INPUT=./assets/css/input.css
CSS_OUTPUT=$(OUTPUT_DIR)/assets/css/styles.css
WASM_SRC=./cmd/apps/qr
WASM_OUTPUT=$(OUTPUT_DIR)/assets/wasm

# Comandos principales
all: clean build

# 1. Limpiar carpeta de salida
clean:
	@echo "🧹 Limpiando..."
	@rm -rf $(OUTPUT_DIR)
	@rm -f $(BINARY_NAME)

# 2. Generar código Go desde Templ
templ:
	@echo "🔥 Generando templates..."
	@templ generate

# 2.1 Compilar Aplicación WebAssembly (WASM)
wasm:
	@echo "⚙️  Compilando Aplicaciones WASM..."
	@mkdir -p $(WASM_OUTPUT)
	@GOOS=js GOARCH=wasm go build -o $(WASM_OUTPUT)/qr.wasm $(WASM_SRC)/main.go $(WASM_SRC)/qr_logic.go
	@GOOS=js GOARCH=wasm go build -o $(WASM_OUTPUT)/markdown.wasm cmd/apps/markdown/main.go

# 3. Generar CSS con Tailwind
css:
	@echo "🎨 Compilando CSS..."
	@# Aseguramos que el directorio exista
	@mkdir -p $(dir $(CSS_OUTPUT))
	@npx @tailwindcss/cli -i $(CSS_INPUT) -o $(CSS_OUTPUT) --minify

# 4. Compilar y Ejecutar el Generador (Go)
build: templ copy-assets wasm
	@echo "🚀 Construyendo sitio..."
	@go build -o $(BINARY_NAME) cmd/builder/main.go
	@# Ejecutamos el binario para crear los HTMLs
	@./$(BINARY_NAME)
	@# Generamos el CSS al final para que Tailwind vea las clases en los .html generados si fuera necesario (aunque ve los .templ)
	@make css
	@echo "✅ ¡Sitio generado en $(OUTPUT_DIR)!"

copy-assets:
	@echo "📂 Copiando assets estáticos..."
	@mkdir -p $(OUTPUT_DIR)/assets
	@cp -r assets/images $(OUTPUT_DIR)/assets/ 2>/dev/null || :
	@cp -r assets/js $(OUTPUT_DIR)/assets/ 2>/dev/null || :

# Desarrollo: Escucha cambios (requiere 'air' instalado, opcional)
dev:
	@air

# 5. Servidor de Desarrollo
serve:
	@go run cmd/server/main.go

# 6. Testing
install-htmltest:
	@echo "📥 Instalando htmltest..."
	@go install github.com/wjdp/htmltest@latest

test-html:
	@echo "🔍 Ejecutando htmltest..."
	@htmltest

test-e2e:
	@echo "🎭 Ejecutando Playwright..."
	@npx playwright test

# 7. Gestión de Dependencias (NPM)
deps:
	@echo "📦 Instalando dependencias de NPM..."
	@npm install

clean-deps:
	@echo "🗑️  Eliminando node_modules y package-lock.json..."
	@rm -rf node_modules package-lock.json

.PHONY: all clean templ css build dev serve install-htmltest test-html test-e2e deps clean-deps
