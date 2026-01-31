import { Head } from "fresh/runtime";
import Dashboard from "../islands/Dashboard.tsx";

export default function DashboardPage() {
  return (
    <div class="min-h-screen bg-gray-50">
      <Head>
        <title>Dashboard | Wrk_Api</title>
        <link
          href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:wght,FILL@100..700,0..1&display=swap"
          rel="stylesheet"
        />
      </Head>
      <Dashboard />
    </div>
  );
}
