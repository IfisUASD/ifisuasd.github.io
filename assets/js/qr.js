// assets/js/qr.js
(() => {
  const go = new Go();
  let wasmReady = false;

  // Helper: Debounce
  function debounce(fn, delay) {
    let timer;
    return (...args) => {
      clearTimeout(timer);
      timer = setTimeout(() => fn.apply(this, args), delay);
    };
  }

  // Helper: Read file as base64
  function readFileAsBase64(file) {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = () => resolve(reader.result.split(",")[1]);
      reader.onerror = reject;
      reader.readAsDataURL(file);
    });
  }

  async function loadWasm() {
    try {
      // Nota: La ruta del WASM será relativa a la raíz del sitio generado
      const response = await fetch("/assets/wasm/qr.wasm");
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

      const loader = document.getElementById("qr-loader");
      const app = document.getElementById("qr-app");
      if (loader) loader.classList.add("hidden");
      if (app) app.classList.remove("hidden");

      updateQR();
    } catch (err) {
      console.error("Error loading WASM:", err);
      const loader = document.getElementById("qr-loader");
      if (loader)
        loader.innerHTML = `<span class="text-error">Error: ${err.message}</span>`;
    }
  }

  async function updateQR() {
    if (!wasmReady) return;

    const text =
      document.getElementById("textInput").value || "https://ifis.uasd.edu.do";
    const size = parseInt(document.getElementById("sizeSlider").value, 10);
    const fg = document.getElementById("fgColor").value;
    const bg = document.getElementById("bgColor").value;
    const logoInput = document.getElementById("logoInput");
    const level = document.getElementById("levelSelect").value;
    let logoBase64 = "";

    document.getElementById("sizeValue").textContent = size;

    if (logoInput.files && logoInput.files[0]) {
      try {
        logoBase64 = await readFileAsBase64(logoInput.files[0]);
      } catch (e) {
        console.warn("Failed to read logo:", e);
      }
    }

    // generateQR es expuesto por Go (main.go)
    const dataURL = generateQR(text, size, fg, bg, logoBase64, level);
    const img = document.getElementById("qrImage");
    img.src = dataURL;
  }

  const debouncedUpdate = debounce(updateQR, 300);

  // Init
  document.addEventListener("DOMContentLoaded", () => {
    loadWasm();

    const elements = [
      "textInput",
      "sizeSlider",
      "fgColor",
      "bgColor",
      "levelSelect",
    ];
    elements.forEach((id) => {
      const el = document.getElementById(id);
      if (el) el.addEventListener("input", debouncedUpdate);
    });

    const logoInput = document.getElementById("logoInput");
    if (logoInput) logoInput.addEventListener("change", debouncedUpdate);

    // Download
    document.getElementById("downloadBtn").addEventListener("click", () => {
      const img = document.getElementById("qrImage");
      if (!img.src) return;
      const a = document.createElement("a");
      a.href = img.src;
      a.download = `qrcode-${Date.now()}.png`;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
    });

    // Copy
    document.getElementById("copyBtn").addEventListener("click", async () => {
      const img = document.getElementById("qrImage");
      if (!img.src) return;
      try {
        const response = await fetch(img.src);
        const blob = await response.blob();
        await navigator.clipboard.write([
          new ClipboardItem({ [blob.type]: blob }),
        ]);
        const btn = document.getElementById("copyBtn");
        btn.classList.add("text-success");
        setTimeout(() => btn.classList.remove("text-success"), 2000);
      } catch (e) {
        alert("Error al copiar");
      }
    });
  });
})();
