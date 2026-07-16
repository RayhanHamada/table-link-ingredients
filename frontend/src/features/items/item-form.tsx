import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { useForm } from "@tanstack/react-form";
import { useQuery } from "@tanstack/react-query";
import { toast } from "sonner";
import * as v from "valibot";

import { useItem, useCreateItem, useUpdateItem } from "@/hooks/use-items";
import { ingredientApi, ApiError } from "@/services/api";
import { ingredientKeys } from "@/hooks/query-keys";
import { valibotField, type ItemFormValues } from "@/lib/validations";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Spinner } from "@/components/ui/spinner";
import { ErrorState } from "@/components/shared/error-state";
import {
  Combobox,
  ComboboxChips,
  ComboboxChip,
  ComboboxChipsInput,
  ComboboxContent,
  ComboboxList,
  ComboboxItem,
  ComboboxCollection,
  ComboboxEmpty,
} from "@/components/ui/combobox";

export function ItemForm() {
  const navigate = useNavigate();
  const { uuid } = useParams<{ uuid: string }>();
  const isEdit = !!uuid;

  const { data: item, isLoading, isError, refetch } = useItem(uuid);
  const createMutation = useCreateItem();
  const updateMutation = useUpdateItem();

  // Fetch all ingredients for the multi-select
  const [ingredientPage] = useState(1);
  const [ingredientPageSize] = useState(100);
  const { data: ingredientData } = useQuery({
    queryKey: ingredientKeys.list(ingredientPage, ingredientPageSize),
    queryFn: () => ingredientApi.list(ingredientPage, ingredientPageSize),
    staleTime: 60_000,
  });

  const availableIngredients = ingredientData?.data ?? [];

  const form = useForm<ItemFormValues>({
    defaultValues: {
      name: "",
      price: 0,
      status: 1,
      ingredients: [] as string[],
    },
    onSubmit: async ({ value }) => {
      try {
        if (isEdit && uuid) {
          await updateMutation.mutateAsync({ uuid, data: value });
          toast.success(`"${value.name}" updated successfully`);
        } else {
          await createMutation.mutateAsync(value);
          toast.success(`"${value.name}" created successfully`);
        }
        navigate("/items");
      } catch (err) {
        if (err instanceof ApiError && err.errors) {
          for (const [field, message] of Object.entries(err.errors)) {
            if (field in value) {
              form.setFieldMeta(field as keyof ItemFormValues, (prev) => ({
                ...prev,
                errorMap: {
                  ...prev.errorMap,
                  onServer: message,
                },
              }));
            }
          }
        }
        const message = err instanceof ApiError ? err.message : "Something went wrong";
        toast.error(message);
      }
    },
  });

  useEffect(() => {
    if (item && isEdit) {
      form.setFieldValue("name", item.name);
      form.setFieldValue("price", item.price);
      form.setFieldValue("status", item.status);
      form.setFieldValue("ingredients", item.ingredients ?? []);
    }
  }, [item, isEdit]);

  const STATUS_OPTIONS = [
    { value: "0", label: "Inactive" },
    { value: "1", label: "Active" },
  ];

  const isPending = createMutation.isPending || updateMutation.isPending;

  if (isEdit && isLoading) {
    return (
      <div className="flex items-center justify-center py-16">
        <Spinner className="size-8" />
      </div>
    );
  }

  if (isEdit && isError) {
    return <ErrorState message="Failed to load item" onRetry={() => refetch()} />;
  }

  const getIngredientName = (ingUuid: string) => {
    const found = availableIngredients.find((ing) => ing.uuid === ingUuid);
    return found?.name ?? ingUuid;
  };

  return (
    <div className="mx-auto max-w-lg space-y-6">
      <div>
        <h1 className="font-heading text-2xl font-semibold tracking-tight">
          {isEdit ? "Edit Item" : "New Item"}
        </h1>
        <p className="mt-1 text-sm text-muted-foreground">
          {isEdit ? "Update item details" : "Add a new item to the menu"}
        </p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Item Details</CardTitle>
        </CardHeader>
        <CardContent>
          <form
            onSubmit={(e) => {
              e.preventDefault();
              form.handleSubmit();
            }}
            className="space-y-5"
          >
            <form.Field
              name="name"
              validators={{
                onChange: valibotField(v.pipe(v.string(), v.nonEmpty("Name is required"))),
              }}
            >
              {(field) => (
                <div className="space-y-2">
                  <Label htmlFor="name">Name</Label>
                  <Input
                    id="name"
                    value={field.state.value}
                    onChange={(e) => field.handleChange(e.target.value)}
                    placeholder="Enter item name"
                  />
                  {field.state.meta.errors.length > 0 && (
                    <p className="text-sm text-destructive">
                      {field.state.meta.errors.join(", ")}
                    </p>
                  )}
                </div>
              )}
            </form.Field>

            <form.Field
              name="price"
              validators={{
                onChange: valibotField(v.pipe(v.number(), v.minValue(1, "Price must be greater than 0"))),
              }}
            >
              {(field) => (
                <div className="space-y-2">
                  <Label htmlFor="price">Price (IDR)</Label>
                  <Input
                    id="price"
                    type="number"
                    min={0}
                    value={field.state.value}
                    onChange={(e) => field.handleChange(Number(e.target.value))}
                    placeholder="Enter price"
                  />
                  {field.state.meta.errors.length > 0 && (
                    <p className="text-sm text-destructive">
                      {field.state.meta.errors.join(", ")}
                    </p>
                  )}
                </div>
              )}
            </form.Field>

            <form.Field
              name="status"
              validators={{
                onChange: valibotField(v.pipe(v.number(), v.integer(), v.minValue(0), v.maxValue(1))),
              }}
            >
              {(field) => (
                <div className="space-y-2">
                  <Label htmlFor="status">Status</Label>
                  <Select
                    value={String(field.state.value)}
                    onValueChange={(v) => field.handleChange(Number(v))}
                  >
                    <SelectTrigger id="status" className="w-full">
                      <SelectValue>
                        {STATUS_OPTIONS.find((o) => o.value === String(field.state.value))?.label ?? "Select status"}
                      </SelectValue>
                    </SelectTrigger>
                    <SelectContent>
                      {STATUS_OPTIONS.map((opt) => (
                        <SelectItem key={opt.value} value={opt.value}>
                          {opt.label}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  {field.state.meta.errors.length > 0 && (
                    <p className="text-sm text-destructive">
                      {field.state.meta.errors.join(", ")}
                    </p>
                  )}
                </div>
              )}
            </form.Field>

            <form.Field
              name="ingredients"
              validators={{
                onChange: valibotField(v.pipe(v.array(v.string()), v.minLength(1, "At least one ingredient is required"))),
              }}
            >
              {(field) => (
                <div className="space-y-2">
                  <Label>Ingredients</Label>
                  <Combobox
                    multiple
                    value={field.state.value}
                    onValueChange={(values) => {
                      field.handleChange(
                        Array.isArray(values) ? (values as string[]) : []
                      );
                    }}
                  >
                    <ComboboxChips>
                      {field.state.value.map((ingUuid: string) => (
                        <ComboboxChip key={ingUuid} value={ingUuid}>
                          {getIngredientName(ingUuid)}
                        </ComboboxChip>
                      ))}
                      <ComboboxChipsInput placeholder="Select ingredients..." />
                    </ComboboxChips>
                    <ComboboxContent>
                      <ComboboxList>
                        <ComboboxCollection>
                          {availableIngredients.map((ing) => (
                            <ComboboxItem key={ing.uuid} value={ing.uuid}>
                              {ing.name}
                            </ComboboxItem>
                          ))}
                        </ComboboxCollection>
                        {availableIngredients.length === 0 && (
                          <ComboboxEmpty>No ingredients found</ComboboxEmpty>
                        )}
                      </ComboboxList>
                    </ComboboxContent>
                  </Combobox>
                  {field.state.meta.errors.length > 0 && (
                    <p className="text-sm text-destructive">
                      {field.state.meta.errors.join(", ")}
                    </p>
                  )}
                </div>
              )}
            </form.Field>

            <div className="flex items-center gap-3 pt-2">
              <Button type="submit" disabled={isPending}>
                {isPending ? (
                  <>
                    <Spinner className="size-4" />
                    Saving...
                  </>
                ) : isEdit ? (
                  "Update Item"
                ) : (
                  "Create Item"
                )}
              </Button>
              <Button
                type="button"
                variant="outline"
                onClick={() => navigate("/items")}
                disabled={isPending}
              >
                Cancel
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}

