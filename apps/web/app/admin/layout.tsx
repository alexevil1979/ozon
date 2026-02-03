import type { Metadata } from 'next';
import { redirect } from 'next/navigation';
import { cookies } from 'next/headers';
import Link from 'next/link';
import { Button } from '@marketplace/ui';

export const metadata: Metadata = {
  title: { default: 'Админ-панель', template: '%s | Админ' },
  description: 'Управление пользователями и модерация товаров',
};

const ADMIN_COOKIE = 'admin_secret';

export default async function AdminLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const cookieStore = await cookies();
  const secret = cookieStore.get(ADMIN_COOKIE)?.value;
  const expected = process.env.ADMIN_SECRET ?? 'dev_admin_secret_change_in_production';
  if (secret !== expected) {
    redirect('/admin/login');
  }
  return (
    <div className="min-h-screen flex flex-col">
      <header className="border-b bg-muted/30">
        <div className="container flex h-14 items-center justify-between">
          <Link href="/admin" className="font-semibold">
            Админ-панель
          </Link>
          <nav className="flex items-center gap-4">
            <Link href="/admin/users">
              <Button variant="ghost" size="sm">
                Пользователи
              </Button>
            </Link>
            <Link href="/admin/products">
              <Button variant="ghost" size="sm">
                Товары (модерация)
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
