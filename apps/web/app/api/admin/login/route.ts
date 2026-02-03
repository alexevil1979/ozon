import { NextResponse } from 'next/server';
import { cookies } from 'next/headers';

const ADMIN_COOKIE = 'admin_secret';
const MAX_AGE = 24 * 60 * 60; // 24h

export async function POST(request: Request) {
  const body = await request.json().catch(() => ({}));
  const secret = typeof body.secret === 'string' ? body.secret : '';
  const expected = process.env.ADMIN_SECRET ?? 'dev_admin_secret_change_in_production';
  if (secret !== expected) {
    return NextResponse.json({ ok: false }, { status: 401 });
  }
  const cookieStore = await cookies();
  cookieStore.set(ADMIN_COOKIE, secret, { path: '/', maxAge: MAX_AGE, sameSite: 'lax', httpOnly: true });
  return NextResponse.json({ ok: true });
}
