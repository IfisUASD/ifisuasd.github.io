# Instituto de Física - UASD (Web)

Este repositorio contiene el código fuente del sitio web oficial del Instituto de Física de la Universidad Autónoma de Santo Domingo (UASD).

## 🛠️ Tech Stack (GOTH Stack)

El sitio está construido utilizando una arquitectura estática moderna y eficiente:

*   **Go (Golang) v1.25+**: Lenguaje principal para la generación del sitio.
*   **Templ**: Motor de plantillas type-safe para Go (usado como `go tool`).
*   **TailwindCSS v4**: Framework de utilidades CSS para el diseño.
*   **DaisyUI**: Librería de componentes para Tailwind.
*   **Playwright**: Pruebas End-to-End (E2E).
*   **GitHub Actions**: CI/CD para despliegue automático.

## 🚀 Requisitos Previos

Asegúrate de tener instaladas las siguientes herramientas:

*   [Go 1.25+](https://go.dev/dl/)
*   [Node.js v20+](https://nodejs.org/) (para Tailwind y Playwright)
*   Make (para ejecutar comandos de construcción)

## 📦 Instalación

1.  **Clonar el repositorio:**

    ```bash
    git clone https://github.com/IfisUASD/ifisuasd.github.io.git
    cd ifisuasd.github.io
    ```

2.  **Instalar dependencias:**

    ```bash
    make deps           # Instala dependencias de NPM
    go mod download     # Descarga módulos de Go
    make install-htmltest # Instala herramienta de testing HTML
    npx playwright install # Instala navegadores para E2E
    ```

## 🏃‍♂️ Ejecución Local

### Desarrollo

Para desarrollo, puedes usar el servidor local que sirve el contenido generado:

1.  **Generar el sitio:**
    ```bash
    make build
    ```

2.  **Iniciar el servidor:**
    ```bash
    make serve
    ```
    El sitio estará disponible en `http://localhost:8180`.

### Comandos Útiles (Makefile)

*   `make build`: Genera todo el sitio (Templates, CSS, HTML, Assets).
*   `make css`: Compila solo el CSS con Tailwind.
*   `make templ`: Genera el código Go de las plantillas `.templ`.
*   `make clean`: Limpia el directorio `output`.
*   `make clean-deps`: Elimina `node_modules` para una instalación limpia.

## 🧪 Testing

El proyecto incluye varios niveles de pruebas para asegurar la calidad:

*   **Unit Tests (Go):**
    ```bash
    go test ./...
    ```

*   **HTML Validation (htmltest):**
    Verifica enlaces rotos y estructura HTML.
    ```bash
    make test-html
    ```

*   **End-to-End Tests (Playwright):**
    Prueba la funcionalidad del sitio en un navegador real.
    ```bash
    make test-e2e
    ```

## 📂 Estructura del Proyecto

*   `/content`: Archivos Markdown y BibTeX (Fuente de Verdad).
    *   `/people`: Perfiles de investigadores.
    *   `/projects`: Proyectos de investigación.
    *   `/news`: Noticias.
    *   `/references`: Archivo `.bib` con las publicaciones.
*   `/templates`: Componentes y páginas `.templ`.
*   `/assets`: Imágenes, CSS y JS estáticos.
*   `/cmd`: Puntos de entrada de la aplicación (Generador, Servidor).
*   `/internal`: Lógica de negocio (Parsers, Linker, Tipos).
*   `/output`: Sitio estático generado (listo para desplegar).

## 🚢 Despliegue

El despliegue es automático a través de **GitHub Actions**.

1.  Al hacer push a la rama `main`, se dispara el workflow.
2.  El sitio se construye y valida.
3.  Se despliega automáticamente a la rama `gh-pages`.
4.  El sitio es accesible en [ifis.edu.do](https://ifis.edu.do).

## 📝 Licencia

Este proyecto es propiedad del Instituto de Física de la UASD.
