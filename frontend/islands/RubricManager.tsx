import { useState, useEffect } from "preact/hooks";
import { rubrics, auth, type Rubric } from "../api_client.ts";

interface RubricManagerProps {
  projectId?: string;
}

export default function RubricManager({ projectId }: RubricManagerProps) {
  const [rubricList, setRubricList] = useState<Rubric[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);

  // Form State
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [criteria, setCriteria] = useState([{ name: "", description: "", maxScore: 10, weight: 1 }]);

  const user = auth.getUser();
  const isTeacher = user?.role === "ADMIN" || user?.role === "DOCENTE";

  useEffect(() => {
    fetchRubrics();
  }, [projectId]);

  const fetchRubrics = async () => {
    try {
      setLoading(true);
      const data = await rubrics.getAll(projectId);
      setRubricList(data);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const addCriteria = () => {
    setCriteria([...criteria, { name: "", description: "", maxScore: 10, weight: 1 }]);
  };

  const removeCriteria = (index: number) => {
    setCriteria(criteria.filter((_, i) => i !== index));
  };

  const updateCriteria = (index: number, field: string, value: any) => {
    const newCriteria = [...criteria];
    (newCriteria[index] as any)[field] = value;
    setCriteria(newCriteria);
  };

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    try {
      await rubrics.create({
        projectId,
        name,
        description,
        criteria
      });
      setShowForm(false);
      setName("");
      setDescription("");
      setCriteria([{ name: "", description: "", maxScore: 10, weight: 1 }]);
      fetchRubrics();
    } catch (err) {
      alert("Error al crear rúbrica");
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm("¿Eliminar esta rúbrica?")) return;
    try {
      await rubrics.delete(id);
      fetchRubrics();
    } catch (err) {
      alert("Error al eliminar");
    }
  };

  if (loading) return <div class="p-4 text-center">Cargando rúbricas...</div>;

  return (
    <div class="space-y-6">
      <div class="flex justify-between items-center">
        <h2 class="text-xl font-bold text-gray-900">Rúbricas de Evaluación</h2>
        {isTeacher && !showForm && (
          <button
            onClick={() => setShowForm(true)}
            class="px-4 py-2 bg-primary text-white rounded-lg text-sm font-bold hover:bg-primary-hover transition-colors"
          >
            Nueva Rúbrica
          </button>
        )}
      </div>

      {showForm && (
        <div class="bg-white p-6 rounded-2xl border border-primary/20 shadow-sm space-y-4">
          <div class="flex justify-between items-center">
            <h3 class="font-bold text-gray-800">Configurar Rúbrica</h3>
            <button onClick={() => setShowForm(false)} class="text-gray-400 hover:text-gray-600">
              <span class="material-symbols-outlined">close</span>
            </button>
          </div>
          <form onSubmit={handleSubmit} class="space-y-4">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Nombre</label>
                <input
                  type="text" required class="w-full px-3 py-2 border rounded-lg"
                  value={name} onInput={e => setName(e.currentTarget.value)}
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Descripción</label>
                <input
                  type="text" class="w-full px-3 py-2 border rounded-lg"
                  value={description} onInput={e => setDescription(e.currentTarget.value)}
                />
              </div>
            </div>

            <div class="space-y-3">
              <div class="flex justify-between items-center">
                <h4 class="text-sm font-bold text-gray-700 uppercase">Criterios</h4>
                <button type="button" onClick={addCriteria} class="text-primary text-sm font-bold">+ Añadir Criterio</button>
              </div>
              {criteria.map((c, i) => (
                <div key={i} class="p-4 bg-gray-50 rounded-xl border border-gray-100 flex flex-col md:flex-row gap-3 items-end">
                  <div class="flex-1 w-full">
                    <label class="block text-[10px] uppercase font-bold text-gray-500">Nombre del Criterio</label>
                    <input
                      type="text" required class="w-full px-2 py-1 border rounded"
                      value={c.name} onInput={e => updateCriteria(i, "name", e.currentTarget.value)}
                    />
                  </div>
                  <div class="w-24">
                    <label class="block text-[10px] uppercase font-bold text-gray-500">Puntaje Máx</label>
                    <input
                      type="number" required class="w-full px-2 py-1 border rounded"
                      value={c.maxScore} onInput={e => updateCriteria(i, "maxScore", parseInt(e.currentTarget.value))}
                    />
                  </div>
                  <button type="button" onClick={() => removeCriteria(i)} class="p-1 text-red-400 hover:text-red-600">
                    <span class="material-symbols-outlined">delete</span>
                  </button>
                </div>
              ))}
            </div>

            <button type="submit" class="w-full py-2 bg-primary text-white rounded-lg font-bold">Guardar Rúbrica</button>
          </form>
        </div>
      )}

      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        {rubricList.map(rubric => (
          <div key={rubric.id} class="bg-white p-5 rounded-2xl border border-gray-100 shadow-sm group">
            <div class="flex justify-between items-start mb-3">
              <div>
                <h4 class="font-bold text-gray-900">{rubric.name}</h4>
                <p class="text-xs text-gray-500">{rubric.description}</p>
              </div>
              {isTeacher && (
                <button onClick={() => handleDelete(rubric.id)} class="text-gray-300 hover:text-red-500 opacity-0 group-hover:opacity-100 transition-all">
                  <span class="material-symbols-outlined text-[20px]">delete</span>
                </button>
              )}
            </div>
            <div class="flex flex-wrap gap-2">
              {rubric.criteria?.map(c => (
                <span key={c.id} class="text-[10px] bg-blue-50 text-blue-600 px-2 py-1 rounded-full font-bold uppercase">
                  {c.name} ({c.maxScore} pts)
                </span>
              ))}
            </div>
          </div>
        ))}
        {rubricList.length === 0 && !showForm && (
          <div class="col-span-2 p-10 text-center bg-gray-50 rounded-2xl border border-dashed border-gray-200 text-gray-500 italic">
            No hay rúbricas definidas para este proyecto.
          </div>
        )}
      </div>
    </div>
  );
}
