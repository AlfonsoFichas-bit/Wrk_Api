import { Head } from "fresh/runtime";
import SprintDetail from "../../../islands/SprintDetail.tsx";
import { define } from "../../../utils.ts";

export default define.page(function SprintDetailPage(ctx) {
  const { id } = ctx.params;

  return (
    <div class="min-h-screen bg-gray-50 p-4 md:p-8">
      <Head>
        <title>Tablero del Sprint | WorkflowS</title>
        <link
          href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:wght,FILL@100..700,0..1&display=swap"
          rel="stylesheet"
        />
      </Head>

      <div class="max-w-[1280px] mx-auto">
        <SprintDetail id={id} />
      </div>
    </div>
  );
});
