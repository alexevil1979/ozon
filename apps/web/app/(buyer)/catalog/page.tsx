import type { Metadata } from 'next';
import { Suspense } from 'react';
import { CatalogFilters } from '@/components/buyer/CatalogFilters';
import { ProductCard } from '@/components/buyer/ProductCard';
import { Skeleton } from '@marketplace/ui';
import { searchProducts } from '@/lib/api';

export const metadata: Metadata = {
  title: 'Каталог',
  description: 'Каталог товаров с поиском и фильтрами',
};

const LIMIT = 20;

interface CatalogPageProps {
  searchParams: Promise<{ q?: string; category_id?: string; brand?: string; min_price?: string; max_price?: string; page?: string }>;
}

export default async function CatalogPage({ searchParams }: CatalogPageProps) {
  const params = await searchParams;
  const page = Math.max(1, parseInt(params.page ?? '1', 10));
  const offset = (page - 1) * LIMIT;

  let data: Awaited<ReturnType<typeof searchProducts>> | null = null;
  let error: string | null = null;
  try {
    data = await searchProducts({
      q: params.q,
      category_id: params.category_id,
      brand: params.brand,
      min_price: params.min_price ? parseFloat(params.min_price) : undefined,
      max_price: params.max_price ? parseFloat(params.max_price) : undefined,
      limit: LIMIT,
      offset,
    });
  } catch (e) {
    error = 'Не удалось загрузить каталог. Проверьте, что сервис поиска запущен.';
  }

  return (
    <div className="container py-8">
      <h1 className="text-2xl font-bold mb-6">Каталог</h1>
      <div className="flex flex-col lg:flex-row gap-8">
        <aside className="lg:w-64 shrink-0">
          <Suspense fallback={<Skeleton className="h-64 w-full" />}>
            <CatalogFilters
              facets={data?.facets}
              currentQ={params.q}
              currentCategoryId={params.category_id}
              currentBrand={params.brand}
              currentMinPrice={params.min_price}
              currentMaxPrice={params.max_price}
            />
          </Suspense>
        </aside>
        <div className="flex-1">
          {error ? (
            <p className="text-destructive">{error}</p>
          ) : data ? (
            <>
              <p className="text-muted-foreground text-sm mb-4">
                Найдено: {data.total}
              </p>
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
                {data.hits.map((hit) => (
                  <ProductCard
                    key={hit.id}
                    id={hit.id}
                    name={hit.name}
                    price={hit.price}
                    brand={hit.brand || undefined}
                    slug={hit.slug}
                  />
                ))}
              </div>
              {data.hits.length === 0 ? (
                <p className="text-muted-foreground py-8">Товары не найдены.</p>
              ) : null}
              <div className="mt-8 flex justify-center gap-2">
                {page > 1 ? (
                  <a
                    href={buildCatalogUrl(params, page - 1)}
                    className="text-primary hover:underline"
                  >
                    ← Назад
                  </a>
                ) : null}
                {offset + data.hits.length < data.total ? (
                  <a
                    href={buildCatalogUrl(params, page + 1)}
                    className="text-primary hover:underline"
                  >
                    Вперёд →
                  </a>
                ) : null}
              </div>
            </>
          ) : null}
        </div>
      </div>
    </div>
  );
}

function buildCatalogUrl(
  params: Record<string, string | undefined>,
  page: number
): string {
  const sp = new URLSearchParams();
  if (params.q) sp.set('q', params.q);
  if (params.category_id) sp.set('category_id', params.category_id);
  if (params.brand) sp.set('brand', params.brand);
  if (params.min_price) sp.set('min_price', params.min_price);
  if (params.max_price) sp.set('max_price', params.max_price);
  if (page > 1) sp.set('page', String(page));
  const qs = sp.toString();
  return qs ? `/catalog?${qs}` : '/catalog';
}
