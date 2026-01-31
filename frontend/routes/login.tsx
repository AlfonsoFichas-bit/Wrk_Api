import { Head } from "fresh/runtime";
import AuthForm from "../islands/AuthForm.tsx";

export default function LoginPage() {
  return (
    <div class="min-h-screen bg-gray-100 flex flex-col items-center justify-center p-4">
      <Head>
        <title>Iniciar Sesión</title>
      </Head>
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-800">Wrk_Api Frontend</h1>
      </div>
      <AuthForm mode="login" />
      <a href="/" class="mt-4 text-sm text-gray-500 hover:underline">Volver al inicio</a>
    </div>
  );
}
