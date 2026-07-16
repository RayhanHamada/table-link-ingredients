import * as v from "valibot";

/**
 * Creates a TanStack Form onChange validator from a Valibot schema.
 * Returns the first issue message on failure, or undefined on success.
 */
export function valibotField<T extends v.GenericSchema>(
  schema: T
): (opts: { value: unknown }) => string | undefined {
  return ({ value }) => {
    const result = v.safeParse(schema, value);
    if (result.success) return undefined;
    return result.issues[0]?.message;
  };
}

// ── Ingredient ─────────────────────────────────────────────

export const ingredientSchema = v.object({
  name: v.pipe(v.string(), v.nonEmpty("Name is required")),
  cause_alergy: v.boolean(),
  type: v.pipe(v.number(), v.minValue(0), v.maxValue(2)),
  status: v.pipe(v.number(), v.minValue(0), v.maxValue(1)),
});

export type IngredientFormValues = v.InferInput<typeof ingredientSchema>;

// ── Item ───────────────────────────────────────────────────

export const itemSchema = v.object({
  name: v.pipe(v.string(), v.nonEmpty("Name is required")),
  price: v.pipe(v.number(), v.minValue(1, "Price must be greater than 0")),
  status: v.pipe(v.number(), v.minValue(0), v.maxValue(1)),
  ingredients: v.pipe(
    v.array(v.string()),
    v.minLength(1, "At least one ingredient is required")
  ),
});

export type ItemFormValues = v.InferInput<typeof itemSchema>;
