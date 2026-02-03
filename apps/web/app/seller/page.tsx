import Link from 'next/link';
import { Button } from '@marketplace/ui';

export default function SellerDashboardPage() {
  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Кабинет продавца</h1>
      <p className="text-muted-foreground">
        Управляйте своими товарами: добавляйте новые, редактируйте и отслеживайте статус модерации.
      </p>
      <div className="flex gap-4">
        <Button asChild>
          <Link href="/seller/products">Мои товары</Link>
        </Button>
        <Button variant="outline" asChild>
          <Link href="/seller/products/new">Добавить товар</Link>
        </Button>
      </div>
    </div>
  );
}
