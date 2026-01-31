import { useState } from "preact/hooks";
import { projects, auth } from "../api_client.ts";

export default function CreateProjectForm() {
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [showModal, setShowModal] = useState(false);

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    const user = auth.getUser();
    if (!user) {
      setError("Usuario no autenticado");
      setLoading(false);
      return;
    }

    try {
      await projects.create({
        name,
        description,
        ownerId: user.id,
        status: "ACTIVE"
      });
      window.location.reload();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error al crear proyecto");
    } finally {
      setLoading(false);
    }
  };

  const user = auth.getUser();
  const canCreate = user && (user.role === "ADMIN" || user.role === "SCRUM_MASTER" || user.role === "DOCENTE");

  if (!canCreate) return null;

  return (
    <>
      <button
        onClick={() => setShowModal(true)}
        class="flex items-center gap-2 px-4 py-2 bg-primary text-white rounded-lg font-bold hover:bg-primary-hover transition-colors shadow-lg shadow-primary/20"
      >
        <span class="material-symbols-outlined">add</span>
        Nuevo Proyecto
      </button>

      {showModal && (
        <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm">
          <div class="bg-white rounded-xl shadow-2xl w-full max-w-md overflow-hidden">
            <div class="px-6 py-4 border-b border-gray-100 flex justify-between items-center bg-gray-50">
              <h3 class="font-bold text-gray-900">Crear Nuevo Proyecto</h3>
              <button onClick={() => setShowModal(false)} class="text-gray-400 hover:text-gray-600">
                <span class="material-symbols-outlined">close</span>
              </button>
            </div>
            <form onSubmit={handleSubmit} class="p-6 space-y-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Nombre del Proyecto</label>
                <input
                  type="text"
                  required
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-primary focus:border-primary outline-none transition-all"
                  placeholder="Ej: WorkflowS"
                  value={name}
                  onInput={(e) => setName(e.currentTarget.value)}
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Descripción</label>
                <textarea
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-primary focus:border-primary outline-none transition-all min-h-[100px]"
                  placeholder="Breve descripción del proyecto..."
                  value={description}
                  onInput={(e) => setDescription(e.currentTarget.value)}
                ></textarea>
              </div>

              {error && <p class="text-sm text-red-600 bg-red-50 p-2 rounded">{error}</p>}

              <div class="flex gap-3 pt-2">
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
                  {loading ? "Creando..." : "Crear Proyecto"}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </>
  );
}
