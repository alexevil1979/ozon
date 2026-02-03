import { AdminProductList } from '@/components/admin/AdminProductList';

export default function AdminProductsPage() {
  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Товары (модерация)</h1>
      <AdminProductList />
    </div>
  );
}
