/**
 * Server-side API helpers. Use in RSC only.
 * Env: API_SEARCH_URL, API_CATALOG_URL (defaults for local dev).
 */

const getSearchBase = () =>
  process.env.API_SEARCH_URL ?? 'http://localhost:8083';
const getCatalogBase = () =>
  process.env.API_CATALOG_URL ?? 'http://localhost:8082';

export interface SearchParams {
  q?: string;
  category_id?: string;
  brand?: string;
  min_price?: number;
  max_price?: number;
  limit?: number;
  offset?: number;
}

export async function searchProducts(
  params: SearchParams = {}
): Promise<{
  hits: Array<{
    id: string;
    category_id: string;
    name: string;
    slug: string;
    description: string;
    price: number;
    brand: string;
    seller_id: string;
    status: string;
    created_at: string;
    updated_at: string;
    attributes?: Record<string, string>;
  }>;
  total: number;
  facets: {
    brands: { key: string; count: number }[];
    prices: { from?: number; to?: number; count: number }[];
  };
}> {
  const searchParams = new URLSearchParams();
  if (params.q) searchParams.set('q', params.q);
  if (params.category_id) searchParams.set('category_id', params.category_id);
  if (params.brand) searchParams.set('brand', params.brand);
  if (params.min_price != null) searchParams.set('min_price', String(params.min_price));
  if (params.max_price != null) searchParams.set('max_price', String(params.max_price));
  searchParams.set('limit', String(params.limit ?? 20));
  searchParams.set('offset', String(params.offset ?? 0));
  const res = await fetch(`${getSearchBase()}/search?${searchParams.toString()}`, {
    next: { revalidate: 30 },
  });
  if (!res.ok) throw new Error('Search failed');
  return res.json();
}

export async function getProduct(id: string): Promise<{
  id: string;
  category_id: string;
  name: string;
  slug: string;
  description: string;
  price: number;
  brand: string;
  seller_id: string;
  status: string;
  created_at: string;
  updated_at: string;
  attributes?: { product_id: string; attribute_id: string; value: string }[];
} | null> {
  const res = await fetch(`${getCatalogBase()}/products/${id}`, {
    next: { revalidate: 60 },
  });
  if (res.status === 404) return null;
  if (!res.ok) throw new Error('Product fetch failed');
  return res.json();
}

export async function getCategories(parentId?: string): Promise<
  Array<{
    id: string;
    name: string;
    slug: string;
    parent_id: string | null;
    created_at: string;
    updated_at: string;
  }>
> {
  const searchParams = new URLSearchParams();
  searchParams.set('limit', '100');
  if (parentId) searchParams.set('parent_id', parentId);
  const res = await fetch(`${getCatalogBase()}/categories?${searchParams.toString()}`, {
    next: { revalidate: 300 },
  });
  if (!res.ok) return [];
  return res.json();
}
