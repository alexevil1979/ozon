import Link from 'next/link';
import { notFound } from 'next/navigation';
import { Button, Card, CardContent, CardHeader } from '@marketplace/ui';
import { getProduct } from '@/lib/api';
import { addToCart } from '@/lib/actions/cart';

interface ProductPageProps {
  params: Promise<{ id: string }>;
}

export default async function ProductPage({ params }: ProductPageProps) {
  const { id } = await params;
  let product: Awaited<ReturnType<typeof getProduct>> = null;
  try {
    product = await getProduct(id);
  } catch {
    notFound();
  }
  if (!product) notFound();

  return (
    <div className="container py-8">
      <nav className="text-sm text-muted-foreground mb-6">
        <Link href="/catalog" className="hover:text-foreground">
          Каталог
        </Link>
        <span className="mx-2">/</span>
        <span className="text-foreground">{product.name}</span>
      </nav>
      <div className="grid gap-8 lg:grid-cols-2">
        <div className="aspect-square rounded-lg bg-muted shrink-0 max-w-md" />
        <div>
          <Card>
            <CardHeader>
              <h1 className="text-2xl font-bold">{product.name}</h1>
              {product.brand ? (
                <p className="text-muted-foreground">{product.brand}</p>
              ) : null}
            </CardHeader>
            <CardContent className="space-y-4">
              <p className="text-2xl font-semibold">
                {new Intl.NumberFormat('ru-RU', {
                  style: 'currency',
                  currency: 'RUB',
                  maximumFractionDigits: 0,
                }).format(product.price)}
              </p>
              {product.description ? (
                <p className="text-muted-foreground whitespace-pre-wrap">
                  {product.description}
                </p>
              ) : null}
              {product.attributes && product.attributes.length > 0 ? (
                <dl className="grid gap-2 text-sm">
                  {product.attributes.map((a) => (
                    <div key={a.attribute_id} className="flex gap-2">
                      <dt className="text-muted-foreground">Атрибут:</dt>
                      <dd>{a.value}</dd>
                    </div>
                  ))}
                </dl>
              ) : null}
              <div className="pt-4">
                <AddToCartButton productId={product.id} name={product.name} price={product.price} />
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}

function AddToCartButton({
  productId,
  name,
  price,
}: {
  productId: string;
  name: string;
  price: number;
}) {
  return (
    <form action={addToCart}>
      <input type="hidden" name="product_id" value={productId} />
      <input type="hidden" name="name" value={name} />
      <input type="hidden" name="price" value={String(price)} />
      <Button type="submit" size="lg">
        В корзину
      </Button>
    </form>
  );
}
