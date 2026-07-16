import { useEffect } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { useForm } from "@tanstack/react-form";
import { toast } from "sonner";
import * as v from "valibot";

import { useIngredient, useCreateIngredient, useUpdateIngredient } from "@/hooks/use-ingredients";
import { ApiError } from "@/services/api";
import { valibotField, type IngredientFormValues } from "@/lib/validations";

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
import { Checkbox } from "@/components/ui/checkbox";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Spinner } from "@/components/ui/spinner";

export function IngredientForm() {
  const navigate = useNavigate();
  const { uuid } = useParams<{ uuid: string }>();
  const isEdit = !!uuid;

  const { data: ingredient, isLoading, isError, refetch } = useIngredient(uuid);
  const createMutation = useCreateIngredient();
  const updateMutation = useUpdateIngredient();

  const form = useForm<IngredientFormValues>({
    defaultValues: {
      name: "",
      cause_alergy: false,
      type: 0,
      status: 1,
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
        navigate("/ingredients");
      } catch (err) {
        if (err instanceof ApiError && err.errors) {
          for (const [field, message] of Object.entries(err.errors)) {
            if (field in value) {
              form.setFieldMeta(field as keyof IngredientFormValues, (prev) => ({
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
    if (ingredient && isEdit) {
      form.setFieldValue("name", ingredient.name);
      form.setFieldValue("cause_alergy", ingredient.cause_alergy);
      form.setFieldValue("type", ingredient.type);
      form.setFieldValue("status", ingredient.status);
    }
  }, [ingredient, isEdit]);

  const TYPE_OPTIONS = [
    { value: "0", label: "None" },
    { value: "1", label: "Veggie" },
    { value: "2", label: "Vegan" },
  ];
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
    return <ErrorState message="Failed to load ingredient" onRetry={() => refetch()} />;
  }

  return (
    <div className="mx-auto max-w-lg space-y-6">
      <div>
        <h1 className="font-heading text-2xl font-semibold tracking-tight">
          {isEdit ? "Edit Ingredient" : "New Ingredient"}
        </h1>
        <p className="mt-1 text-sm text-muted-foreground">
          {isEdit ? "Update ingredient details" : "Add a new ingredient to the catalog"}
        </p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Ingredient Details</CardTitle>
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
                    placeholder="Enter ingredient name"
                  />
                  {field.state.meta.errors.length > 0 && (
                    <p className="text-sm text-destructive">
                      {field.state.meta.errors.join(", ")}
                    </p>
                  )}
                </div>
              )}
            </form.Field>

            <form.Field name="cause_alergy">
              {(field) => (
                <div className="flex items-center gap-3">
                  <Checkbox
                    id="cause_alergy"
                    checked={field.state.value}
                    onCheckedChange={(checked) =>
                      field.handleChange(checked === true)
                    }
                  />
                  <Label htmlFor="cause_alergy" className="cursor-pointer">
                    Causes allergy
                  </Label>
                </div>
              )}
            </form.Field>

            <form.Field
              name="type"
              validators={{
                onChange: valibotField(v.pipe(v.number(), v.integer(), v.minValue(0), v.maxValue(2))),
              }}
            >
              {(field) => (
                <div className="space-y-2">
                  <Label htmlFor="type">Type</Label>
                  <Select
                    value={String(field.state.value)}
                    onValueChange={(v) => field.handleChange(Number(v))}
                  >
                    <SelectTrigger id="type" className="w-full">
                      <SelectValue>
                        {TYPE_OPTIONS.find((o) => o.value === String(field.state.value))?.label ?? "Select type"}
                      </SelectValue>
                    </SelectTrigger>
                    <SelectContent>
                      {TYPE_OPTIONS.map((opt) => (
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

            <div className="flex items-center gap-3 pt-2">
              <Button type="submit" disabled={isPending}>
                {isPending ? (
                  <>
                    <Spinner className="size-4" />
                    Saving...
                  </>
                ) : isEdit ? (
                  "Update Ingredient"
                ) : (
                  "Create Ingredient"
                )}
              </Button>
              <Button
                type="button"
                variant="outline"
                onClick={() => navigate("/ingredients")}
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

