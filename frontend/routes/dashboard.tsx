import { Head } from "fresh/runtime";
import Dashboard from "../islands/Dashboard.tsx";

export default function DashboardPage() {
  return (
    <div class="min-h-screen bg-gray-50">
      <Head>
        <title>Dashboard | Wrk_Api</title>
      </Head>
      <Dashboard />
    </div>
  );
}
