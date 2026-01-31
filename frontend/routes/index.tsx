import { Head } from "fresh/runtime";
import { define } from "../utils.ts";

export default define.page(function Home() {
  return (
    <div class="bg-background-light dark:bg-background-dark font-display text-gray-900 dark:text-white antialiased overflow-x-hidden">
      <Head>
        <title>WorkflowS - Gestión de Proyectos Académicos</title>
        <link
          href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;700;900&display=swap"
          rel="stylesheet"
        />
        <link
          href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:wght,FILL@100..700,0..1&display=swap"
          rel="stylesheet"
        />
      </Head>

      <header class="sticky top-0 z-50 w-full border-b border-solid border-border-dark bg-background-dark/80 backdrop-blur-md">
        <div class="px-4 md:px-10 py-3 flex items-center justify-between max-w-[1280px] mx-auto">
          <div class="flex items-center gap-4 text-white">
            <div class="flex items-center justify-center size-8 rounded bg-primary/20 text-primary">
              <span class="material-symbols-outlined text-[20px]">dataset</span>
            </div>
            <h2 class="text-white text-lg font-bold leading-tight tracking-[-0.015em]">
              WorkflowS
            </h2>
          </div>
          <nav class="hidden md:flex items-center gap-8">
            <a
              class="text-white text-sm font-medium hover:text-primary transition-colors"
              href="#"
            >
              Inicio
            </a>
            <a
              class="text-white text-sm font-medium hover:text-primary transition-colors"
              href="#features"
            >
              Características
            </a>
            <a
              class="text-white text-sm font-medium hover:text-primary transition-colors"
              href="#roles"
            >
              Roles
            </a>
            <a
              class="text-white text-sm font-medium hover:text-primary transition-colors"
              href="#tech"
            >
              Tecnología
            </a>
          </nav>
          <div class="flex items-center gap-4">
            <a
              href="/login"
              class="flex items-center justify-center overflow-hidden rounded-lg h-9 px-4 bg-primary hover:bg-primary-hover text-white text-sm font-bold transition-all shadow-lg shadow-primary/20"
            >
              <span class="truncate">Iniciar Sesión</span>
            </a>
          </div>
        </div>
      </header>

      <main class="flex flex-col items-center w-full">
        {/* Hero Section */}
        <section class="relative w-full py-16 md:py-24 px-4 md:px-10 max-w-[1280px]">
          <div class="flex flex-col gap-10 lg:flex-row items-center">
            {/* Text Content */}
            <div class="flex flex-col gap-6 lg:w-1/2">
              <div class="flex flex-col gap-4">
                <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-primary/10 border border-primary/20 w-fit">
                  <span class="material-symbols-outlined text-primary text-[16px]">
                    school
                  </span>
                  <span class="text-primary text-xs font-bold uppercase tracking-wider">
                    Diseñado para U. La Salle
                  </span>
                </div>
                <h1 class="text-white text-4xl md:text-5xl lg:text-6xl font-black leading-tight tracking-tight">
                  Gestión de Proyectos Académicos con{" "}
                  <span class="text-primary">Agilidad Real</span>{" "}
                  y Rigor Evaluativo
                </h1>
                <p class="text-text-secondary text-lg leading-relaxed max-w-[600px]">
                  La primera plataforma que une la metodología Scrum con
                  rúbricas de evaluación dinámicas. Organiza, colabora y evalúa
                  en un solo lugar.
                </p>
              </div>
              <div class="flex flex-wrap gap-4 pt-2">
                <a
                  href="/register"
                  class="flex items-center justify-center rounded-lg h-12 px-6 bg-primary hover:bg-primary-hover text-white text-base font-bold shadow-lg shadow-primary/25 transition-all"
                >
                  Empieza a gestionar
                </a>
                <button class="flex items-center justify-center rounded-lg h-12 px-6 bg-card-dark border border-border-dark hover:border-primary/50 text-white text-base font-bold transition-all">
                  <span class="material-symbols-outlined mr-2 text-[20px]">
                    play_circle
                  </span>
                  Ver Demo
                </button>
              </div>
            </div>
            {/* Hero Image */}
            <div class="lg:w-1/2 w-full relative">
              <div class="absolute -inset-1 bg-gradient-to-r from-primary to-purple-600 rounded-2xl blur opacity-30 animate-pulse"></div>
              <div class="relative w-full aspect-video rounded-xl overflow-hidden border border-border-dark shadow-2xl bg-card-dark">
                {/* Simulated UI */}
                <div class="w-full h-full flex flex-col bg-[#0f1520]">
                  <div class="h-12 border-b border-border-dark flex items-center px-4 gap-3">
                    <div class="w-3 h-3 rounded-full bg-red-500"></div>
                    <div class="w-3 h-3 rounded-full bg-yellow-500"></div>
                    <div class="w-3 h-3 rounded-full bg-green-500"></div>
                    <div class="ml-4 h-2 w-32 bg-border-dark rounded-full"></div>
                  </div>
                  <div class="flex-1 p-4 flex gap-4 overflow-hidden">
                    <div class="flex-1 bg-[#1a2333] rounded-lg p-3 flex flex-col gap-3">
                      <div class="flex justify-between items-center mb-1">
                        <div class="h-3 w-20 bg-text-secondary/30 rounded"></div>
                        <div class="h-4 w-6 bg-primary/20 rounded text-xs text-primary text-center">
                          3
                        </div>
                      </div>
                      <div class="bg-border-dark p-3 rounded border-l-2 border-l-blue-400">
                        <div class="h-2 w-3/4 bg-text-secondary/50 rounded mb-2"></div>
                        <div class="h-2 w-1/2 bg-text-secondary/30 rounded"></div>
                      </div>
                      <div class="bg-border-dark p-3 rounded border-l-2 border-l-purple-400">
                        <div class="h-2 w-full bg-text-secondary/50 rounded mb-2"></div>
                        <div class="flex gap-2 mt-2">
                          <div class="size-5 rounded-full bg-primary/40"></div>
                        </div>
                      </div>
                    </div>
                    <div class="flex-1 bg-[#1a2333] rounded-lg p-3 flex flex-col gap-3">
                      <div class="flex justify-between items-center mb-1">
                        <div class="h-3 w-24 bg-text-secondary/30 rounded"></div>
                        <div class="h-4 w-6 bg-orange-500/20 rounded text-xs text-orange-500 text-center">
                          2
                        </div>
                      </div>
                      <div class="bg-border-dark p-3 rounded border-l-2 border-l-orange-400">
                        <div class="h-2 w-2/3 bg-text-secondary/50 rounded mb-2"></div>
                        <div class="h-2 w-full bg-text-secondary/30 rounded mb-2"></div>
                        <div class="w-full bg-gray-700 h-1.5 rounded-full mt-2">
                          <div
                            class="bg-orange-500 h-1.5 rounded-full"
                            style="width: 60%"
                          >
                          </div>
                        </div>
                      </div>
                    </div>
                    <div class="flex-1 bg-[#1a2333] rounded-lg p-3 flex flex-col gap-3 opacity-80">
                      <div class="flex justify-between items-center mb-1">
                        <div class="h-3 w-16 bg-text-secondary/30 rounded"></div>
                        <div class="h-4 w-6 bg-green-500/20 rounded text-xs text-green-500 text-center">
                          5
                        </div>
                      </div>
                      <div class="bg-border-dark p-3 rounded border-l-2 border-l-green-400 opacity-50">
                        <div class="h-2 w-1/2 bg-text-secondary/50 rounded mb-2"></div>
                      </div>
                      <div class="bg-border-dark p-3 rounded border-l-2 border-l-green-400 opacity-50">
                        <div class="h-2 w-3/4 bg-text-secondary/50 rounded mb-2"></div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* Problem vs Solution */}
        <section class="w-full py-16 px-4 md:px-10 max-w-[1280px]">
          <div class="flex flex-col gap-10">
            <div class="max-w-[720px]">
              <h2 class="text-white text-3xl font-bold leading-tight mb-4">
                El Problema vs. La Solución
              </h2>
              <p class="text-text-secondary text-lg">
                Olvídate del caos de gestionar proyectos por WhatsApp y Excel.
                WorkflowS centraliza la comunicación y visualiza el avance real
                para que te enfoques en aprender.
              </p>
            </div>
            <div class="grid md:grid-cols-2 gap-6">
              {/* Problem Card */}
              <div class="group relative flex flex-col gap-4 rounded-xl border border-red-900/30 bg-[#1a1111] p-6 md:p-8 hover:border-red-900/60 transition-colors">
                <div class="absolute top-6 right-6 p-2 rounded-full bg-red-500/10 text-red-500">
                  <span class="material-symbols-outlined">warning</span>
                </div>
                <h3 class="text-white text-xl font-bold">Caos Actual</h3>
                <p class="text-text-secondary">
                  "¿En qué archivo está la última versión?", "¿Qué tarea me
                  toca?", "¿Por qué tengo baja nota si trabajé más?". <br />
                  <br />
                  La falta de visibilidad, la evaluación subjetiva y los
                  archivos perdidos en chats grupales destruyen la productividad
                  del equipo.
                </p>
                <div class="mt-auto pt-4 flex gap-2">
                  <div class="h-1 flex-1 bg-red-900/30 rounded-full"></div>
                </div>
              </div>
              {/* Solution Card */}
              <div class="group relative flex flex-col gap-4 rounded-xl border border-primary/30 bg-primary/5 p-6 md:p-8 hover:border-primary/60 transition-colors shadow-inner shadow-primary/5">
                <div class="absolute top-6 right-6 p-2 rounded-full bg-primary/20 text-primary">
                  <span class="material-symbols-outlined">check_circle</span>
                </div>
                <h3 class="text-white text-xl font-bold">Solución WorkflowS</h3>
                <p class="text-text-secondary">
                  <strong class="text-white">Centralización:</strong>{" "}
                  Chat, documentos y tareas en un solo lugar. <br />
                  <strong class="text-white">Transparencia:</strong>{" "}
                  Visualización en tiempo real del avance del Sprint. <br />
                  <br />
                  Todo sincronizado. Nadie se queda atrás.
                </p>
                <div class="mt-auto pt-4 flex gap-2">
                  <div class="h-1 flex-1 bg-primary/50 rounded-full"></div>
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* Key Features Grid */}
        <section class="w-full py-16 px-4 md:px-10 max-w-[1280px]" id="features">
          <h2 class="text-white text-3xl font-bold leading-tight mb-12 text-center">
            Características Clave
          </h2>
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            <div class="flex flex-col gap-4 p-6 rounded-xl bg-card-dark border border-border-dark hover:border-primary/40 transition-all hover:-translate-y-1">
              <div class="size-12 rounded-lg bg-blue-500/10 text-blue-500 flex items-center justify-center">
                <span class="material-symbols-outlined text-[28px]">
                  view_kanban
                </span>
              </div>
              <div>
                <h3 class="text-white text-lg font-bold mb-2">
                  Tablero Kanban Interactivo
                </h3>
                <p class="text-text-secondary text-sm leading-relaxed">
                  Gestión visual con Drag & Drop. Mueve tus tareas de
                  'Pendiente' a 'Completado' y sincroniza el estado con tu
                  equipo al instante.
                </p>
              </div>
            </div>
            <div class="flex flex-col gap-4 p-6 rounded-xl bg-card-dark border border-border-dark hover:border-primary/40 transition-all hover:-translate-y-1">
              <div class="size-12 rounded-lg bg-purple-500/10 text-purple-500 flex items-center justify-center">
                <span class="material-symbols-outlined text-[28px]">
                  sprint
                </span>
              </div>
              <div>
                <h3 class="text-white text-lg font-bold mb-2">
                  Sprints y Backlog
                </h3>
                <p class="text-text-secondary text-sm leading-relaxed">
                  Planificación profesional. Prioriza historias de usuario y
                  define ciclos de trabajo claros (Sprints) como en la industria
                  del software.
                </p>
              </div>
            </div>
            <div class="flex flex-col gap-4 p-6 rounded-xl bg-card-dark border border-border-dark hover:border-primary/40 transition-all hover:-translate-y-1">
              <div class="size-12 rounded-lg bg-emerald-500/10 text-emerald-500 flex items-center justify-center">
                <span class="material-symbols-outlined text-[28px]">
                  fact_check
                </span>
              </div>
              <div>
                <h3 class="text-white text-lg font-bold mb-2">
                  Rúbricas Dinámicas
                </h3>
                <p class="text-text-secondary text-sm leading-relaxed">
                  Adiós a la subjetividad. Los docentes evalúan entregables con
                  criterios predefinidos y claros directamente en la plataforma.
                </p>
              </div>
            </div>
            <div class="flex flex-col gap-4 p-6 rounded-xl bg-card-dark border border-border-dark hover:border-primary/40 transition-all hover:-translate-y-1">
              <div class="size-12 rounded-lg bg-orange-500/10 text-orange-500 flex items-center justify-center">
                <span class="material-symbols-outlined text-[28px]">
                  monitoring
                </span>
              </div>
              <div>
                <h3 class="text-white text-lg font-bold mb-2">
                  Métricas de Rendimiento
                </h3>
                <p class="text-text-secondary text-sm leading-relaxed">
                  Gráficos Burndown automáticos para visualizar el progreso
                  ideal vs. real y la velocidad del equipo sin esfuerzo manual.
                </p>
              </div>
            </div>
            <div class="flex flex-col gap-4 p-6 rounded-xl bg-card-dark border border-border-dark hover:border-primary/40 transition-all hover:-translate-y-1 md:col-span-2 lg:col-span-2">
              <div class="size-12 rounded-lg bg-pink-500/10 text-pink-500 flex items-center justify-center">
                <span class="material-symbols-outlined text-[28px]">
                  folder_shared
                </span>
              </div>
              <div>
                <h3 class="text-white text-lg font-bold mb-2">
                  Gestión Documental con Versionado
                </h3>
                <p class="text-text-secondary text-sm leading-relaxed">
                  Mantén el historial de tus archivos. Sube entregables y accede
                  siempre a la versión más reciente sin perder los cambios
                  anteriores. Integrado directamente en las tareas.
                </p>
              </div>
            </div>
          </div>
        </section>

        {/* Audience Section */}
        <section
          class="w-full py-16 px-4 md:px-10 bg-[#161e2e] border-y border-border-dark"
          id="roles"
        >
          <div class="max-w-[1280px] mx-auto">
            <h2 class="text-white text-3xl font-bold leading-tight mb-12 text-center">
              ¿Para quién es WorkflowS?
            </h2>
            <div class="flex flex-col md:flex-row gap-8">
              <div class="flex-1 flex flex-col items-center text-center p-8 rounded-2xl bg-gradient-to-b from-card-dark to-transparent border border-border-dark">
                <div class="p-4 bg-primary/10 rounded-full mb-6">
                  <span class="material-symbols-outlined text-primary text-[40px]">
                    person_apron
                  </span>
                </div>
                <h3 class="text-white text-2xl font-bold mb-3">Para Docentes</h3>
                <p class="text-text-secondary mb-6">Stakeholders & Evaluadores</p>
                <p class="text-gray-300 leading-relaxed max-w-sm">
                  Automatiza el seguimiento y reduce la carga administrativa.
                  Genera reportes de avance y exporta notas en CSV fácilmente.
                  Ten visibilidad total del aporte de cada alumno.
                </p>
              </div>
              <div class="flex-1 flex flex-col items-center text-center p-8 rounded-2xl bg-gradient-to-b from-card-dark to-transparent border border-border-dark">
                <div class="p-4 bg-green-500/10 rounded-full mb-6">
                  <span class="material-symbols-outlined text-green-500 text-[40px]">
                    groups
                  </span>
                </div>
                <h3 class="text-white text-2xl font-bold mb-3">Para Estudiantes</h3>
                <p class="text-text-secondary mb-6">Scrum Team</p>
                <p class="text-gray-300 leading-relaxed max-w-sm">
                  Aprendan roles reales: Product Owner, Scrum Master y Developer.
                  Eviten conflictos de equipo con una asignación clara de tareas
                  y notificaciones automáticas.
                </p>
              </div>
            </div>
          </div>
        </section>

        {/* Tech Stack */}
        <section class="w-full py-20 px-4 md:px-10 max-w-[1280px]" id="tech">
          <div class="flex flex-col items-center text-center gap-10">
            <div class="max-w-2xl">
              <h2 class="text-white text-3xl font-bold mb-4">
                Stack Tecnológico de Vanguardia
              </h2>
              <p class="text-text-secondary">
                Construido con las últimas tecnologías para garantizar
                velocidad, escalabilidad y una experiencia de usuario fluida.
              </p>
            </div>
            <div class="flex flex-wrap justify-center gap-4 md:gap-8">
              <div class="flex items-center gap-3 px-5 py-3 rounded-full bg-card-dark border border-border-dark hover:border-primary/50 transition-colors">
                <span class="material-symbols-outlined text-blue-400">code</span>
                <div class="flex flex-col text-left">
                  <span class="text-white font-bold text-sm">
                    React 19 + TypeScript
                  </span>
                  <span class="text-xs text-text-secondary">Frontend</span>
                </div>
              </div>
              <div class="flex items-center gap-3 px-5 py-3 rounded-full bg-card-dark border border-border-dark hover:border-primary/50 transition-colors">
                <span class="material-symbols-outlined text-green-400">dns</span>
                <div class="flex flex-col text-left">
                  <span class="text-white font-bold text-sm">
                    Go + Gin
                  </span>
                  <span class="text-xs text-text-secondary">Backend</span>
                </div>
              </div>
              <div class="flex items-center gap-3 px-5 py-3 rounded-full bg-card-dark border border-border-dark hover:border-primary/50 transition-colors">
                <span class="material-symbols-outlined text-purple-400">
                  storage
                </span>
                <div class="flex flex-col text-left">
                  <span class="text-white font-bold text-sm">
                    GORM + SQLite
                  </span>
                  <span class="text-xs text-text-secondary">Database</span>
                </div>
              </div>
              <div class="flex items-center gap-3 px-5 py-3 rounded-full bg-card-dark border border-border-dark hover:border-primary/50 transition-colors">
                <span class="material-symbols-outlined text-orange-400">
                  bug_report
                </span>
                <div class="flex flex-col text-left">
                  <span class="text-white font-bold text-sm">Playwright</span>
                  <span class="text-xs text-text-secondary">E2E Testing</span>
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* Footer */}
        <footer class="w-full border-t border-border-dark bg-[#0b1019] pt-12 pb-8 px-4 md:px-10">
          <div class="max-w-[1280px] mx-auto flex flex-col gap-10">
            <div class="flex flex-col md:flex-row justify-between gap-8">
              <div class="flex flex-col gap-4 max-w-sm">
                <div class="flex items-center gap-2 text-white">
                  <span class="material-symbols-outlined text-primary">
                    dataset
                  </span>
                  <span class="text-lg font-bold">WorkflowS</span>
                </div>
                <p class="text-text-secondary text-sm">
                  La plataforma definitiva para la gestión de proyectos
                  académicos en la Universidad La Salle.
                </p>
              </div>
              <div class="flex flex-col md:flex-row gap-10 md:gap-20">
                <div class="flex flex-col gap-4">
                  <h4 class="text-white font-bold text-sm uppercase tracking-wider">
                    Plataforma
                  </h4>
                  <a
                    class="text-text-secondary hover:text-primary text-sm transition-colors"
                    href="#"
                  >
                    Documentación
                  </a>
                  <a
                    class="text-text-secondary hover:text-primary text-sm transition-colors"
                    href="#"
                  >
                    Guía de Usuario
                  </a>
                  <a
                    class="text-text-secondary hover:text-primary text-sm transition-colors"
                    href="#"
                  >
                    Release Notes
                  </a>
                </div>
                <div class="flex flex-col gap-4">
                  <h4 class="text-white font-bold text-sm uppercase tracking-wider">
                    Contacto
                  </h4>
                  <a
                    class="text-text-secondary hover:text-primary text-sm transition-colors"
                    href="#"
                  >
                    Soporte Técnico
                  </a>
                  <a
                    class="text-text-secondary hover:text-primary text-sm transition-colors"
                    href="#"
                  >
                    Universidad La Salle
                  </a>
                </div>
              </div>
            </div>
            <div class="h-px w-full bg-border-dark"></div>
            <div class="flex flex-col md:flex-row justify-between items-center gap-4 text-xs text-text-secondary">
              <p>© 2025 WorkflowS. Todos los derechos reservados.</p>
              <p class="flex items-center gap-1">
                Desarrollado por{" "}
                <span class="text-white font-medium">Ronald Choque Sillo</span>
                {" "}- Ingeniería de Sistemas 2025
              </p>
            </div>
          </div>
        </footer>
      </main>
    </div>
  );
});
