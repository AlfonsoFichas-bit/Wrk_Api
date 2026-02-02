import { useState } from "preact/hooks";
import { tasks, type ProjectMember, type UserStory } from "../api_client.ts";

interface CreateTaskFormProps {
  projectId: string;
  sprintId?: string;
  userStories: UserStory[];
  members: ProjectMember[];
}

export default function CreateTaskForm({ projectId, sprintId, userStories, members }: CreateTaskFormProps) {
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [assigneeId, setAssigneeId] = useState("");
  const [userStoryId, setUserStoryId] = useState("");
  const [priority, setPriority] = useState("MEDIUM");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [showModal, setShowModal] = useState(false);

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      await tasks.create({
        title,
        description,
        projectId,
        sprintId,
        userStoryId: userStoryId || undefined,
        assigneeId: assigneeId || undefined,
        priority,
        status: "TODO"
      });
      window.location.reload();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error al crear tarea");
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <button
        onClick={() => setShowModal(true)}
        class="flex items-center gap-2 px-3 py-1.5 bg-blue-600 text-white rounded-lg text-sm font-bold hover:bg-blue-700 transition-all shadow-sm"
      >
        <span class="material-symbols-outlined text-[18px]">playlist_add</span>
        Añadir Tarea
      </button>

      {showModal && (
        <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm">
          <div class="bg-white rounded-xl shadow-2xl w-full max-w-lg overflow-hidden">
            <div class="px-6 py-4 border-b border-gray-100 flex justify-between items-center bg-gray-50">
              <h3 class="font-bold text-gray-900">Nueva Tarea</h3>
              <button onClick={() => setShowModal(false)} class="text-gray-400 hover:text-gray-600">
                <span class="material-symbols-outlined">close</span>
              </button>
            </div>
            <form onSubmit={handleSubmit} class="p-6 space-y-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Título de la tarea</label>
                <input
                  type="text"
                  required
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-primary focus:border-primary outline-none transition-all"
                  placeholder="Ej: Implementar Login"
                  value={title}
                  onInput={(e) => setTitle(e.currentTarget.value)}
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Descripción</label>
                <textarea
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-primary focus:border-primary outline-none transition-all min-h-[80px]"
                  placeholder="Detalles de la tarea..."
                  value={description}
                  onInput={(e) => setDescription(e.currentTarget.value)}
                ></textarea>
              </div>

              <div class="grid grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Historia Relacionada</label>
                  <select
                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-primary focus:border-primary outline-none transition-all bg-white text-sm"
                    value={userStoryId}
                    onChange={(e) => setUserStoryId(e.currentTarget.value)}
                  >
                    <option value="">Ninguna</option>
                    {userStories.map(s => (
                      <option key={s.id} value={s.id}>{s.title}</option>
                    ))}
                  </select>
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Responsable</label>
                  <select
                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-primary focus:border-primary outline-none transition-all bg-white text-sm"
                    value={assigneeId}
                    onChange={(e) => setAssigneeId(e.currentTarget.value)}
                  >
                    <option value="">Sin asignar</option>
                    {members.map(m => (
                      <option key={m.userId} value={m.userId}>{m.user?.name || "Cargando..."}</option>
                    ))}
                  </select>
                </div>
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Prioridad</label>
                <div class="flex gap-4">
                  {["LOW", "MEDIUM", "HIGH"].map((p) => (
                    <label key={p} class="flex items-center gap-2 cursor-pointer">
                      <input
                        type="radio"
                        name="priority"
                        value={p}
                        checked={priority === p}
                        onChange={() => setPriority(p)}
                        class="text-primary focus:ring-primary"
                      />
                      <span class="text-xs font-medium text-gray-600">{p}</span>
                    </label>
                  ))}
                </div>
              </div>

              {error && <p class="text-sm text-red-600 bg-red-50 p-2 rounded">{error}</p>}

              <div class="flex gap-3 pt-4">
                <button
                  type="button"
                  onClick={() => setShowModal(false)}
                  class="flex-1 px-4 py-2 border border-gray-300 text-gray-700 rounded-lg font-medium hover:bg-gray-50 transition-colors"
                >
                  Cancelar
                </button>
                <button
                  type="submit"
                  disabled={loading}
                  class="flex-1 px-4 py-2 bg-blue-600 text-white rounded-lg font-bold hover:bg-blue-700 transition-colors disabled:opacity-50"
                >
                  {loading ? "Guardando..." : "Crear Tarea"}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </>
  );
}
