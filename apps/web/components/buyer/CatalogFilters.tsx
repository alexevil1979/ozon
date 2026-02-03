'use client';

import { useRouter, useSearchParams } from 'next/navigation';
import { useCallback } from 'react';
import { Button, Input } from '@marketplace/ui';

interface CatalogFiltersProps {
  facets?: {
    brands: { key: string; count: number }[];
    prices: { from?: number; to?: number; count: number }[];
  };
  currentQ?: string;
  currentCategoryId?: string;
  currentBrand?: string;
  currentMinPrice?: string;
  currentMaxPrice?: string;
}

export function CatalogFilters({
  facets,
  currentQ,
  currentBrand,
  currentMinPrice,
  currentMaxPrice,
}: CatalogFiltersProps) {
  const router = useRouter();
  const searchParams = useSearchParams();

  const updateParams = useCallback(
    (updates: Record<string, string | undefined>) => {
      const sp = new URLSearchParams(searchParams.toString());
      Object.entries(updates).forEach(([k, v]) => {
        if (v === undefined || v === '') sp.delete(k);
        else sp.set(k, v);
      });
      sp.delete('page');
      router.push(`/catalog?${sp.toString()}`);
    },
    [router, searchParams]
  );

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const form = e.currentTarget;
    const q = (form.querySelector('[name="q"]') as HTMLInputElement)?.value;
    const min = (form.querySelector('[name="min_price"]') as HTMLInputElement)?.value;
    const max = (form.querySelector('[name="max_price"]') as HTMLInputElement)?.value;
    updateParams({ q: q || undefined, min_price: min || undefined, max_price: max || undefined });
  };

  return (
    <div className="space-y-6">
      <form onSubmit={handleSubmit} className="space-y-4">
        <label className="text-sm font-medium">Поиск</label>
        <Input
          name="q"
          defaultValue={currentQ}
          placeholder="Название или бренд"
          className="w-full"
        />
        <label className="text-sm font-medium block">Цена</label>
        <div className="flex gap-2">
          <Input
            name="min_price"
            type="number"
            min={0}
            step={100}
            defaultValue={currentMinPrice}
            placeholder="От"
            className="w-full"
          />
          <Input
            name="max_price"
            type="number"
            min={0}
            step={100}
            defaultValue={currentMaxPrice}
            placeholder="До"
            className="w-full"
          />
        </div>
        <Button type="submit" size="sm" className="w-full">
          Применить
        </Button>
      </form>
      {facets?.brands && facets.brands.length > 0 ? (
        <div>
          <label className="text-sm font-medium block mb-2">Бренд</label>
          <ul className="space-y-1">
            <li>
              <button
                type="button"
                onClick={() => updateParams({ brand: undefined })}
                className={`text-sm w-full text-left px-2 py-1 rounded ${!currentBrand ? 'bg-muted font-medium' : 'hover:bg-muted'}`}
              >
                Все
              </button>
            </li>
            {facets.brands.map((b) => (
              <li key={b.key}>
                <button
                  type="button"
                  onClick={() => updateParams({ brand: b.key === '—' ? undefined : b.key })}
                  className={`text-sm w-full text-left px-2 py-1 rounded truncate ${currentBrand === b.key ? 'bg-muted font-medium' : 'hover:bg-muted'}`}
                >
                  {b.key} ({b.count})
                </button>
              </li>
            ))}
          </ul>
        </div>
      ) : null}
    </div>
  );
}
