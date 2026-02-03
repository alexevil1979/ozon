import type { Metadata } from 'next';
import { ProductForm } from '@/components/seller/ProductForm';

export const metadata: Metadata = {
  title: 'Добавить товар',
};

export default function NewProductPage() {
  return (
    <div className="max-w-xl space-y-6">
      <h1 className="text-2xl font-bold">Добавить товар</h1>
      <ProductForm />
    </div>
  );
}
