import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { ingredientApi, type CreateIngredientPayload, type UpdateIngredientPayload } from "@/services/api";
import { ingredientKeys } from "@/hooks/query-keys";

export function useIngredients(page: number, pageSize: number) {
  return useQuery({
    queryKey: ingredientKeys.list(page, pageSize),
    queryFn: () => ingredientApi.list(page, pageSize),
  });
}

export function useIngredient(uuid: string | undefined) {
  return useQuery({
    queryKey: ingredientKeys.detail(uuid!),
    queryFn: () => ingredientApi.get(uuid!),
    enabled: !!uuid,
  });
}

export function useCreateIngredient() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateIngredientPayload) => ingredientApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ingredientKeys.lists() });
    },
  });
}

export function useUpdateIngredient() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ uuid, data }: { uuid: string; data: UpdateIngredientPayload }) =>
      ingredientApi.update(uuid, data),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: ingredientKeys.lists() });
      queryClient.invalidateQueries({ queryKey: ingredientKeys.detail(variables.uuid) });
    },
  });
}

export function useDeleteIngredient() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (uuid: string) => ingredientApi.delete(uuid),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ingredientKeys.lists() });
    },
  });
}
