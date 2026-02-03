'use client';

import { useEffect, useState } from 'react';
import { Button, Badge } from '@marketplace/ui';

interface ProductRow {
  id: string;
  name: string;
  slug: string;
  price: number;
  brand?: string;
  seller_id: string;
  status: string;
  created_at: string;
  updated_at: string;
}

export function AdminProductList() {
  const [products, setProducts] = useState<ProductRow[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [filter, setFilter] = useState<string>('');

  useEffect(() => {
    const qs = filter ? `?status=${filter}` : '';
    fetch(`/api/admin/products${qs}`)
      .then((res) => {
        if (res.status === 401) {
          setError('Доступ запрещён');
          return [];
        }
        if (!res.ok) return [];
        return res.json();
      })
      .then((data) => {
        setProducts(Array.isArray(data) ? data : []);
      })
      .catch(() => setError('Ошибка загрузки'))
      .finally(() => setLoading(false));
  }, [filter]);

  async function setStatus(id: string, status: string) {
    const res = await fetch(`/api/admin/products/${id}/status`, {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ status }),
    });
    if (!res.ok) return;
    setProducts((prev) =>
      prev.map((p) => (p.id === id ? { ...p, status } : p))
    );
  }

  if (loading) return <p className="text-muted-foreground">Загрузка…</p>;
  if (error) return <p className="text-destructive">{error}</p>;

  return (
    <div className="space-y-4">
      <div className="flex gap-2">
        <button
          type="button"
          onClick={() => setFilter('')}
          className={`rounded-md px-3 py-1.5 text-sm ${filter === '' ? 'bg-muted font-medium' : 'hover:bg-muted'}`}
        >
          Все
        </button>
        <button
          type="button"
          onClick={() => setFilter('draft')}
          className={`rounded-md px-3 py-1.5 text-sm ${filter === 'draft' ? 'bg-muted font-medium' : 'hover:bg-muted'}`}
        >
          Черновики
        </button>
        <button
          type="button"
          onClick={() => setFilter('active')}
          className={`rounded-md px-3 py-1.5 text-sm ${filter === 'active' ? 'bg-muted font-medium' : 'hover:bg-muted'}`}
        >
          Активные
        </button>
      </div>
      {products.length === 0 ? (
        <p className="text-muted-foreground">Нет товаров.</p>
      ) : (
        <div className="rounded-lg border">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b bg-muted/50">
                <th className="text-left p-4 font-medium">Название</th>
                <th className="text-left p-4 font-medium">Цена</th>
                <th className="text-left p-4 font-medium">Статус</th>
                <th className="text-right p-4 font-medium">Действия</th>
              </tr>
            </thead>
            <tbody>
              {products.map((p) => (
                <tr key={p.id} className="border-b last:border-0">
                  <td className="p-4">
                    <span className="font-medium">{p.name}</span>
                    {p.brand ? (
                      <span className="text-muted-foreground ml-2">({p.brand})</span>
                    ) : null}
                  </td>
                  <td className="p-4">
                    {new Intl.NumberFormat('ru-RU', {
                      style: 'currency',
                      currency: 'RUB',
                      maximumFractionDigits: 0,
                    }).format(p.price)}
                  </td>
                  <td className="p-4">
                    <Badge variant={p.status === 'active' ? 'default' : 'secondary'}>
                      {p.status}
                    </Badge>
                  </td>
                  <td className="p-4 text-right">
                    {p.status === 'draft' && (
                      <Button
                        size="sm"
                        variant="outline"
                        onClick={() => setStatus(p.id, 'active')}
                      >
                        Одобрить
                      </Button>
                    )}
                    {p.status === 'active' && (
                      <Button
                        size="sm"
                        variant="ghost"
                        onClick={() => setStatus(p.id, 'draft')}
                      >
                        В черновик
                      </Button>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
