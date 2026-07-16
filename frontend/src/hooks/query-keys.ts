export const ingredientKeys = {
  all: ["ingredients"] as const,
  lists: () => [...ingredientKeys.all, "list"] as const,
  list: (page: number, pageSize: number) =>
    [...ingredientKeys.lists(), { page, pageSize }] as const,
  details: () => [...ingredientKeys.all, "detail"] as const,
  detail: (uuid: string) => [...ingredientKeys.details(), uuid] as const,
};

export const itemKeys = {
  all: ["items"] as const,
  lists: () => [...itemKeys.all, "list"] as const,
  list: (page: number, pageSize: number) =>
    [...itemKeys.lists(), { page, pageSize }] as const,
  details: () => [...itemKeys.all, "detail"] as const,
  detail: (uuid: string) => [...itemKeys.details(), uuid] as const,
};

export const itemIngredientKeys = {
  all: ["item-ingredients"] as const,
  byItem: (itemUuid: string) => [...itemIngredientKeys.all, itemUuid] as const,
};
