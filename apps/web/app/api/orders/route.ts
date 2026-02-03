import { cookies } from 'next/headers';

const ORDER_API = process.env.ORDER_API_URL ?? 'http://localhost:8085';

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
  const url = qs ? `${ORDER_API}/orders?${qs}` : `${ORDER_API}/orders`;
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
