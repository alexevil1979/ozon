'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Button, Badge } from '@marketplace/ui';

interface OrderItem {
  product_id: string;
  name: string;
  price: number;
  quantity: number;
}

interface Order {
  id: string;
  status: string;
  total: number;
  payment_id?: string;
  created_at: string;
  items: OrderItem[];
}

export function OrderDetail({ order }: { order: Order }) {
  const router = useRouter();
  const [paying, setPaying] = useState(false);
  const [error, setError] = useState('');

  async function handlePay() {
    setError('');
    setPaying(true);
    try {
      const res = await fetch(`/api/orders/${order.id}/pay`, { method: 'POST' });
      const data = await res.json().catch(() => ({}));
      if (!res.ok) {
        setError((data as { error?: string }).error ?? 'Ошибка оплаты');
        return;
      }
      router.refresh();
    } catch {
      setError('Ошибка сети');
    } finally {
      setPaying(false);
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <Badge variant={order.status === 'paid' ? 'default' : 'secondary'}>{order.status}</Badge>
        <p className="text-lg font-semibold">
          {new Intl.NumberFormat('ru-RU', { style: 'currency', currency: 'RUB', maximumFractionDigits: 0 }).format(order.total)}
        </p>
      </div>
      <p className="text-sm text-muted-foreground">
        Создан: {new Date(order.created_at).toLocaleString('ru-RU')}
      </p>
      {order.items?.length ? (
        <div className="rounded-lg border">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b bg-muted/50">
                <th className="text-left p-4 font-medium">Товар</th>
                <th className="text-right p-4 font-medium">Цена</th>
                <th className="text-right p-4 font-medium">Кол-во</th>
                <th className="text-right p-4 font-medium">Сумма</th>
              </tr>
            </thead>
            <tbody>
              {order.items.map((it, idx) => (
                <tr key={idx} className="border-b last:border-0">
                  <td className="p-4">{it.name}</td>
                  <td className="p-4 text-right">{new Intl.NumberFormat('ru-RU', { style: 'currency', currency: 'RUB', maximumFractionDigits: 0 }).format(it.price)}</td>
                  <td className="p-4 text-right">{it.quantity}</td>
                  <td className="p-4 text-right">{new Intl.NumberFormat('ru-RU', { style: 'currency', currency: 'RUB', maximumFractionDigits: 0 }).format(it.price * it.quantity)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      ) : null}
      {order.status === 'pending' && (
        <div>
          {error && <p className="text-sm text-destructive mb-2">{error}</p>}
          <Button onClick={handlePay} disabled={paying}>
            {paying ? 'Оплата…' : 'Оплатить (заглушка)'}
          </Button>
        </div>
      )}
      {order.status === 'paid' && order.payment_id && (
        <p className="text-sm text-muted-foreground">Оплачено (payment_id: {order.payment_id})</p>
      )}
    </div>
  );
}
