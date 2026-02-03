import { Skeleton } from '@marketplace/ui';

export default function CatalogLoading() {
  return (
    <div className="container py-8">
      <Skeleton className="h-8 w-48 mb-6" />
      <div className="flex gap-8">
        <Skeleton className="h-64 w-64 shrink-0" />
        <div className="flex-1 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {Array.from({ length: 6 }).map((_, i) => (
            <Skeleton key={i} className="h-64" />
          ))}
        </div>
      </div>
    </div>
  );
}
