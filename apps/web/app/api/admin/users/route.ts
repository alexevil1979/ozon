import { cookies } from 'next/headers';

const ADMIN_API = process.env.ADMIN_API_URL ?? 'http://localhost:8086';
const ADMIN_COOKIE_NAME = 'admin_secret';

export async function GET() {
  const cookieStore = await cookies();
  const secret = cookieStore.get(ADMIN_COOKIE_NAME)?.value;
  const expected = process.env.ADMIN_SECRET ?? 'dev_admin_secret_change_in_production';
  if (secret !== expected) {
    return new Response(JSON.stringify({ error: 'unauthorized' }), {
      status: 401,
      headers: { 'Content-Type': 'application/json' },
    });
  }
  const res = await fetch(`${ADMIN_API}/admin/users?limit=100`, {
    headers: { 'X-Admin-Key': secret },
    cache: 'no-store',
  });
  const data = await res.json().catch(() => ({}));
  return new Response(JSON.stringify(data), {
    status: res.status,
    headers: { 'Content-Type': 'application/json' },
  });
}
