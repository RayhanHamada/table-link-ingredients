import { Badge } from "@/components/ui/badge";

const typeLabels: Record<number, string> = {
  0: "None",
  1: "Veggie",
  2: "Vegan",
};

const typeVariants: Record<number, "default" | "secondary" | "outline"> = {
  0: "secondary",
  1: "outline",
  2: "default",
};

interface TypeBadgeProps {
  type: number; // 0=none, 1=veggie, 2=vegan
}

export function TypeBadge({ type }: TypeBadgeProps) {
  return (
    <Badge variant={typeVariants[type] ?? "secondary"}>
      {typeLabels[type] ?? "Unknown"}
    </Badge>
  );
}
