import { useState } from "preact/hooks";
import { userStories } from "../api_client.ts";

interface CreateUserStoryFormProps {
  projectId: string;
}

export default function CreateUserStoryForm({ projectId }: CreateUserStoryFormProps) {
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [acceptance, setAcceptance] = useState("");
  const [priority, setPriority] = useState("MEDIUM");
  const [storyPoints, setStoryPoints] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [showModal, setShowModal] = useState(false);

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      await userStories.create({
        title,
        description,
        acceptance,
        priority,
        storyPoints,
        projectId,
        status: "BACKLOG"
      });
      window.location.reload(); // Refresh to show new story
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error al crear historia de usuario");
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <button
        onClick={() => setShowModal(true)}
        class="flex items-center gap-2 px-3 py-1.5 bg-white border border-gray-200 text-gray-700 rounded-lg text-sm font-bold hover:border-primary hover:text-primary transition-all shadow-sm"
      >
        <span class="material-symbols-outlined text-[18px]">add</span>
        Añadir Historia
      </button>

      {showModal && (
        <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm">
          <div class="bg-white rounded-xl shadow-2xl w-full max-w-lg overflow-hidden">
            <div class="px-6 py-4 border-b border-gray-100 flex justify-between items-center bg-gray-50">
              <h3 class="font-bold text-gray-900">Nueva Historia de Usuario</h3>
              <button onClick={() => setShowModal(false)} class="text-gray-400 hover:text-gray-600">
                <span class="material-symbols-outlined">close</span>
              </button>
            </div>
            <form onSubmit={handleSubmit} class="p-6 space-y-4 max-h-[80vh] overflow-y-auto">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Título</label>
                <input
                  type="text"
                  required
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-primary focus:border-primary outline-none transition-all"
                  placeholder="Como [rol] quiero [acción] para [beneficio]"
                  value={title}
                  onInput={(e) => setTitle(e.currentTarget.value)}
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Descripción</label>
                <textarea
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-primary focus:border-primary outline-none transition-all min-h-[80px]"
                  placeholder="Detalles adicionales..."
                  value={description}
                  onInput={(e) => setDescription(e.currentTarget.value)}
                ></textarea>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Criterios de Aceptación</label>
                <textarea
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-primary focus:border-primary outline-none transition-all min-h-[80px]"
                  placeholder="Dado que..., Cuando..., Entonces..."
                  value={acceptance}
                  onInput={(e) => setAcceptance(e.currentTarget.value)}
                ></textarea>
              </div>

              <div class="grid grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Prioridad</label>
                  <select
                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-primary focus:border-primary outline-none transition-all bg-white"
                    value={priority}
                    onChange={(e) => setPriority(e.currentTarget.value)}
                  >
                    <option value="LOW">Baja</option>
                    <option value="MEDIUM">Media</option>
                    <option value="HIGH">Alta</option>
                    <option value="URGENT">Urgente</option>
                  </select>
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Story Points</label>
                  <input
                    type="number"
                    min="0"
                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-primary focus:border-primary outline-none transition-all"
                    value={storyPoints}
                    onInput={(e) => setStoryPoints(parseInt(e.currentTarget.value) || 0)}
                  />
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
                  class="flex-1 px-4 py-2 bg-primary text-white rounded-lg font-bold hover:bg-primary-hover transition-colors disabled:opacity-50"
                >
                  {loading ? "Guardando..." : "Crear Historia"}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </>
  );
}
