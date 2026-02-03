'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Button, Input } from '@marketplace/ui';

const AUTH_API = process.env.NEXT_PUBLIC_AUTH_URL ?? 'http://localhost:8080';

export function LoginForm() {
  const router = useRouter();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      const res = await fetch(`${AUTH_API}/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password }),
      });
      const data = await res.json().catch(() => ({}));
      if (!res.ok) {
        setError(data.error ?? 'Ошибка входа');
        return;
      }
      const token = data.access_token;
      if (token) {
        document.cookie = `marketplace_token=${token}; path=/; max-age=86400; samesite=lax`;
      }
      const next = new URLSearchParams(typeof window !== 'undefined' ? window.location.search : '').get('next') ?? '/seller';
      router.push(next);
      router.refresh();
    } catch {
      setError('Сервис авторизации недоступен');
    } finally {
      setLoading(false);
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div>
        <label htmlFor="email" className="text-sm font-medium">
          Email
        </label>
        <Input
          id="email"
          name="email"
          type="email"
          autoComplete="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
          className="mt-1 w-full"
        />
      </div>
      <div>
        <label htmlFor="password" className="text-sm font-medium">
          Пароль
        </label>
        <Input
          id="password"
          name="password"
          type="password"
          autoComplete="current-password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
          className="mt-1 w-full"
        />
      </div>
      {error ? <p className="text-sm text-destructive">{error}</p> : null}
      <Button type="submit" className="w-full" disabled={loading}>
        {loading ? 'Вход…' : 'Войти'}
      </Button>
    </form>
  );
}
