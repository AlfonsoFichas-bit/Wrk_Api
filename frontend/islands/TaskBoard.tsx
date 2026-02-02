import { useState, useEffect } from "preact/hooks";
import { tasks as tasksApi, type Task, type ProjectMember } from "../api_client.ts";

interface TaskBoardProps {
  projectId: string;
  sprintId?: string;
  members: ProjectMember[];
}

const COLUMNS = [
  { id: "TODO", title: "Por hacer", color: "bg-gray-100" },
  { id: "IN_PROGRESS", title: "En progreso", color: "bg-blue-50" },
  { id: "DONE", title: "Finalizado", color: "bg-green-50" },
];

export default function TaskBoard({ projectId, sprintId, members }: TaskBoardProps) {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    fetchTasks();
  }, [projectId, sprintId]);

  const fetchTasks = async () => {
    try {
      setLoading(true);
      const allTasks = await tasksApi.getAll(projectId);
      // Filter by sprint if provided
      const filtered = sprintId
        ? allTasks.filter(t => t.sprintId === sprintId)
        : allTasks;
      setTasks(filtered);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error al cargar tareas");
    } finally {
      setLoading(false);
    }
  };

  const updateTaskStatus = async (taskId: string, newStatus: string) => {
    try {
      await tasksApi.update(taskId, { status: newStatus });
      setTasks(tasks.map(t => t.id === taskId ? { ...t, status: newStatus } : t));
    } catch (err) {
      alert("Error al actualizar estado");
    }
  };

  if (loading) return <div class="text-center p-10">Cargando tablero...</div>;

  return (
    <div class="space-y-6">
      <div class="flex justify-between items-center px-2">
        <h2 class="text-xl font-bold text-gray-900">Tablero Kanban</h2>
        <button
          onClick={fetchTasks}
          class="text-sm text-primary hover:underline flex items-center gap-1"
        >
          <span class="material-symbols-outlined text-[18px]">refresh</span>
          Actualizar
        </button>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 h-full min-h-[500px]">
        {COLUMNS.map((column) => (
          <div key={column.id} class={`rounded-2xl border border-gray-200 flex flex-col ${column.color}`}>
            <div class="px-4 py-3 border-b border-gray-200 bg-white/50 rounded-t-2xl flex justify-between items-center">
              <h3 class="font-bold text-gray-700 uppercase text-xs tracking-wider">{column.title}</h3>
              <span class="bg-white px-2 py-0.5 rounded-full text-xs font-bold text-gray-500 border border-gray-100">
                {tasks.filter(t => t.status === column.id).length}
              </span>
            </div>

            <div class="p-3 flex-1 space-y-3 overflow-y-auto max-h-[600px]">
              {tasks.filter(t => t.status === column.id).map((task) => (
                <div key={task.id} class="bg-white p-4 rounded-xl border border-gray-200 shadow-sm hover:shadow-md transition-shadow group relative">
                  <div class="flex items-center gap-2 mb-2">
                    <span class={`text-[10px] font-bold px-1.5 py-0.5 rounded uppercase ${
                      task.priority === 'HIGH' ? 'bg-red-100 text-red-700' :
                      task.priority === 'MEDIUM' ? 'bg-orange-100 text-orange-700' :
                      'bg-blue-100 text-blue-700'
                    }`}>
                      {task.priority}
                    </span>
                  </div>
                  <h4 class="font-semibold text-gray-900 text-sm mb-1">{task.title}</h4>
                  <p class="text-xs text-gray-500 line-clamp-2 mb-3">{task.description}</p>

                  <div class="flex justify-between items-center mt-3 pt-3 border-t border-gray-50">
                    <div class="flex items-center gap-2">
                      <div class="size-6 rounded-full bg-gray-100 flex items-center justify-center text-[10px] font-bold text-gray-500 border border-white">
                        {task.assignee?.name?.charAt(0) || '?'}
                      </div>
                      <span class="text-[10px] text-gray-400 font-medium">
                        {task.assignee?.name || "Sin asignar"}
                      </span>
                    </div>
                  </div>

                  {/* Actions overlay for mobile/simplicity */}
                  <div class="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity flex gap-1">
                    {column.id !== "TODO" && (
                      <button
                        onClick={() => updateTaskStatus(task.id, "TODO")}
                        class="p-1 hover:bg-gray-100 rounded text-gray-400"
                        title="Mover a TODO"
                      >
                        <span class="material-symbols-outlined text-sm">first_page</span>
                      </button>
                    )}
                    {column.id === "TODO" && (
                      <button
                        onClick={() => updateTaskStatus(task.id, "IN_PROGRESS")}
                        class="p-1 hover:bg-gray-100 rounded text-blue-400"
                        title="Mover a EN PROGRESO"
                      >
                        <span class="material-symbols-outlined text-sm">arrow_forward</span>
                      </button>
                    )}
                    {column.id === "IN_PROGRESS" && (
                      <button
                        onClick={() => updateTaskStatus(task.id, "DONE")}
                        class="p-1 hover:bg-gray-100 rounded text-green-400"
                        title="Mover a FINALIZADO"
                      >
                        <span class="material-symbols-outlined text-sm">check_circle</span>
                      </button>
                    )}
                  </div>
                </div>
              ))}

              {tasks.filter(t => t.status === column.id).length === 0 && (
                <div class="text-center py-10 opacity-30 select-none">
                  <span class="material-symbols-outlined text-3xl">inbox</span>
                  <p class="text-[10px] font-medium mt-1">Vacío</p>
                </div>
              )}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
