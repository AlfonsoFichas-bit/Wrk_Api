import { useEffect, useState } from "preact/hooks";
import { projects, type Project } from "../../api_client.ts";
import UserStoryList from "./UserStoryList.tsx";
import CreateUserStoryForm from "./CreateUserStoryForm.tsx";
import SprintList from "./SprintList.tsx";
import CreateSprintForm from "./CreateSprintForm.tsx";
import TaskBoard from "./TaskBoard.tsx";
import CreateTaskForm from "./CreateTaskForm.tsx";
import RubricManager from "./RubricManager.tsx";

export default function ProjectDetail({ id }: { id: string }) {
  const [project, setProject] = useState<Project | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    projects.getById(id)
      .then(setProject)
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, [id]);

  if (loading) return <div class="text-center py-20">Cargando detalles del proyecto...</div>;
  if (error) return <div class="p-8 bg-red-50 text-red-700 rounded-xl m-8">{error}</div>;
  if (!project) return <div class="p-8 text-center">Proyecto no encontrado</div>;

  return (
    <div class="space-y-8">
      <header class="bg-white p-8 rounded-2xl border border-gray-200 shadow-sm">
        <div class="flex flex-col md:flex-row justify-between items-start gap-4 mb-6">
          <div class="space-y-2">
            <div class="flex items-center gap-3">
              <div class="p-3 bg-primary/10 rounded-xl text-primary">
                <span class="material-symbols-outlined text-3xl">dataset</span>
              </div>
              <div>
                <h1 class="text-3xl font-bold text-gray-900">{project.name}</h1>
                <p class="text-gray-500">{project.status}</p>
              </div>
            </div>
          </div>
          <div class="flex gap-2">
            <button class="px-4 py-2 border border-gray-300 rounded-lg font-medium hover:bg-gray-50 transition-colors">
              Editar Proyecto
            </button>
            <button class="px-4 py-2 bg-primary text-white rounded-lg font-bold hover:bg-primary-hover transition-colors shadow-lg shadow-primary/20">
              Nuevo Sprint
            </button>
          </div>
        </div>
        <p class="text-gray-700 max-w-3xl leading-relaxed">
          {project.description || "Sin descripción proporcionada para este proyecto."}
        </p>
      </header>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Main Content: Sprints/Tasks */}
        <div class="lg:col-span-2 space-y-10">
          <div class="space-y-4">
            <div class="flex justify-between items-center px-2">
              <h2 class="text-xl font-bold text-gray-900">Ejecución y Tareas</h2>
              <CreateTaskForm
                projectId={id}
                userStories={project.userStories || []}
                members={project.members || []}
              />
            </div>
            <TaskBoard
              projectId={id}
              members={project.members || []}
            />
          </div>

          <div class="space-y-4">
            <div class="flex justify-between items-center px-2">
              <h2 class="text-xl font-bold text-gray-900">Planificación de Sprints</h2>
              <CreateSprintForm projectId={id} />
            </div>
            <SprintList sprints={project.sprints || []} />
          </div>

          <div class="space-y-4">
            <div class="flex justify-between items-center px-2">
              <h2 class="text-xl font-bold text-gray-900">Backlog del Proyecto</h2>
              <CreateUserStoryForm projectId={id} />
            </div>
            <UserStoryList projectId={id} initialStories={project.userStories || []} sprints={project.sprints || []} />
          </div>
        </div>

        {/* Sidebar: Members/Owner */}
        <div class="space-y-6">
          <div class="bg-white rounded-2xl border border-gray-200 shadow-sm overflow-hidden">
            <div class="px-6 py-4 border-b border-gray-100 bg-gray-50 flex justify-between items-center">
              <h2 class="font-bold text-gray-900">Equipo</h2>
              <button class="text-primary hover:underline text-sm font-bold">+ Añadir</button>
            </div>
            <div class="p-4 space-y-4">
              <div class="flex items-center gap-3">
                <div class="size-10 rounded-full bg-blue-100 flex items-center justify-center text-blue-600 font-bold">
                  {project.owner?.name?.charAt(0)}
                </div>
                <div>
                  <p class="text-sm font-bold text-gray-900">{project.owner?.name}</p>
                  <p class="text-xs text-gray-500">Owner / Stakeholder</p>
                </div>
              </div>
              {project.members?.map((member) => (
                <div key={member.id} class="flex items-center gap-3">
                  <div class="size-10 rounded-full bg-gray-100 flex items-center justify-center text-gray-600 font-bold">
                    {member.user?.name?.charAt(0)}
                  </div>
                  <div>
                    <p class="text-sm font-bold text-gray-900">{member.user?.name}</p>
                    <p class="text-xs text-gray-500">{member.role}</p>
                  </div>
                </div>
              ))}
            </div>
          </div>

          <div class="bg-white rounded-2xl border border-gray-200 shadow-sm overflow-hidden">
            <div class="px-6 py-4 border-b border-gray-100 bg-gray-50">
              <h2 class="font-bold text-gray-900">Rúbricas</h2>
            </div>
            <div class="p-4">
               <RubricManager projectId={id} />
            </div>
          </div>

          <div class="bg-white rounded-2xl border border-gray-200 shadow-sm overflow-hidden">
            <div class="px-6 py-4 border-b border-gray-100 bg-gray-50">
              <h2 class="font-bold text-gray-900">Métricas Rápidas</h2>
            </div>
            <div class="p-6 space-y-4">
              <div class="flex justify-between items-center">
                <span class="text-sm text-gray-500">Velocidad Promedio</span>
                <span class="font-bold text-gray-900">0 pts</span>
              </div>
              <div class="flex justify-between items-center">
                <span class="text-sm text-gray-500">Historias Completadas</span>
                <span class="font-bold text-gray-900">0%</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
