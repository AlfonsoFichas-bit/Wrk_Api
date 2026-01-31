export interface User {
  id: string;
  name: string;
  email: string;
  role: string;
  active?: boolean;
  createdAt?: string;
}

export interface Project {
  id: string;
  name: string;
  description?: string;
  status: string;
  startDate?: string;
  endDate?: string;
  ownerId: string;
  owner?: User;
  members?: ProjectMember[];
}

export interface ProjectMember {
  id: string;
  projectId: string;
  userId: string;
  role: string;
  user?: User;
}

export interface AuthResponse {
  message: string;
  token: string;
  user: User;
}

export interface ApiResponse<T> {
  data: T;
  message?: string;
}

// In a real Fresh project, you might use environment variables.
// For browser-side code, we default to localhost.
const BASE_URL = "http://localhost:5000/api";

export async function apiFetch<T = any>(endpoint: string, options: RequestInit = {}): Promise<T> {
  let token = null;
  if (typeof localStorage !== "undefined") {
    token = localStorage.getItem("token");
  }

  const headers = {
    "Content-Type": "application/json",
    ...options.headers,
  } as Record<string, string>;

  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const response = await fetch(`${BASE_URL}${endpoint}`, {
    ...options,
    headers,
  });

  if (response.status === 204) {
    return null as any;
  }

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}));
    throw new Error(errorData.error || "Something went wrong");
  }

  return response.json();
}

export const projects = {
  async getAll(memberId?: string): Promise<Project[]> {
    const url = memberId ? `/projects?memberId=${memberId}` : "/projects";
    return apiFetch<Project[]>(url);
  },
  async getById(id: string): Promise<Project> {
    const res = await apiFetch<ApiResponse<Project>>(`/projects/${id}`);
    return res.data;
  },
  async create(projectData: Partial<Project>): Promise<Project> {
    const res = await apiFetch<ApiResponse<Project>>("/projects/", {
      method: "POST",
      body: JSON.stringify(projectData),
    });
    return res.data;
  }
};

export const auth = {
  async login(credentials: any): Promise<AuthResponse> {
    const data = await apiFetch<AuthResponse>("/auth/login", {
      method: "POST",
      body: JSON.stringify(credentials),
    });
    if (typeof localStorage !== "undefined") {
      localStorage.setItem("token", data.token);
      localStorage.setItem("user", JSON.stringify(data.user));
    }
    return data;
  },
  async register(userData: any): Promise<any> {
    return apiFetch("/auth/register", {
      method: "POST",
      body: JSON.stringify(userData),
    });
  },
  logout() {
    if (typeof localStorage !== "undefined") {
      localStorage.removeItem("token");
      localStorage.removeItem("user");
    }
  },
  getUser(): User | null {
    if (typeof localStorage !== "undefined") {
      const user = localStorage.getItem("user");
      return user ? JSON.parse(user) : null;
    }
    return null;
  }
};
