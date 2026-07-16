import { Badge } from "@/components/ui/badge";

interface StatusBadgeProps {
  status: number; // 0=inactive, 1=active
}

export function StatusBadge({ status }: StatusBadgeProps) {
  if (status === 1) {
    return <Badge variant="default">Active</Badge>;
  }
  return <Badge variant="secondary">Inactive</Badge>;
}
