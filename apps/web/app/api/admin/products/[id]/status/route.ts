import { cookies } from 'next/headers';

const ADMIN_API = process.env.ADMIN_API_URL ?? 'http://localhost:8086';

export async function PATCH(
  request: Request,
  { params }: { params: Promise<{ id: string }> }
) {
  const cookieStore = await cookies();
  const secret = cookieStore.get('admin_secret')?.value ?? null; // same as layout ADMIN_COOKIE
  const expected = process.env.ADMIN_SECRET ?? 'dev_admin_secret_change_in_production';
  if (secret !== expected) {
    return new Response(JSON.stringify({ error: 'unauthorized' }), {
      status: 401,
      headers: { 'Content-Type': 'application/json' },
    });
  }
  const { id } = await params;
  const body = await request.json().catch(() => ({}));
  const res = await fetch(`${ADMIN_API}/admin/products/${id}/status`, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
      'X-Admin-Key': secret,
    },
    body: JSON.stringify(body),
  });
  const data = await res.json().catch(() => ({}));
  return new Response(JSON.stringify(data), {
    status: res.status,
    headers: { 'Content-Type': 'application/json' },
  });
}
