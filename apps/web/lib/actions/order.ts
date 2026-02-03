'use server';

import { cookies } from 'next/headers';
import { redirect } from 'next/navigation';
import { getCartItems, clearCart } from './cart';

const ORDER_API = process.env.ORDER_API_URL ?? 'http://localhost:8085';

export type PlaceOrderResult = { orderId: string } | { error: string };

export async function placeOrder(): Promise<PlaceOrderResult> {
  const cart = await getCartItems();
  if (cart.length === 0) {
    return { error: 'Корзина пуста' };
  }
  const cookieStore = await cookies();
  const token = cookieStore.get('marketplace_token')?.value;
  if (!token) {
    return { error: 'Войдите в аккаунт для оформления заказа' };
  }
  const body = {
    items: cart.map((i) => ({
      product_id: i.productId,
      name: i.name,
      price: i.price,
      quantity: i.quantity,
    })),
  };
  const res = await fetch(`${ORDER_API}/orders`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(body),
  });
  const data = await res.json().catch(() => ({}));
  if (!res.ok) {
    return { error: (data as { error?: string }).error ?? 'Ошибка создания заказа' };
  }
  const orderId = (data as { id?: string }).id;
  if (!orderId) {
    return { error: 'Некорректный ответ сервера' };
  }
  await clearCart();
  return { orderId };
}
