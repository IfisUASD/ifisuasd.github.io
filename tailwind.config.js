/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./templates/**/*.templ", // Escuchar cambios en tus componentes Templ
    "./content/**/*.md", // Escuchar clases si las usaras dentro del markdown (opcional)
    "./cmd/**/*.go", // Escuchar si inyectas clases desde Go
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ["Inter", "sans-serif"], // Recomendado: Usar una fuente limpia
      },
    },
  },
  // Activamos DaisyUI para los componentes pre-hechos
  plugins: [require("daisyui")],

  // Configuración de DaisyUI
  daisyui: {
    themes: ["light", "dark"], // Habilitar temas claro y oscuro
    darkTheme: "dark",
  },
};
