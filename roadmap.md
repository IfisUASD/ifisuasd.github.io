# **Roadmap del Proyecto Web: Instituto de Física**

Este documento sirve como la lista maestra de tareas (To-Do List) para el desarrollo del portal web. Está organizado cronológicamente por dependencias.

## **🏁 Fase 0: Configuración y Andamiaje (Scaffolding)**

Estado: ✅ Completado  
Objetivo: Tener un entorno de desarrollo local funcional donde se pueda ejecutar un "Hola Mundo" con Go, Templ y Tailwind.

* [x] **Inicialización del Proyecto**  
  * [x] Crear repositorio en GitHub.  
  * [x] go mod init github.com/tu-usuario/instituto-fisica.  
  * [x] Configurar .gitignore.  
* [x] **Instalación de Herramientas**  
  * [x] Instalar templ.  
  * [x] Inicializar NPM.  
  * [x] Instalar Tailwind \+ DaisyUI.  
  * [x] Generar configuración Tailwind.  
* [x] **Configuración Base**  
  * [x] Configurar tailwind.config.js.  
  * [x] Configurar temas.  
  * [x] Crear estructura de carpetas.  
* [x] **Automatización (Makefile)**  
  * [x] Crear Makefile.

## **🧠 Fase 1: Motor de Datos (Backend en Go)**

Estado: 🚧 En Progreso  
Objetivo: Que Go sea capaz de leer todos los archivos Markdown, YAML y BibTeX y relacionarlos en memoria antes de generar una sola línea de HTML.

* [x] **Modelado de Datos (Structs)**  
  * [x] Definir struct Person (con campos calculados para papers/proyectos).  
  * [x] Definir struct Project.  
  * [x] Definir struct Paper (campos BibTeX \+ x-orcids).  
  * [x] Definir struct NewsItem y BlogPost.  
* [x] **Parsers (Lectura de Archivos)**  
  * [x] Implementar parser de **Markdown \+ Frontmatter** (usando goldmark \+ yaml).  
  * [x] Implementar parser de **BibTeX** (leyendo campos custom x-\*).  
  * [x] Implementar lógica de lectura recursiva de directorios.  
* [x] **Lógica de Internacionalización (i18n)**  
  * [x] Crear función para detectar pares de archivos (.es.md y .en.md).  
  * [x] Cargar archivo i18n.json (diccionario de etiquetas UI).  
* [x] **Algoritmo de Vinculación (The linker)**  
  * [x] **TEST CRÍTICO:** Crear función LinkData(db \*Database).  
  * [x] Lógica: Asignar Papers a Personas mediante x-orcids.  
  * [x] Lógica: Asignar Papers a Proyectos mediante x-project.  
  * [x] Lógica: Asignar Integrantes a Proyectos mediante IDs en YAML.  
* [x] **Pruebas Unitarias (Fase 1\)**  
  * [x] go test: Verificar que el parser BibTeX no falla con datos vacíos.  
  * [x] go test: Verificar que las relaciones se crean correctamente en memoria.

## **🎨 Fase 2: Sistema de Diseño y UI (Frontend)**

**Objetivo:** Definir la apariencia visual y los componentes reutilizables.

* [x] **Layout Base**  
  * [x] Crear templates/layouts/base.templ.  
  * [x] Implementar \<head\> con metadatos SEO básicos.  
  * [x] Integrar script de **Dark Mode**.  
* [x] **Navegación (Navbar)**  
  * [x] Diseñar Navbar responsiva (DaisyUI).  
  * [x] Implementar menú hamburguesa.  
  * [x] Implementar **Selector de Idioma**.  
  * [x] Implementar **Toggle Dark/Light**.  
* [x] **Componentes Reutilizables**  
  * [x] Card, Avatar, Badge, PublicationRow.  
* [x] **Footer**  
  * [x] Diseñar pie de página.

## **📄 Fase 3: Generación de Páginas (Contenido)**

**Objetivo:** Generar los archivos HTML finales combinando los Datos (Fase 1\) con la UI (Fase 2).

* [x] **Páginas Estáticas Simples** (index.html, Historia, Transparencia).  
* [x] **Generadores Dinámicos** (Gente, Proyectos, Publicaciones, Blog).  
* [x] **Internacionalización de Páginas**.

## **🔍 Fase 4: Interactividad y Pulido**

**Objetivo:** Añadir funcionalidades que requieren JavaScript mínimo.

* [x] **Buscador / Filtros** (JS puro).  
* [x] **Mejoras UX** (Breadcrumbs, Lazy loading, OpenGraph).

## **🧪 Fase 5: Aseguramiento de Calidad (QA)**

**Objetivo:** "Zero Broken Pages".

* [x] **Configuración de Herramientas** (htmltest, Playwright).  
* [x] **Link Integrity Testing**.  
* [x] **E2E Testing**.

## **🚀 Fase 6: Despliegue y CI/CD**

**Objetivo:** Automatizar todo en GitHub.

* [x] **GitHub Actions** (deploy.yml).  
* [x] **Documentación** (README.md).