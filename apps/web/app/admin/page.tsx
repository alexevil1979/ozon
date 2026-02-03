import Link from 'next/link';
import { Button } from '@marketplace/ui';

export default function AdminDashboardPage() {
  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Админ-панель</h1>
      <p className="text-muted-foreground">
        Управление пользователями и модерация товаров.
      </p>
      <div className="flex gap-4">
        <Button asChild>
          <Link href="/admin/users">Пользователи</Link>
        </Button>
        <Button variant="outline" asChild>
          <Link href="/admin/products">Товары (модерация)</Link>
        </Button>
      </div>
    </div>
  );
}
