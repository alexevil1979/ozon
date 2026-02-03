import { cookies } from 'next/headers';

const ORDER_API = process.env.ORDER_API_URL ?? 'http://localhost:8085';

async function getToken(): Promise<string | null> {
  const cookieStore = await cookies();
  return cookieStore.get('marketplace_token')?.value ?? null;
}

export async function GET(
  _request: Request,
  { params }: { params: Promise<{ id: string }> }
) {
  const token = await getToken();
  if (!token) {
    return new Response(JSON.stringify({ error: 'unauthorized' }), {
      status: 401,
      headers: { 'Content-Type': 'application/json' },
    });
  }
  const { id } = await params;
  const res = await fetch(`${ORDER_API}/orders/${id}`, {
    headers: { Authorization: `Bearer ${token}` },
    cache: 'no-store',
  });
  const data = await res.json().catch(() => ({}));
  return new Response(JSON.stringify(data), {
    status: res.status,
    headers: { 'Content-Type': 'application/json' },
  });
}
