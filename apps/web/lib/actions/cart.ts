'use server';

import { cookies } from 'next/headers';

export interface CartItemCookie {
  productId: string;
  name: string;
  price: number;
  quantity: number;
}

const CART_COOKIE = 'marketplace_cart';
const MAX_AGE = 7 * 24 * 60 * 60; // 7 days

async function getCart(): Promise<CartItemCookie[]> {
  const cookieStore = await cookies();
  const raw = cookieStore.get(CART_COOKIE)?.value;
  if (!raw) return [];
  try {
    const parsed = JSON.parse(raw) as CartItemCookie[];
    return Array.isArray(parsed) ? parsed : [];
  } catch {
    return [];
  }
}

async function setCart(items: CartItemCookie[]) {
  const cookieStore = await cookies();
  cookieStore.set(CART_COOKIE, JSON.stringify(items), {
    path: '/',
    maxAge: MAX_AGE,
    sameSite: 'lax',
  });
}

export async function addToCart(productId: string, name: string, price: number, quantity?: number): Promise<void>;
export async function addToCart(formData: FormData): Promise<void>;
export async function addToCart(
  productIdOrFormData: string | FormData,
  name?: string,
  price?: number,
  quantity = 1
): Promise<void> {
  let productId: string;
  let itemName: string;
  let itemPrice: number;
  if (productIdOrFormData instanceof FormData) {
    productId = (productIdOrFormData.get('product_id') as string) ?? '';
    itemName = (productIdOrFormData.get('name') as string) ?? '';
    itemPrice = parseFloat((productIdOrFormData.get('price') as string) ?? '0');
  } else {
    productId = productIdOrFormData;
    itemName = name ?? '';
    itemPrice = price ?? 0;
  }
  if (!productId || !itemName) return;
  const cart = await getCart();
  const existing = cart.find((i) => i.productId === productId);
  if (existing) {
    existing.quantity += quantity;
  } else {
    cart.push({ productId, name: itemName, price: itemPrice, quantity });
  }
  await setCart(cart);
}

export async function removeFromCart(productId: string): Promise<void>;
export async function removeFromCart(formData: FormData): Promise<void>;
export async function removeFromCart(productIdOrFormData: string | FormData): Promise<void> {
  const productId =
    typeof productIdOrFormData === 'string'
      ? productIdOrFormData
      : (productIdOrFormData.get('product_id') as string | null);
  if (!productId) return;
  const cart = await getCart().then((items) => items.filter((i) => i.productId !== productId));
  await setCart(cart);
}

export async function updateQuantity(productId: string, quantity: number) {
  const cart = await getCart();
  const item = cart.find((i) => i.productId === productId);
  if (!item) return;
  if (quantity <= 0) {
    await setCart(cart.filter((i) => i.productId !== productId));
    return;
  }
  item.quantity = quantity;
  await setCart(cart);
}

export async function getCartItems(): Promise<CartItemCookie[]> {
  return getCart();
}

export async function clearCart(): Promise<void> {
  const cookieStore = await cookies();
  cookieStore.set(CART_COOKIE, '[]', { path: '/', maxAge: 0, sameSite: 'lax' });
}
