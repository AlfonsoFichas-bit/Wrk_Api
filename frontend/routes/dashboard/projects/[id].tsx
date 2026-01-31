import { Head } from "fresh/runtime";
import ProjectDetail from "../../../islands/ProjectDetail.tsx";
import { define } from "../../../utils.ts";

export default define.page(function ProjectDetailPage(ctx) {
  const { id } = ctx.params;

  return (
    <div class="min-h-screen bg-gray-50 p-4 md:p-8">
      <Head>
        <title>Detalle del Proyecto | WorkflowS</title>
        <link
          href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:wght,FILL@100..700,0..1&display=swap"
          rel="stylesheet"
        />
      </Head>

      <div class="max-w-[1280px] mx-auto">
        <div class="mb-6">
          <a href="/dashboard/projects" class="text-primary hover:underline flex items-center gap-1 font-medium">
            <span class="material-symbols-outlined text-[18px]">arrow_back</span>
            Volver a Proyectos
          </a>
        </div>

        <ProjectDetail id={id} />
      </div>
    </div>
  );
});
