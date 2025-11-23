# **Instrucciones del Proyecto: Instituto de Física Web (GOTH Stack)**

Este documento define el contexto, las convenciones y las reglas de arquitectura para cualquier agente de IA (Gemini, Copilot, etc.) que asista en el desarrollo de este proyecto.

## **1\. Tech Stack (El "GOTH" Stack Estático)**

- **Lenguaje:** Go (Golang) v1.23+
- **Templating:** templ (https://templ.guide) \- **NUNCA** generar HTML puro en strings de Go, usar siempre componentes .templ.
- **Estilos:** TailwindCSS \+ DaisyUI.
- **Build System:** Makefile (targets: build, css, templ).
- **Testing:** go test (Standard Library).

## **2\. Filosofía de Arquitectura**

### **A. Fuente Única de Verdad (SSOT)**

- La "base de datos" son archivos planos en /content.
- **BibTeX es Rey:** La relación _Paper-Persona_ y _Paper-Proyecto_ se define **EXCLUSIVAMENTE** en el archivo .bib mediante campos personalizados (x-orcids, x-project). No se deben duplicar listas de papers en los archivos Markdown.

### **B. Identificadores vs. Slugs**

- **Humanos (URLs/Archivos):** Usamos slugs legibles.
  - Bien: content/people/vladimir-perez.md
  - Mal: content/people/0000-0002-1825-0097.md
- **Máquinas (Lógica Interna):** Usamos IDs robustos para vincular.
  - **Personas:** orcid (Ej: 0000-0002-1825-0097).
  - **Proyectos:** project_id (Ej: FONDOCYT-2024-10).
  - **Papers:** doi (Ej: 10.1103/PhysRev.1.1).

### **C. Internacionalización (i18n)**

- Sufijos de archivo: archivo.es.md (Español/Default) vs archivo.en.md (Inglés).
- La lógica de Go debe buscar primero el inglés si se está generando el sitio en inglés; si no existe, hacer fallback o ignorar.

## **3\. Estructura de Datos (Referencia Rápida)**

Al generar código Go, respeta estrictamente los modelos definidos en internal/types/types.go.

### **Person (/content/people/\*.md)**

\---  
orcid: "0000-0000-0000-0000" \# KEY PRINCIPAL  
name: "Nombre Legible"  
role: "Director"  
type: "academic" \# academic | staff | student  
\---

### **Project (/content/projects/\*.md)**

\---  
project_id: "CODIGO-UNICO" \# KEY PRINCIPAL  
title: "Título del Proyecto"  
principal_investigator: "ORCID-DEL-PI"  
coinvestigator: \["ORCID-1", "ORCID-2"\]  
research_assistant: \["ORCID-3"\]  
\---

### **Paper (/content/references/institute.bib)**

@article{key,  
 title \= {...},  
 doi \= {...},  
 x-orcids \= {ORCID-1, ORCID-2}, \<-- VINCULACIÓN  
 x-project \= {CODIGO-PROYECTO} \<-- VINCULACIÓN  
}

## **4\. Reglas de Codificación**

1. **TDD (Test Driven Development):**
   - Antes de crear una función nueva en parsers o linker, crea (o pide) el Test Unitario fallido.
   - Usa mocks de datos (\[\]byte o structs en memoria) en lugar de leer archivos reales en los tests.
2. **Manejo de Errores:**
   - Si un ID de referencia (ej. en x-orcids) no existe en la base de datos, **NO** hacer panic. Loguear un WARNING y continuar. La página no debe romperse por un typo en un ID.
3. **Rutas y Directorios:**
   - Input de contenido: ./content/...
   - Output generado: ./output/...
   - **IMPORTANTE:** Al generar HTML, asegurar que las rutas de assets (CSS/Imágenes) sean relativas o absolutas correctas según el entorno (GitHub Pages suele estar en un subdirectorio /repo-name/).

## **5\. Comandos Frecuentes**

- **Regenerar todo:** make build
- **Solo CSS:** make css
- **Correr Tests:** go test ./...
