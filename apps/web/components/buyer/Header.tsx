import Link from 'next/link';
import { Button } from '@marketplace/ui';

export function Header() {
  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur">
      <div className="container flex h-14 items-center justify-between">
        <Link href="/" className="font-semibold">
          Marketplace
        </Link>
        <nav className="flex items-center gap-4">
          <Link href="/catalog">
            <Button variant="ghost" size="sm">
              Каталог
            </Button>
          </Link>
          <Link href="/cart">
            <Button variant="ghost" size="sm">
              Корзина
            </Button>
          </Link>
          <Link href="/orders">
            <Button variant="ghost" size="sm">
              Заказы
            </Button>
          </Link>
          <Link href="/seller">
            <Button variant="ghost" size="sm">
              Продавец
            </Button>
          </Link>
          <Link href="/login">
            <Button variant="ghost" size="sm">
              Войти
            </Button>
          </Link>
          <Link href="/admin">
            <Button variant="ghost" size="sm">
              Админ
            </Button>
          </Link>
        </nav>
      </div>
    </header>
  );
}
