import { useEffect, useState } from "preact/hooks";
import { projects, auth, type Project } from "../api_client.ts";

export default function ProjectList() {
  const [projectList, setProjectList] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    const user = auth.getUser();
    if (!user) return;

    projects.getAll(user.id)
      .then((data) => {
        if (Array.isArray(data)) {
          setProjectList(data);
        }
      })
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <div class="text-center py-10">Cargando proyectos...</div>;

  return (
    <div class="space-y-6">
      <div class="flex justify-between items-center">
        <h2 class="text-2xl font-bold text-gray-800">Mis Proyectos</h2>
      </div>

      {error && (
        <div class="p-4 bg-red-50 text-red-700 rounded-lg border border-red-200">
          {error}
        </div>
      )}

      {projectList.length === 0 ? (
        <div class="bg-white p-10 text-center rounded-xl border border-dashed border-gray-300">
          <span class="material-symbols-outlined text-gray-400 text-5xl mb-4">folder_open</span>
          <p class="text-gray-500">No tienes proyectos asignados actualmente.</p>
        </div>
      ) : (
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {projectList.map((project) => (
            <a
              key={project.id}
              href={`/dashboard/projects/${project.id}`}
              class="block bg-white p-6 rounded-xl border border-gray-200 hover:border-primary/50 hover:shadow-lg transition-all group"
            >
              <div class="flex justify-between items-start mb-4">
                <div class="p-2 bg-primary/10 rounded-lg text-primary group-hover:bg-primary group-hover:text-white transition-colors">
                  <span class="material-symbols-outlined">dataset</span>
                </div>
                <span class={`text-xs font-bold px-2 py-1 rounded-full ${
                  project.status === 'ACTIVE' ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-700'
                }`}>
                  {project.status}
                </span>
              </div>
              <h3 class="text-lg font-bold text-gray-900 mb-2 group-hover:text-primary transition-colors">
                {project.name}
              </h3>
              <p class="text-sm text-gray-500 line-clamp-2 mb-4">
                {project.description || "Sin descripción"}
              </p>
              <div class="flex items-center gap-2 mt-auto pt-4 border-t border-gray-50">
                <span class="material-symbols-outlined text-sm text-gray-400">person</span>
                <span class="text-xs text-gray-500">Propietario: {project.owner?.name || "N/A"}</span>
              </div>
            </a>
          ))}
        </div>
      )}
    </div>
  );
}
