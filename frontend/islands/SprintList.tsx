import { type Sprint } from "../api_client.ts";

interface SprintListProps {
  sprints: Sprint[];
}

export default function SprintList({ sprints }: SprintListProps) {
  return (
    <div class="bg-white rounded-2xl border border-gray-200 shadow-sm overflow-hidden">
      <div class="px-6 py-4 border-b border-gray-100 bg-gray-50 flex justify-between items-center">
        <h2 class="font-bold text-gray-900">Sprints del Proyecto</h2>
        <span class="text-xs font-bold bg-blue-100 text-blue-700 px-2 py-1 rounded-full">
          {sprints.length} total
        </span>
      </div>

      <div class="divide-y divide-gray-100">
        {sprints.length === 0 ? (
          <div class="p-8 text-center text-gray-500 italic">
            No hay sprints configurados. Crea uno para empezar a planificar.
          </div>
        ) : (
          sprints.map((sprint) => (
            <div key={sprint.id} class="p-6 hover:bg-gray-50 transition-colors">
              <div class="flex justify-between items-start mb-2">
                <div>
                  <h3 class="font-bold text-gray-900 text-lg">{sprint.name}</h3>
                  <p class="text-sm text-gray-500">
                    {sprint.startDate ? new Date(sprint.startDate).toLocaleDateString() : 'N/A'} - {sprint.endDate ? new Date(sprint.endDate).toLocaleDateString() : 'N/A'}
                  </p>
                </div>
                <span class={`text-xs font-bold px-2 py-1 rounded-full ${
                  sprint.status === 'ACTIVE' ? 'bg-green-100 text-green-700' :
                  sprint.status === 'PLANNING' ? 'bg-orange-100 text-orange-700' :
                  'bg-gray-100 text-gray-700'
                }`}>
                  {sprint.status}
                </span>
              </div>
              <p class="text-sm text-gray-600 mb-4 line-clamp-2">
                {sprint.description || "Sin descripción"}
              </p>
              <div class="flex justify-between items-center">
                <div class="flex -space-x-2">
                   <span class="text-xs text-gray-400 font-medium">
                     {sprint.userStories?.length || 0} Historias de Usuario
                   </span>
                </div>
                <button
                  class="text-primary text-sm font-bold hover:underline flex items-center gap-1"
                >
                  Ver Tablero
                  <span class="material-symbols-outlined text-sm">arrow_forward</span>
                </button>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
