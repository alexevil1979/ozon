import type { Metadata } from 'next';
import { redirect } from 'next/navigation';
import { cookies } from 'next/headers';
import Link from 'next/link';
import { Button } from '@marketplace/ui';

export const metadata: Metadata = {
  title: { default: 'Кабинет продавца', template: '%s | Продавец' },
  description: 'Управление товарами',
};

export default async function SellerLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const cookieStore = await cookies();
  const token = cookieStore.get('marketplace_token')?.value;
  if (!token) {
    redirect('/login?next=/seller');
  }
  return (
    <div className="min-h-screen flex flex-col">
      <header className="border-b bg-background">
        <div className="container flex h-14 items-center justify-between">
          <Link href="/seller" className="font-semibold">
            Кабинет продавца
          </Link>
          <nav className="flex items-center gap-4">
            <Link href="/seller/products">
              <Button variant="ghost" size="sm">
                Мои товары
              </Button>
            </Link>
            <Link href="/seller/products/new">
              <Button variant="ghost" size="sm">
                Добавить товар
              </Button>
            </Link>
            <Link href="/">
              <Button variant="ghost" size="sm">
                На главную
              </Button>
            </Link>
          </nav>
        </div>
      </header>
      <main className="flex-1 container py-8">{children}</main>
    </div>
  );
}
