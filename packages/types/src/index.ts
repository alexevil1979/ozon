/** Shared domain types. Extended per phase. */

export type UserId = string;
export type ProductId = string;
export type OrderId = string;

export interface Pagination {
  page: number;
  limit: number;
}

export interface PageResult<T> {
  items: T[];
  total: number;
  page: number;
  limit: number;
}

/** Search API (OpenSearch) response. */
export interface SearchHit {
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
}

export interface BrandBucket {
  key: string;
  count: number;
}

export interface PriceBucket {
  from?: number;
  to?: number;
  count: number;
}

export interface SearchFacets {
  brands: BrandBucket[];
  prices: PriceBucket[];
}

export interface SearchResult {
  hits: SearchHit[];
  total: number;
  facets: SearchFacets;
}

/** Catalog API product (detail). */
export interface CatalogProduct {
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
}

/** Cart item (client stub). */
export interface CartItem {
  productId: string;
  name: string;
  price: number;
  quantity: number;
  imageUrl?: string;
}
