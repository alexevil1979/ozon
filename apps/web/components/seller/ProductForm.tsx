'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { Button, Input } from '@marketplace/ui';

interface ProductFormProps {
  product?: {
    id: string;
    category_id: string;
    name: string;
    slug: string;
    description?: string;
    price: number;
    brand?: string;
    status: string;
  };
}

export function ProductForm({ product }: ProductFormProps) {
  const router = useRouter();
  const [categories, setCategories] = useState<{ id: string; name: string }[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [form, setForm] = useState({
    category_id: product?.category_id ?? '',
    name: product?.name ?? '',
    slug: product?.slug ?? '',
    description: product?.description ?? '',
    price: product?.price ?? 0,
    brand: product?.brand ?? '',
    status: product?.status ?? 'draft',
  });

  useEffect(() => {
    fetch((process.env.NEXT_PUBLIC_API_CATALOG_URL ?? 'http://localhost:8082') + '/categories?limit=100')
      .then((r) => r.json())
      .then((data) => setCategories(Array.isArray(data) ? data : []))
      .catch(() => {});
  }, []);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      const url = product ? `/api/seller/products/${product.id}` : '/api/seller/products';
      const method = product ? 'PUT' : 'POST';
      const body = {
        category_id: form.category_id,
        name: form.name,
        slug: form.slug,
        description: form.description || undefined,
        price: Number(form.price),
        brand: form.brand || undefined,
        status: form.status,
      };
      const res = await fetch(url, {
        method,
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      });
      if (!res.ok) {
        const data = await res.json().catch(() => ({}));
        setError(data.error ?? 'Ошибка сохранения');
        return;
      }
      router.push('/seller/products');
      router.refresh();
    } catch {
      setError('Ошибка сети');
    } finally {
      setLoading(false);
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div>
        <label className="text-sm font-medium">Категория *</label>
        <select
          required
          value={form.category_id}
          onChange={(e) => setForm((f) => ({ ...f, category_id: e.target.value }))}
          className="mt-1 flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
        >
          <option value="">Выберите</option>
          {categories.map((c) => (
            <option key={c.id} value={c.id}>{c.name}</option>
          ))}
        </select>
      </div>
      <div>
        <label className="text-sm font-medium">Название *</label>
        <Input
          required
          value={form.name}
          onChange={(e) => setForm((f) => ({ ...f, name: e.target.value }))}
          className="mt-1"
        />
      </div>
      <div>
        <label className="text-sm font-medium">Slug *</label>
        <Input
          required
          value={form.slug}
          onChange={(e) => setForm((f) => ({ ...f, slug: e.target.value }))}
          placeholder="product-slug"
          className="mt-1"
        />
      </div>
      <div>
        <label className="text-sm font-medium">Описание</label>
        <textarea
          value={form.description}
          onChange={(e) => setForm((f) => ({ ...f, description: e.target.value }))}
          rows={3}
          className="mt-1 flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
        />
      </div>
      <div>
        <label className="text-sm font-medium">Цена (₽) *</label>
        <Input
          type="number"
          required
          min={0}
          step={0.01}
          value={form.price}
          onChange={(e) => setForm((f) => ({ ...f, price: parseFloat(e.target.value) || 0 }))}
          className="mt-1"
        />
      </div>
      <div>
        <label className="text-sm font-medium">Бренд</label>
        <Input
          value={form.brand}
          onChange={(e) => setForm((f) => ({ ...f, brand: e.target.value }))}
          className="mt-1"
        />
      </div>
      <div>
        <label className="text-sm font-medium">Статус</label>
        <select
          value={form.status}
          onChange={(e) => setForm((f) => ({ ...f, status: e.target.value }))}
          className="mt-1 flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
        >
          <option value="draft">Черновик</option>
          <option value="active">Активный</option>
        </select>
      </div>
      {error && <p className="text-sm text-destructive">{error}</p>}
      <div className="flex gap-2">
        <Button type="submit" disabled={loading}>{loading ? 'Сохранение…' : 'Сохранить'}</Button>
        <Button type="button" variant="outline" asChild>
          <Link href="/seller/products">Отмена</Link>
        </Button>
      </div>
    </form>
  );
}
