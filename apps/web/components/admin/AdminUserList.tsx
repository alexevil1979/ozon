'use client';

import { useEffect, useState } from 'react';

interface UserRow {
  id: string;
  email: string;
  created_at: string;
}

export function AdminUserList() {
  const [users, setUsers] = useState<UserRow[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    fetch('/api/admin/users')
      .then((res) => {
        if (res.status === 401) {
          setError('Доступ запрещён');
          return [];
        }
        if (!res.ok) return [];
        return res.json();
      })
      .then((data) => {
        setUsers(Array.isArray(data) ? data : []);
      })
      .catch(() => setError('Ошибка загрузки'))
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <p className="text-muted-foreground">Загрузка…</p>;
  if (error) return <p className="text-destructive">{error}</p>;
  if (users.length === 0) return <p className="text-muted-foreground">Нет пользователей.</p>;

  return (
    <div className="rounded-lg border">
      <table className="w-full text-sm">
        <thead>
          <tr className="border-b bg-muted/50">
            <th className="text-left p-4 font-medium">ID</th>
            <th className="text-left p-4 font-medium">Email</th>
            <th className="text-left p-4 font-medium">Дата регистрации</th>
          </tr>
        </thead>
        <tbody>
          {users.map((u) => (
            <tr key={u.id} className="border-b last:border-0">
              <td className="p-4 font-mono text-xs">{u.id}</td>
              <td className="p-4">{u.email}</td>
              <td className="p-4 text-muted-foreground">
                {new Date(u.created_at).toLocaleString('ru-RU')}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
