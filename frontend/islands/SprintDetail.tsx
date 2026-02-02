import { useEffect, useState } from "preact/hooks";
import { sprints, projects, type Sprint, type Project } from "../../api_client.ts";
import TaskBoard from "./TaskBoard.tsx";
import CreateTaskForm from "./CreateTaskForm.tsx";

export default function SprintDetail({ id }: { id: string }) {
  const [sprint, setSprint] = useState<Sprint | null>(null);
  const [project, setProject] = useState<Project | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    fetchData();
  }, [id]);

  const fetchData = async () => {
    try {
      const s = await sprints.getById(id);
      setSprint(s);
      const p = await projects.getById(s.projectId);
      setProject(p);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error al cargar sprint");
    } finally {
      setLoading(false);
    }
  };

  if (loading) return <div class="text-center py-20">Cargando tablero del sprint...</div>;
  if (error) return <div class="p-8 bg-red-50 text-red-700 rounded-xl m-8">{error}</div>;
  if (!sprint || !project) return <div class="p-8 text-center">Sprint no encontrado</div>;

  return (
    <div class="space-y-8">
      <div class="mb-6 flex justify-between items-center">
        <a href={`/dashboard/projects/${sprint.projectId}`} class="text-primary hover:underline flex items-center gap-1 font-medium">
          <span class="material-symbols-outlined text-[18px]">arrow_back</span>
          Volver al Proyecto
        </a>
        <div class="flex items-center gap-2">
           <span class={`text-xs font-bold px-3 py-1 rounded-full ${
              sprint.status === 'ACTIVE' ? 'bg-green-100 text-green-700' : 'bg-orange-100 text-orange-700'
           }`}>
             {sprint.status}
           </span>
        </div>
      </div>

      <header class="bg-white p-8 rounded-2xl border border-gray-200 shadow-sm">
        <div class="flex flex-col md:flex-row justify-between items-start gap-4 mb-6">
          <div class="space-y-2">
            <h1 class="text-3xl font-bold text-gray-900">{sprint.name}</h1>
            <p class="text-gray-500">
              {new Date(sprint.startDate).toLocaleDateString()} - {new Date(sprint.endDate).toLocaleDateString()}
            </p>
          </div>
          <CreateTaskForm
            projectId={sprint.projectId}
            sprintId={sprint.id}
            userStories={sprint.userStories || []}
            members={project.members || []}
          />
        </div>
        <p class="text-gray-700 max-w-3xl leading-relaxed">
          {sprint.description || "Sin descripción para este sprint."}
        </p>
      </header>

      <TaskBoard
        projectId={sprint.projectId}
        sprintId={sprint.id}
        members={project.members || []}
      />
    </div>
  );
}
