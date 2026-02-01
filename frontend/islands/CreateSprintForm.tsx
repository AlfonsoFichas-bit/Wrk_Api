import { useState } from "preact/hooks";
import { sprints } from "../api_client.ts";

interface CreateSprintFormProps {
  projectId: string;
}

export default function CreateSprintForm({ projectId }: CreateSprintFormProps) {
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [startDate, setStartDate] = useState("");
  const [endDate, setEndDate] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [showModal, setShowModal] = useState(false);

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      await sprints.create({
        name,
        description,
        projectId,
        startDate: startDate ? new Date(startDate).toISOString() : undefined,
        endDate: endDate ? new Date(endDate).toISOString() : undefined,
        status: "PLANNING"
      });
      window.location.reload();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error al crear sprint");
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <button
        onClick={() => setShowModal(true)}
        class="flex items-center gap-2 px-3 py-1.5 bg-primary text-white rounded-lg text-sm font-bold hover:bg-primary-hover transition-all shadow-sm"
      >
        <span class="material-symbols-outlined text-[18px]">add</span>
        Nuevo Sprint
      </button>

      {showModal && (
        <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm">
          <div class="bg-white rounded-xl shadow-2xl w-full max-w-lg overflow-hidden">
            <div class="px-6 py-4 border-b border-gray-100 flex justify-between items-center bg-gray-50">
              <h3 class="font-bold text-gray-900">Configurar Nuevo Sprint</h3>
              <button onClick={() => setShowModal(false)} class="text-gray-400 hover:text-gray-600">
                <span class="material-symbols-outlined">close</span>
              </button>
            </div>
            <form onSubmit={handleSubmit} class="p-6 space-y-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Nombre del Sprint</label>
                <input
                  type="text"
                  required
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-primary focus:border-primary outline-none transition-all"
                  placeholder="Ej: Sprint 1 - Core MVP"
                  value={name}
                  onInput={(e) => setName(e.currentTarget.value)}
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Objetivo del Sprint</label>
                <textarea
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-primary focus:border-primary outline-none transition-all min-h-[80px]"
                  placeholder="¿Qué se espera lograr en este ciclo?"
                  value={description}
                  onInput={(e) => setDescription(e.currentTarget.value)}
                ></textarea>
              </div>

              <div class="grid grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Fecha Inicio</label>
                  <input
                    type="date"
                    required
                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-primary focus:border-primary outline-none transition-all"
                    value={startDate}
                    onInput={(e) => setStartDate(e.currentTarget.value)}
                  />
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Fecha Fin</label>
                  <input
                    type="date"
                    required
                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-primary focus:border-primary outline-none transition-all"
                    value={endDate}
                    onInput={(e) => setEndDate(e.currentTarget.value)}
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
                  {loading ? "Creando..." : "Crear Sprint"}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </>
  );
}
