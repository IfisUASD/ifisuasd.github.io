(() => {
  const go = new Go();
  let wasmReady = false;
  let editor = null;

  // Debounce para no renderizar en cada tecla pulsada
  function debounce(fn, delay) {
    let timer;
    return (...args) => {
      clearTimeout(timer);
      timer = setTimeout(() => fn.apply(this, args), delay);
    };
  }

  async function loadWasm() {
    try {
      const response = await fetch("/assets/wasm/markdown.wasm");
      if (!response.ok) throw new Error("Failed to fetch wasm");

      let result;
      if (WebAssembly.instantiateStreaming) {
        result = await WebAssembly.instantiateStreaming(
          response,
          go.importObject
        );
      } else {
        const buffer = await response.arrayBuffer();
        result = await WebAssembly.instantiate(buffer, go.importObject);
      }

      go.run(result.instance);
      wasmReady = true;

      // Ocultar loader y mostrar editor
      document.getElementById("md-loader").classList.add("hidden");
      document.getElementById("md-app").classList.remove("hidden");

      // Inicializar CodeMirror
      initEditor();

      // Renderizado inicial
      updatePreview();
    } catch (err) {
      console.error(err);
    }
  }

  function initEditor() {
    const container = document.getElementById("editor-container");
    
    // Determinar el tema según el modo actual
    const currentTheme = document.documentElement.getAttribute('data-theme');
    const cmTheme = currentTheme === 'dark' ? 'github-dark' : 'github';
    
    editor = CodeMirror(container, {
      mode: 'gfm',
      theme: cmTheme,
      lineNumbers: true,
      lineWrapping: true,
      autofocus: true,
      viewportMargin: Infinity,
      autoCloseBrackets: true,
      matchBrackets: true,
      placeholder: "*Escribe tu contenido en Markdown aquí...*",
      extraKeys: {
        "Enter": "newlineAndIndentContinueMarkdownList"
      }
    });

    // Actualizar preview cuando cambia el contenido
    editor.on('change', debounce(updatePreview, 150));
    
    // Escuchar cambios de tema
    const observer = new MutationObserver((mutations) => {
      mutations.forEach((mutation) => {
        if (mutation.type === 'attributes' && mutation.attributeName === 'data-theme') {
          const newTheme = document.documentElement.getAttribute('data-theme');
          const newCmTheme = newTheme === 'dark' ? 'github-dark' : 'github';
          editor.setOption('theme', newCmTheme);
        }
      });
    });
    
    observer.observe(document.documentElement, {
      attributes: true,
      attributeFilter: ['data-theme']
    });
  }

  function updatePreview() {
    if (!wasmReady || !editor) return;

    const input = editor.getValue();
    const outputDiv = document.getElementById("mdOutput");

    // Llamada a la función Go
    const html = renderMarkdown(input);

    outputDiv.innerHTML = html;

    // Aplicar resaltado de sintaxis con Highlight.js
    if (window.hljs) {
      outputDiv.querySelectorAll('pre code').forEach((block) => {
        hljs.highlightElement(block);
      });
    }

    // Renderizar ecuaciones matemáticas con KaTeX
    if (window.renderMathInElement) {
      renderMathInElement(outputDiv, {
        delimiters: [
          {left: '$$', right: '$$', display: true},
          {left: '$', right: '$', display: false},
          {left: '\\[', right: '\\]', display: true},
          {left: '\\(', right: '\\)', display: false}
        ],
        throwOnError: false
      });
    }

    // Renderizar diagramas Mermaid
    if (window.mermaid) {
      // Buscar bloques de código con clase 'language-mermaid'
      outputDiv.querySelectorAll('pre code.language-mermaid').forEach((block, index) => {
        const code = block.textContent;
        const id = `mermaid-${Date.now()}-${index}`;
        const container = document.createElement('div');
        container.className = 'mermaid';
        container.textContent = code;
        container.id = id;
        block.parentElement.replaceWith(container);
      });
      
      // Re-inicializar mermaid para los nuevos elementos
      mermaid.run({
        querySelector: '.mermaid'
      });
    }
  }

  document.addEventListener("DOMContentLoaded", () => {
    loadWasm();
  });
})();
