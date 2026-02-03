'use client';

import Link from 'next/link';
import { removeFromCart } from '@/lib/actions/cart';
import { Button } from '@marketplace/ui';
import type { CartItemCookie } from '@/lib/actions/cart';

interface CartItemRowProps {
  item: CartItemCookie;
}

export function CartItemRow({ item }: CartItemRowProps) {
  return (
    <div className="flex items-center justify-between gap-4 rounded-lg border p-4">
      <div className="flex-1 min-w-0">
        <Link href={`/product/${item.productId}`} className="font-medium hover:underline truncate block">
          {item.name}
        </Link>
        <p className="text-sm text-muted-foreground">
          {new Intl.NumberFormat('ru-RU', {
            style: 'currency',
            currency: 'RUB',
            maximumFractionDigits: 0,
          }).format(item.price)}{' '}
          × {item.quantity}
        </p>
      </div>
      <p className="font-semibold shrink-0">
        {new Intl.NumberFormat('ru-RU', {
          style: 'currency',
          currency: 'RUB',
          maximumFractionDigits: 0,
        }).format(item.price * item.quantity)}
      </p>
      <form action={removeFromCart}>
        <input type="hidden" name="product_id" value={item.productId} />
        <Button type="submit" variant="ghost" size="sm">
          Удалить
        </Button>
      </form>
    </div>
  );
}
