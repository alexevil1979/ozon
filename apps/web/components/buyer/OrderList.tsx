'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { Button } from '@marketplace/ui';
import { Badge } from '@marketplace/ui';

interface OrderRow {
  id: string;
  buyer_id: string;
  status: string;
  total: number;
  created_at: string;
}

export function OrderList() {
  const [orders, setOrders] = useState<OrderRow[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    fetch('/api/orders')
      .then((res) => {
        if (res.status === 401) {
          setError('Войдите в аккаунт');
          return [];
        }
        if (!res.ok) return [];
        return res.json();
      })
      .then((data) => {
        setOrders(Array.isArray(data) ? data : []);
      })
      .catch(() => setError('Ошибка загрузки'))
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <p className="text-muted-foreground">Загрузка…</p>;
  if (error) return <p className="text-destructive">{error}</p>;
  if (orders.length === 0) {
    return (
      <p className="text-muted-foreground">
        У вас пока нет заказов. <Link href="/catalog" className="text-primary hover:underline">Перейти в каталог</Link>.
      </p>
    );
  }

  return (
    <div className="space-y-4">
      {orders.map((o) => (
        <div key={o.id} className="flex items-center justify-between rounded-lg border p-4">
          <div>
            <Link href={`/orders/${o.id}`} className="font-medium hover:underline">
              Заказ #{o.id.slice(0, 8)}
            </Link>
            <p className="text-sm text-muted-foreground">
              {new Date(o.created_at).toLocaleString('ru-RU')} ·{' '}
              {new Intl.NumberFormat('ru-RU', { style: 'currency', currency: 'RUB', maximumFractionDigits: 0 }).format(o.total)}
            </p>
          </div>
          <Badge variant={o.status === 'paid' ? 'default' : 'secondary'}>{o.status}</Badge>
          <Button variant="outline" size="sm" asChild>
            <Link href={`/orders/${o.id}`}>Подробнее</Link>
          </Button>
        </div>
      ))}
    </div>
  );
}
