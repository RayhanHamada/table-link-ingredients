import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { itemApi, type CreateItemPayload, type UpdateItemPayload } from "@/services/api";
import { itemKeys } from "@/hooks/query-keys";

export function useItems(page: number, pageSize: number) {
  return useQuery({
    queryKey: itemKeys.list(page, pageSize),
    queryFn: () => itemApi.list(page, pageSize),
  });
}

export function useItem(uuid: string | undefined) {
  return useQuery({
    queryKey: itemKeys.detail(uuid!),
    queryFn: () => itemApi.get(uuid!),
    enabled: !!uuid,
  });
}

export function useCreateItem() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateItemPayload) => itemApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: itemKeys.lists() });
    },
  });
}

export function useUpdateItem() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ uuid, data }: { uuid: string; data: UpdateItemPayload }) =>
      itemApi.update(uuid, data),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: itemKeys.lists() });
      queryClient.invalidateQueries({ queryKey: itemKeys.detail(variables.uuid) });
    },
  });
}

export function useDeleteItem() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (uuid: string) => itemApi.delete(uuid),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: itemKeys.lists() });
    },
  });
}
