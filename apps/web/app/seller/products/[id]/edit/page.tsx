import type { Metadata } from 'next';
import { notFound } from 'next/navigation';
import { ProductForm } from '@/components/seller/ProductForm';

const SELLER_API = process.env.SELLER_API_URL ?? 'http://localhost:8084';

async function getProduct(id: string, token: string) {
  const res = await fetch(`${SELLER_API}/seller/products/${id}`, {
    headers: { Authorization: `Bearer ${token}` },
    cache: 'no-store',
  });
  if (res.status === 404) return null;
  if (!res.ok) return null;
  return res.json();
}

export default async function EditProductPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;
  const { cookies } = await import('next/headers');
  const token = (await cookies()).get('marketplace_token')?.value;
  if (!token) return null;
  const product = await getProduct(id, token);
  if (!product) notFound();
  return (
    <div className="max-w-xl space-y-6">
      <h1 className="text-2xl font-bold">Редактировать товар</h1>
      <ProductForm product={product} />
    </div>
  );
}
