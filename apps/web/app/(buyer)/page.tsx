import Link from 'next/link';
import { Button } from '@marketplace/ui';

export default function BuyerHomePage() {
  return (
    <section className="container flex flex-col items-center justify-center gap-8 py-24">
      <h1 className="text-4xl font-bold tracking-tight text-center">
        Marketplace
      </h1>
      <p className="text-muted-foreground text-center max-w-md">
        Покупайте и продавайте. Каталог товаров с поиском и фильтрами.
      </p>
      <Link href="/catalog">
        <Button size="lg">Перейти в каталог</Button>
      </Link>
    </section>
  );
}
