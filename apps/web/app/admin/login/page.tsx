'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Button, Input } from '@marketplace/ui';

export default function AdminLoginPage() {
  const router = useRouter();
  const [secret, setSecret] = useState('');
  const [error, setError] = useState('');

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError('');
    fetch('/api/admin/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ secret }),
    })
      .then((res) => res.json())
      .then((data) => {
        if (data.ok) {
          router.push('/admin');
          router.refresh();
        } else {
          setError('Неверный ключ');
        }
      })
      .catch(() => setError('Ошибка'));
  }

  return (
    <div className="container flex min-h-[60vh] items-center justify-center">
      <form onSubmit={handleSubmit} className="w-full max-w-sm space-y-4">
        <h1 className="text-2xl font-bold text-center">Вход в админ-панель</h1>
        <div>
          <label htmlFor="secret" className="text-sm font-medium">
            Секретный ключ
          </label>
          <Input
            id="secret"
            type="password"
            value={secret}
            onChange={(e) => setSecret(e.target.value)}
            required
            className="mt-1"
          />
        </div>
        {error && <p className="text-sm text-destructive">{error}</p>}
        <Button type="submit" className="w-full">
          Войти
        </Button>
      </form>
    </div>
  );
}
