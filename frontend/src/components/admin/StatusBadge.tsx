import { Badge } from '@/components/ui/badge';

type Status =
  | 'draft'
  | 'published'
  | 'full'
  | 'in_progress'
  | 'completed'
  | 'cancelled'
  | 'pending'
  | 'confirmed'
  | 'active'
  | 'inactive';

interface StatusBadgeProps {
  status: Status | string;
  label?: string;
  className?: string;
}

const statusConfig: Record<
  Status,
  { label: string; variant: 'default' | 'secondary' | 'destructive' | 'outline' }
> = {
  // Trip statuses
  draft: { label: 'Borrador', variant: 'outline' },
  published: { label: 'Publicado', variant: 'default' },
  full: { label: 'Completo', variant: 'secondary' },
  in_progress: { label: 'En Progreso', variant: 'default' },
  completed: { label: 'Completado', variant: 'secondary' },
  cancelled: { label: 'Cancelado', variant: 'destructive' },

  // Booking statuses
  pending: { label: 'Pendiente', variant: 'outline' },
  confirmed: { label: 'Confirmada', variant: 'default' },

  // General statuses
  active: { label: 'Activo', variant: 'default' },
  inactive: { label: 'Inactivo', variant: 'outline' }
};

export default function StatusBadge({ status, label, className = '' }: StatusBadgeProps) {
  const config = statusConfig[status as Status] || { label: status, variant: 'outline' as const };

  return (
    <Badge variant={config.variant} className={className}>
      {label || config.label}
    </Badge>
  );
}
