import { useState, useEffect } from "preact/hooks";
import { rubrics as rubricsApi, evaluations as evalsApi, auth, type Rubric, type Evaluation } from "../api_client.ts";

interface EvaluationModuleProps {
  projectId: string;
  taskId?: string;
  sprintId?: string;
  title: string;
}

export default function EvaluationModule({ projectId, taskId, sprintId, title }: EvaluationModuleProps) {
  const [rubrics, setRubrics] = useState<Rubric[]>([]);
  const [selectedRubric, setSelectedRubric] = useState<Rubric | null>(null);
  const [scores, setScores] = useState<Record<string, number>>({});
  const [feedback, setFeedback] = useState("");
  const [loading, setLoading] = useState(false);
  const [showForm, setShowForm] = useState(false);
  const [existingEvals, setExistingEvals] = useState<Evaluation[]>([]);

  const user = auth.getUser();
  const isTeacher = user?.role === "ADMIN" || user?.role === "DOCENTE";

  useEffect(() => {
    if (showForm) {
      rubricsApi.getAll(projectId).then(setRubrics);
    }
  }, [showForm, projectId]);

  useEffect(() => {
    if (taskId) {
      evalsApi.getByTask(taskId).then(setExistingEvals);
    } else if (sprintId) {
      evalsApi.getBySprint(sprintId).then(setExistingEvals);
    }
  }, [taskId, sprintId]);

  const handleRubricSelect = (e: any) => {
    const rubric = rubrics.find(r => r.id === e.target.value);
    setSelectedRubric(rubric || null);
    if (rubric?.criteria) {
      const initialScores: Record<string, number> = {};
      rubric.criteria.forEach(c => initialScores[c.id] = c.maxScore);
      setScores(initialScores);
    }
  };

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    if (!selectedRubric || !user) return;

    setLoading(true);
    try {
      const criteriaScores = Object.entries(scores).map(([id, score]) => ({
        criteriaId: id,
        score
      }));

      const totalScore = criteriaScores.reduce((acc, curr) => acc + curr.score, 0);

      await evalsApi.create({
        projectId,
        taskId,
        sprintId,
        evaluatorId: user.id,
        feedback,
        score: totalScore,
        criteriaScores
      });

      setShowForm(false);
      window.location.reload();
    } catch (err) {
      alert("Error al guardar evaluación");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div class="space-y-4">
      <div class="flex justify-between items-center">
        <h3 class="font-bold text-gray-800 flex items-center gap-2">
          <span class="material-symbols-outlined text-primary">analytics</span>
          Evaluación
        </h3>
        {isTeacher && !showForm && (
          <button
            onClick={() => setShowForm(true)}
            class="text-xs font-bold text-primary hover:underline"
          >
            Evaluar ahora
          </button>
        )}
      </div>

      {showForm && (
        <div class="bg-gray-50 p-6 rounded-2xl border border-primary/10 space-y-4">
          <div class="flex justify-between items-center">
            <h4 class="font-bold text-sm">Evaluando: {title}</h4>
            <button onClick={() => setShowForm(false)} class="text-gray-400">
              <span class="material-symbols-outlined">close</span>
            </button>
          </div>

          <div class="space-y-4">
            <div>
              <label class="block text-xs font-bold text-gray-500 uppercase mb-1">Seleccionar Rúbrica</label>
              <select
                class="w-full px-3 py-2 border rounded-lg bg-white text-sm"
                onChange={handleRubricSelect}
              >
                <option value="">-- Elige una rúbrica --</option>
                {rubrics.map(r => <option key={r.id} value={r.id}>{r.name}</option>)}
              </select>
            </div>

            {selectedRubric && (
              <div class="space-y-4">
                <div class="space-y-3">
                  {selectedRubric.criteria?.map(c => (
                    <div key={c.id} class="bg-white p-3 rounded-xl border border-gray-100 flex justify-between items-center">
                      <div class="flex-1">
                        <p class="text-sm font-bold text-gray-800">{c.name}</p>
                        <p class="text-[10px] text-gray-500">{c.description}</p>
                      </div>
                      <div class="flex items-center gap-2">
                        <input
                          type="number" max={c.maxScore} min="0"
                          class="w-16 px-2 py-1 border rounded text-right font-bold text-primary"
                          value={scores[c.id]}
                          onInput={e => setScores({ ...scores, [c.id]: parseInt(e.currentTarget.value) || 0 })}
                        />
                        <span class="text-xs text-gray-400">/ {c.maxScore}</span>
                      </div>
                    </div>
                  ))}
                </div>

                <div>
                  <label class="block text-xs font-bold text-gray-500 uppercase mb-1">Retroalimentación (Feedback)</label>
                  <textarea
                    class="w-full px-3 py-2 border rounded-lg text-sm min-h-[80px]"
                    placeholder="Escribe tus comentarios aquí..."
                    onInput={e => setFeedback(e.currentTarget.value)}
                  ></textarea>
                </div>

                <button
                  onClick={handleSubmit}
                  disabled={loading}
                  class="w-full py-2 bg-primary text-white rounded-lg font-bold hover:bg-primary-hover transition-colors disabled:opacity-50"
                >
                  {loading ? "Guardando..." : "Confirmar Evaluación"}
                </button>
              </div>
            )}
          </div>
        </div>
      )}

      <div class="space-y-3">
        {existingEvals.map(evalu => (
          <div key={evalu.id} class="bg-white p-4 rounded-xl border border-gray-100 shadow-sm">
            <div class="flex justify-between items-center mb-2">
              <span class="text-[10px] font-bold text-gray-400 uppercase tracking-widest">
                Evaluado por {evalu.evaluator?.name}
              </span>
              <span class="text-lg font-black text-primary">{evalu.score} pts</span>
            </div>
            <p class="text-sm text-gray-700 italic">"{evalu.feedback}"</p>
            <p class="text-[10px] text-gray-400 mt-2">
              {new Date(evalu.createdAt).toLocaleDateString()}
            </p>
          </div>
        ))}
        {existingEvals.length === 0 && !showForm && (
          <div class="p-6 text-center text-gray-400 italic text-sm">
            Aún no hay evaluaciones registradas.
          </div>
        )}
      </div>
    </div>
  );
}
