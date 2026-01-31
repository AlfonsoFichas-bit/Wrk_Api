import { Head } from "fresh/runtime";
import ProjectList from "../../islands/ProjectList.tsx";
import CreateProjectForm from "../../islands/CreateProjectForm.tsx";

export default function ProjectsPage() {
  return (
    <div class="min-h-screen bg-gray-50 p-4 md:p-8">
      <Head>
        <title>Mis Proyectos | WorkflowS</title>
        <link
          href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:wght,FILL@100..700,0..1&display=swap"
          rel="stylesheet"
        />
      </Head>

      <div class="max-w-[1280px] mx-auto space-y-8">
        <header class="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
          <div>
            <h1 class="text-3xl font-bold text-gray-900">Gestión de Proyectos</h1>
            <p class="text-gray-500">Visualiza y administra tus espacios de trabajo colaborativos.</p>
          </div>
          <CreateProjectForm />
        </header>

        <div class="h-px w-full bg-gray-200"></div>

        <ProjectList />
      </div>
    </div>
  );
}
