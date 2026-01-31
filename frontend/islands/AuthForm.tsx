import { useState } from "preact/hooks";
import { auth } from "../api_client.ts";

export default function AuthForm({ mode: initialMode }: { mode: "login" | "register" }) {
  const [mode, setMode] = useState(initialMode);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [name, setName] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      if (mode === "login") {
        await auth.login({ email, password });
        window.location.href = "/dashboard";
      } else {
        await auth.register({ name, email, password });
        alert("¡Registro exitoso! Por favor inicia sesión.");
        setMode("login");
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "Ha ocurrido un error");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div class="w-full max-w-md p-8 space-y-6 bg-white rounded shadow-md border border-gray-200">
      <h2 class="text-2xl font-bold text-center text-gray-900">
        {mode === "login" ? "Iniciar Sesión" : "Registrarse"}
      </h2>
      <form class="space-y-4" onSubmit={handleSubmit}>
        {mode === "register" && (
          <div>
            <label class="block text-sm font-medium text-gray-700 text-left">Nombre</label>
            <input
              type="text"
              required
              class="w-full px-3 py-2 mt-1 border rounded-md focus:ring-blue-500 focus:border-blue-500 border-gray-300"
              value={name}
              onInput={(e) => setName(e.currentTarget.value)}
            />
          </div>
        )}
        <div>
          <label class="block text-sm font-medium text-gray-700 text-left">Email</label>
          <input
            type="email"
            required
            class="w-full px-3 py-2 mt-1 border rounded-md focus:ring-blue-500 focus:border-blue-500 border-gray-300"
            value={email}
            onInput={(e) => setEmail(e.currentTarget.value)}
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 text-left">Contraseña</label>
          <input
            type="password"
            required
            class="w-full px-3 py-2 mt-1 border rounded-md focus:ring-blue-500 focus:border-blue-500 border-gray-300"
            value={password}
            onInput={(e) => setPassword(e.currentTarget.value)}
          />
        </div>
        {error && <p class="text-sm text-red-600 bg-red-50 p-2 rounded">{error}</p>}
        <button
          type="submit"
          disabled={loading}
          class="w-full px-4 py-2 text-white bg-blue-600 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 font-medium transition-colors"
        >
          {loading ? "Cargando..." : mode === "login" ? "Entrar" : "Crear cuenta"}
        </button>
      </form>
      <div class="text-center">
        <button
          class="text-sm text-blue-600 hover:underline"
          onClick={() => setMode(mode === "login" ? "register" : "login")}
        >
          {mode === "login"
            ? "¿No tienes cuenta? Regístrate"
            : "¿Ya tienes cuenta? Inicia sesión"}
        </button>
      </div>
    </div>
  );
}
