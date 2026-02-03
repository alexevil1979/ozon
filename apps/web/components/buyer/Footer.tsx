import Link from 'next/link';

export function Footer() {
  return (
    <footer className="border-t bg-muted/30 mt-auto">
      <div className="container py-8">
        <div className="flex flex-col gap-4 sm:flex-row sm:justify-between">
          <Link href="/" className="font-semibold text-muted-foreground hover:text-foreground">
            Marketplace
          </Link>
          <div className="flex gap-6 text-sm text-muted-foreground">
            <Link href="/catalog" className="hover:text-foreground">
              Каталог
            </Link>
            <Link href="/cart" className="hover:text-foreground">
              Корзина
            </Link>
          </div>
        </div>
        <p className="mt-6 text-xs text-muted-foreground">
          © Marketplace v0.1. Buyer · Seller · Admin
        </p>
      </div>
    </footer>
  );
}
