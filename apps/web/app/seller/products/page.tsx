import Link from 'next/link';
import { Button } from '@marketplace/ui';
import { SellerProductList } from '@/components/seller/SellerProductList';

export default function SellerProductsPage() {
  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">Мои товары</h1>
        <Button asChild>
          <Link href="/seller/products/new">Добавить товар</Link>
        </Button>
      </div>
      <SellerProductList />
    </div>
  );
}
