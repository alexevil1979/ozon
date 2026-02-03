import Link from 'next/link';
import { Card, CardContent, CardFooter } from '@marketplace/ui';

interface ProductCardProps {
  id: string;
  name: string;
  price: number;
  brand?: string;
  slug?: string;
}

export function ProductCard({ id, name, price, brand, slug }: ProductCardProps) {
  return (
    <Link href={`/product/${id}`}>
      <Card className="h-full transition-shadow hover:shadow-md">
        <CardContent className="pt-6">
          <div className="aspect-square rounded-md bg-muted mb-4" />
          <h3 className="font-semibold line-clamp-2">{name}</h3>
          {brand ? (
            <p className="text-sm text-muted-foreground mt-1">{brand}</p>
          ) : null}
        </CardContent>
        <CardFooter className="flex justify-between items-center">
          <span className="text-lg font-semibold">
            {new Intl.NumberFormat('ru-RU', {
              style: 'currency',
              currency: 'RUB',
              maximumFractionDigits: 0,
            }).format(price)}
          </span>
        </CardFooter>
      </Card>
    </Link>
  );
}
