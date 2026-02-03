import { cookies } from 'next/headers';

const ADMIN_API = process.env.ADMIN_API_URL ?? 'http://localhost:8086';
const ADMIN_COOKIE = 'admin_secret';

async function getAdminSecret(): Promise<string | null> {
  const cookieStore = await cookies();
  const secret = cookieStore.get('admin_secret')?.value ?? null;
  const expected = process.env.ADMIN_SECRET ?? 'dev_admin_secret_change_in_production';
  return secret === expected ? secret : null;
}

export async function GET(request: Request) {
  const secret = await getAdminSecret();
  if (!secret) {
    return new Response(JSON.stringify({ error: 'unauthorized' }), {
      status: 401,
      headers: { 'Content-Type': 'application/json' },
    });
  }
  const { searchParams } = new URL(request.url);
  const qs = searchParams.toString();
  const url = qs ? `${ADMIN_API}/admin/products?${qs}` : `${ADMIN_API}/admin/products`;
  const res = await fetch(url, {
    headers: { 'X-Admin-Key': secret },
    cache: 'no-store',
  });
  const data = await res.json().catch(() => ({}));
  return new Response(JSON.stringify(data), {
    status: res.status,
    headers: { 'Content-Type': 'application/json' },
  });
}
