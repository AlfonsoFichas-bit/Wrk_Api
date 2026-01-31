import { useSignal } from "@preact/signals";
import { Head } from "fresh/runtime";
import { define } from "../utils.ts";
import Counter from "../islands/Counter.tsx";

export default define.page(function Home(ctx) {
  const count = useSignal(3);

  console.log("Shared value " + ctx.state.shared);

  return (
    <div class="px-4 py-8 mx-auto fresh-gradient min-h-screen">
      <Head>
        <title>Fresh counter</title>
      </Head>
      <div class="max-w-screen-md mx-auto flex flex-col items-center justify-center">
        <img
          class="my-6"
          src="/logo.svg"
          width="128"
          height="128"
          alt="the Fresh logo: a sliced lemon dripping with juice"
        />
        <h1 class="text-4xl font-bold">Welcome to Wrk_Api</h1>
        <p class="my-4 text-center">
          Esta es la interfaz para tu API de Gin.
        </p>
        <div class="flex gap-4 mt-8">
          <a
            href="/login"
            class="px-6 py-3 bg-blue-600 text-white rounded-lg font-bold hover:bg-blue-700 transition-colors"
          >
            Iniciar Sesión
          </a>
          <a
            href="/register"
            class="px-6 py-3 bg-white text-blue-600 border border-blue-600 rounded-lg font-bold hover:bg-blue-50 transition-colors"
          >
            Registrarse
          </a>
        </div>
        <div class="mt-12">
          <Counter count={count} />
        </div>
      </div>
    </div>
  );
});
