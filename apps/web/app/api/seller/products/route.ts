import { cookies } from 'next/headers';

const SELLER_API = process.env.SELLER_API_URL ?? 'http://localhost:8084';

async function getToken(): Promise<string | null> {
  const cookieStore = await cookies();
  return cookieStore.get('marketplace_token')?.value ?? null;
}

export async function GET(request: Request) {
  const token = await getToken();
  if (!token) {
    return new Response(JSON.stringify({ error: 'unauthorized' }), {
      status: 401,
      headers: { 'Content-Type': 'application/json' },
    });
  }
  const { searchParams } = new URL(request.url);
  const qs = searchParams.toString();
  const url = qs ? `${SELLER_API}/seller/products?${qs}` : `${SELLER_API}/seller/products`;
  const res = await fetch(url, {
    headers: { Authorization: `Bearer ${token}` },
    cache: 'no-store',
  });
  const data = await res.json().catch(() => ({}));
  return new Response(JSON.stringify(data), {
    status: res.status,
    headers: { 'Content-Type': 'application/json' },
  });
}

export async function POST(request: Request) {
  const token = await getToken();
  if (!token) {
    return new Response(JSON.stringify({ error: 'unauthorized' }), {
      status: 401,
      headers: { 'Content-Type': 'application/json' },
    });
  }
  const body = await request.json().catch(() => ({}));
  const res = await fetch(`${SELLER_API}/seller/products`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(body),
  });
  const data = await res.json().catch(() => ({}));
  return new Response(JSON.stringify(data), {
    status: res.status,
    headers: { 'Content-Type': 'application/json' },
  });
}
