'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Button } from '@marketplace/ui';
import { placeOrder } from '@/lib/actions/order';

export function CheckoutForm() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      const result = await placeOrder();
      if ('error' in result) {
        setError(result.error);
        return;
      }
      router.push(`/orders/${result.orderId}`);
      router.refresh();
    } catch {
      setError('Ошибка оформления заказа');
    } finally {
      setLoading(false);
    }
  }

  return (
    <form onSubmit={handleSubmit}>
      {error && <p className="text-sm text-destructive mb-4">{error}</p>}
      <Button type="submit" size="lg" disabled={loading}>
        {loading ? 'Оформление…' : 'Оформить заказ'}
      </Button>
    </form>
  );
}
