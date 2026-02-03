import { notFound, redirect } from 'next/navigation';
import { cookies } from 'next/headers';
import Link from 'next/link';
import { Button } from '@marketplace/ui';
import { OrderDetail } from '@/components/buyer/OrderDetail';

const ORDER_API = process.env.ORDER_API_URL ?? 'http://localhost:8085';

async function getOrder(id: string, token: string) {
  const res = await fetch(`${ORDER_API}/orders/${id}`, {
    headers: { Authorization: `Bearer ${token}` },
    cache: 'no-store',
  });
  if (res.status === 404) return null;
  if (!res.ok) return null;
  return res.json();
}

export default async function OrderDetailPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;
  const cookieStore = await cookies();
  const token = cookieStore.get('marketplace_token')?.value;
  if (!token) {
    redirect(`/login?next=/orders/${id}`);
  }
  const order = await getOrder(id, token);
  if (!order) notFound();
  return (
    <div className="container py-8 max-w-2xl">
      <h1 className="text-2xl font-bold mb-6">Заказ #{id.slice(0, 8)}</h1>
      <OrderDetail order={order} />
      <div className="mt-6">
        <Button variant="outline" asChild>
          <Link href="/orders">← К списку заказов</Link>
        </Button>
      </div>
    </div>
  );
}
