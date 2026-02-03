import { redirect } from 'next/navigation';
import Link from 'next/link';
import { getCartItems } from '@/lib/actions/cart';
import { CheckoutForm } from '@/components/buyer/CheckoutForm';
import { Button } from '@marketplace/ui';

export default async function CheckoutPage() {
  const items = await getCartItems();
  if (items.length === 0) {
    redirect('/cart');
  }
  const total = items.reduce((s, i) => s + i.price * i.quantity, 0);
  return (
    <div className="container py-8 max-w-lg">
      <h1 className="text-2xl font-bold mb-6">Оформление заказа</h1>
      <div className="space-y-4 mb-6">
        {items.map((i) => (
          <div key={i.productId} className="flex justify-between text-sm">
            <span>{i.name} × {i.quantity}</span>
            <span>
              {new Intl.NumberFormat('ru-RU', { style: 'currency', currency: 'RUB', maximumFractionDigits: 0 }).format(i.price * i.quantity)}
            </span>
          </div>
        ))}
      </div>
      <p className="text-lg font-semibold mb-6">
        Итого: {new Intl.NumberFormat('ru-RU', { style: 'currency', currency: 'RUB', maximumFractionDigits: 0 }).format(total)}
      </p>
      <p className="text-sm text-muted-foreground mb-4">
        Для оформления заказа необходимо войти в аккаунт.
      </p>
      <CheckoutForm />
      <div className="mt-6">
        <Button variant="outline" asChild>
          <Link href="/cart">← Вернуться в корзину</Link>
        </Button>
      </div>
    </div>
  );
}
