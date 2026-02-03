import Link from 'next/link';
import { redirect } from 'next/navigation';
import { cookies } from 'next/headers';
import { Button } from '@marketplace/ui';
import { OrderList } from '@/components/buyer/OrderList';

export default async function OrdersPage() {
  const cookieStore = await cookies();
  const token = cookieStore.get('marketplace_token')?.value;
  if (!token) {
    redirect('/login?next=/orders');
  }
  return (
    <div className="container py-8">
      <h1 className="text-2xl font-bold mb-6">Мои заказы</h1>
      <OrderList />
      <div className="mt-6">
        <Button variant="outline" asChild>
          <Link href="/catalog">В каталог</Link>
        </Button>
      </div>
    </div>
  );
}
