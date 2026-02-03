import Link from 'next/link';
import { Button } from '@marketplace/ui';
import { getCartItems } from '@/lib/actions/cart';
import { CartItemRow } from '@/components/buyer/CartItemRow';

export default async function CartPage() {
  const items = await getCartItems();
  const total = items.reduce((s, i) => s + i.price * i.quantity, 0);

  return (
    <div className="container py-8">
      <h1 className="text-2xl font-bold mb-6">Корзина</h1>
      {items.length === 0 ? (
        <p className="text-muted-foreground mb-6">Корзина пуста.</p>
      ) : (
        <>
          <div className="space-y-4 mb-8">
            {items.map((item) => (
              <CartItemRow key={item.productId} item={item} />
            ))}
          </div>
          <div className="flex justify-between items-center border-t pt-6">
            <p className="text-lg font-semibold">
              Итого:{' '}
              {new Intl.NumberFormat('ru-RU', {
                style: 'currency',
                currency: 'RUB',
                maximumFractionDigits: 0,
              }).format(total)}
            </p>
            <Button asChild>
              <Link href="/checkout">Оформить заказ</Link>
            </Button>
          </div>
        </>
      )}
      <div className="mt-8">
        <Link href="/catalog" className="text-primary hover:underline">
          ← Вернуться в каталог
        </Link>
      </div>
    </div>
  );
}
