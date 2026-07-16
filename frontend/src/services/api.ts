const BASE_URL = "http://localhost:3000/api/v1";

class ApiError extends Error {
  status: number;
  errors: Record<string, string> | null;

  constructor(status: number, message: string, errors: Record<string, string> | null = null) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.errors = errors;
  }
}

async function request<T>(url: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE_URL}${url}`, {
    headers: {
      "Content-Type": "application/json",
    },
    ...options,
  });

  if (res.status === 204) {
    return undefined as T;
  }

  const data = await res.json();

  if (!res.ok) {
    const message = data?.message ?? data?.error ?? `Request failed with status ${res.status}`;
    const errors = data?.errors ?? null;
    throw new ApiError(res.status, message, errors);
  }

  return data as T;
}

// ── Types ──────────────────────────────────────────────

export interface Ingredient {
  uuid: string;
  name: string;
  cause_alergy: boolean;
  type: number; // 0=none, 1=veggie, 2=vegan
  status: number; // 0=inactive, 1=active
  created_at: string;
  updated_at: string | null;
  deleted_at: string | null;
}

export interface Item {
  uuid: string;
  name: string;
  price: number;
  status: number;
  created_at: string;
  updated_at: string | null;
  deleted_at: string | null;
  ingredients: string[];
}

export interface ItemIngredient {
  uuid_item: string;
  uuid_ingredient: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  page: number;
  page_size: number;
  total: number;
}

export interface CreateIngredientPayload {
  name: string;
  cause_alergy: boolean;
  type: number;
  status: number;
}

export interface UpdateIngredientPayload {
  name: string;
  cause_alergy: boolean;
  type: number;
  status: number;
}

export interface CreateItemPayload {
  name: string;
  price: number;
  status: number;
  ingredients: string[];
}

export interface UpdateItemPayload {
  name: string;
  price: number;
  status: number;
  ingredients: string[];
}

// ── Ingredient API ─────────────────────────────────────

export const ingredientApi = {
  list(page = 1, pageSize = 10) {
    return request<PaginatedResponse<Ingredient>>(
      `/ingredients?page=${page}&page_size=${pageSize}`
    );
  },

  get(uuid: string) {
    return request<Ingredient>(`/ingredients/${uuid}`);
  },

  create(data: CreateIngredientPayload) {
    return request<Ingredient>("/ingredients", {
      method: "POST",
      body: JSON.stringify(data),
    });
  },

  update(uuid: string, data: UpdateIngredientPayload) {
    return request<Ingredient>(`/ingredients/${uuid}`, {
      method: "PUT",
      body: JSON.stringify(data),
    });
  },

  delete(uuid: string) {
    return request<void>(`/ingredients/${uuid}`, {
      method: "DELETE",
    });
  },
};

// ── Item API ───────────────────────────────────────────

export const itemApi = {
  list(page = 1, pageSize = 10) {
    return request<PaginatedResponse<Item>>(
      `/items?page=${page}&page_size=${pageSize}`
    );
  },

  get(uuid: string) {
    return request<Item>(`/items/${uuid}`);
  },

  create(data: CreateItemPayload) {
    return request<Item>("/items", {
      method: "POST",
      body: JSON.stringify(data),
    });
  },

  update(uuid: string, data: UpdateItemPayload) {
    return request<Item>(`/items/${uuid}`, {
      method: "PUT",
      body: JSON.stringify(data),
    });
  },

  delete(uuid: string) {
    return request<void>(`/items/${uuid}`, {
      method: "DELETE",
    });
  },
};

// ── Item-Ingredient API ────────────────────────────────

export const itemIngredientApi = {
  list(itemUuid: string) {
    return request<ItemIngredient[]>(`/items/${itemUuid}/ingredients`);
  },
};

export { ApiError, BASE_URL };
