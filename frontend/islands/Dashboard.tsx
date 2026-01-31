import { useEffect, useState } from "preact/hooks";
import { apiFetch, auth, type User } from "../api_client.ts";

export default function Dashboard() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [currentUser, setCurrentUser] = useState<User | null>(null);

  useEffect(() => {
    const user = auth.getUser();
    if (!user) {
      window.location.href = "/login";
      return;
    }
    setCurrentUser(user);

    apiFetch("/users")
      .then((res) => {
        if (Array.isArray(res.data)) {
          setUsers(res.data);
        } else if (Array.isArray(res)) {
          setUsers(res);
        }
      })
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, []);

  const handleLogout = () => {
    auth.logout();
    window.location.href = "/login";
  };

  if (loading) return <div class="p-8 text-center">Cargando...</div>;

  return (
    <div class="p-8 max-w-6xl mx-auto">
      <div class="flex justify-between items-center mb-8 bg-white p-6 rounded-lg shadow-sm border border-gray-100">
        <h1 class="text-3xl font-bold text-gray-800">Panel de Control</h1>
        <div class="flex items-center gap-4">
          <div class="text-right">
            <p class="font-semibold text-gray-900">{currentUser?.name}</p>
            <p class="text-xs text-gray-500 uppercase">{currentUser?.role}</p>
          </div>
          <button
            onClick={handleLogout}
            class="px-4 py-2 bg-red-500 text-white rounded-md hover:bg-red-600 transition-colors font-medium"
          >
            Cerrar Sesión
          </button>
        </div>
      </div>

      {error && (
        <div class="bg-red-50 border-l-4 border-red-400 p-4 mb-8">
          <div class="flex">
            <div class="flex-shrink-0">
              <span class="text-red-400">⚠️</span>
            </div>
            <div class="ml-3">
              <p class="text-sm text-red-700">{error}</p>
            </div>
          </div>
        </div>
      )}

      <div class="grid grid-cols-1 md:grid-cols-2 gap-8 mb-8">
        <a href="/dashboard/projects" class="group bg-primary p-8 rounded-2xl shadow-lg shadow-primary/20 flex items-center justify-between text-white transition-transform hover:scale-[1.02]">
          <div>
            <h3 class="text-2xl font-bold mb-2">Mis Proyectos</h3>
            <p class="text-white/80">Gestiona tus espacios de trabajo y Sprints.</p>
          </div>
          <span class="material-symbols-outlined text-4xl opacity-50 group-hover:opacity-100 transition-opacity">folder_shared</span>
        </a>
        <div class="bg-white p-8 rounded-2xl border border-gray-200 flex items-center justify-between transition-transform hover:scale-[1.02]">
          <div>
            <h3 class="text-2xl font-bold text-gray-900 mb-2">Mi Perfil</h3>
            <p class="text-gray-500">Configura tus preferencias y rol.</p>
          </div>
          <span class="material-symbols-outlined text-4xl text-gray-300">account_circle</span>
        </div>
      </div>

      <div class="bg-white shadow rounded-lg overflow-hidden border border-gray-200">
        <div class="px-6 py-4 border-b border-gray-200 bg-gray-50">
          <h2 class="text-lg font-semibold text-gray-800">Usuarios del Sistema</h2>
        </div>
        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Nombre</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Email</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Rol</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              {users.length > 0 ? users.map((user: any) => (
                <tr key={user.id} class="hover:bg-gray-50 transition-colors">
                  <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{user.name}</td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{user.email}</td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    <span class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-blue-100 text-blue-800">
                      {user.role}
                    </span>
                  </td>
                </tr>
              )) : (
                <tr>
                  <td colspan={3} class="px-6 py-10 text-center text-gray-500 italic">No hay usuarios disponibles</td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
