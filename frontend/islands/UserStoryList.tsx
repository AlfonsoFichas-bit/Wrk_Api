import { useState } from "preact/hooks";
import { userStories, sprints as sprintsApi, type UserStory, type Sprint } from "../api_client.ts";

interface UserStoryListProps {
  projectId: string;
  initialStories: UserStory[];
  sprints: Sprint[];
}

export default function UserStoryList({ projectId, initialStories, sprints }: UserStoryListProps) {
  const [stories, setStories] = useState<UserStory[]>(initialStories);
  const [loading, setLoading] = useState(false);

  const handleDelete = async (id: string) => {
    if (!confirm("¿Estás seguro de que deseas eliminar esta historia de usuario?")) return;

    try {
      await userStories.delete(id);
      setStories(stories.filter(s => s.id !== id));
    } catch (err) {
      alert("Error al eliminar: " + (err instanceof Error ? err.message : "Error desconocido"));
    }
  };

  const handleMoveToSprint = async (storyId: string, sprintId: string) => {
    if (!sprintId) return;
    try {
      await sprintsApi.addStory(sprintId, storyId);
      window.location.reload();
    } catch (err) {
      alert("Error al mover: " + (err instanceof Error ? err.message : "Error desconocido"));
    }
  };

  const backlogStories = stories.filter(s => s.status === "BACKLOG" && !s.sprintId);

  return (
    <div class="bg-white rounded-2xl border border-gray-200 shadow-sm overflow-hidden">
      <div class="px-6 py-4 border-b border-gray-100 bg-gray-50 flex justify-between items-center">
        <h2 class="font-bold text-gray-900">Backlog del Proyecto</h2>
        <span class="text-xs font-bold bg-gray-200 text-gray-700 px-2 py-1 rounded-full">
          {backlogStories.length} historias
        </span>
      </div>

      <div class="divide-y divide-gray-100">
        {backlogStories.length === 0 ? (
          <div class="p-8 text-center text-gray-500 italic">
            El backlog está vacío. Empieza a añadir historias de usuario.
          </div>
        ) : (
          backlogStories.map((story) => (
            <div key={story.id} class="p-4 hover:bg-gray-50 transition-colors flex justify-between items-center group">
              <div class="flex-1">
                <div class="flex items-center gap-2 mb-1">
                  <span class={`text-[10px] font-bold px-1.5 py-0.5 rounded uppercase ${
                    story.priority === 'HIGH' ? 'bg-red-100 text-red-700' :
                    story.priority === 'MEDIUM' ? 'bg-orange-100 text-orange-700' :
                    'bg-blue-100 text-blue-700'
                  }`}>
                    {story.priority}
                  </span>
                  <span class="text-xs font-medium text-gray-400">#{story.id.slice(-4)}</span>
                </div>
                <h4 class="font-semibold text-gray-900">{story.title}</h4>
                <p class="text-sm text-gray-500 line-clamp-1">{story.description}</p>
              </div>
              <div class="flex items-center gap-4">
                <select
                  class="text-xs border border-gray-200 rounded px-2 py-1 bg-white focus:border-primary outline-none"
                  onChange={(e) => handleMoveToSprint(story.id, e.currentTarget.value)}
                  value=""
                >
                  <option value="" disabled>Mover a Sprint...</option>
                  {sprints.filter(s => s.status === 'PLANNING' || s.status === 'ACTIVE').map(s => (
                    <option key={s.id} value={s.id}>{s.name}</option>
                  ))}
                </select>
                <div class="flex items-center gap-1 bg-gray-100 px-2 py-1 rounded text-xs font-bold text-gray-600">
                  <span class="material-symbols-outlined text-xs">star</span>
                  {story.storyPoints || 0}
                </div>
                <button
                  onClick={() => handleDelete(story.id)}
                  class="text-gray-300 hover:text-red-500 transition-colors opacity-0 group-hover:opacity-100"
                >
                  <span class="material-symbols-outlined text-[20px]">delete</span>
                </button>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
