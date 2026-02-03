'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { Button } from '@marketplace/ui';
import { Badge } from '@marketplace/ui';

interface ProductRow {
  id: string;
  name: string;
  slug: string;
  price: number;
  status: string;
  updated_at: string;
}

export function SellerProductList() {
  const [items, setItems] = useState<ProductRow[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    fetch('/api/seller/products?limit=50')
      .then((res) => {
        if (res.status === 401) {
          setError('Войдите в аккаунт');
          return [];
        }
        if (!res.ok) return [];
        return res.json();
      })
      .then((data) => {
        setItems(Array.isArray(data) ? data : []);
      })
      .catch(() => setError('Ошибка загрузки'))
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <p className="text-muted-foreground">Загрузка…</p>;
  if (error) return <p className="text-destructive">{error}</p>;
  if (items.length === 0) {
    return (
      <p className="text-muted-foreground">
        Нет товаров. <Link href="/seller/products/new" className="text-primary hover:underline">Добавить товар</Link>.
      </p>
    );
  }

  return (
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
          {items.map((p) => (
            <tr key={p.id} className="border-b last:border-0">
              <td className="p-4">
                <Link href={`/seller/products/${p.id}/edit`} className="font-medium hover:underline">
                  {p.name}
                </Link>
              </td>
              <td className="p-4">
                {new Intl.NumberFormat('ru-RU', { style: 'currency', currency: 'RUB', maximumFractionDigits: 0 }).format(p.price)}
              </td>
              <td className="p-4">
                <Badge variant={p.status === 'active' ? 'default' : 'secondary'}>{p.status}</Badge>
              </td>
              <td className="p-4 text-right">
                <Button variant="ghost" size="sm" asChild>
                  <Link href={`/seller/products/${p.id}/edit`}>Редактировать</Link>
                </Button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
