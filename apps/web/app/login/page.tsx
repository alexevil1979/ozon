import type { Metadata } from 'next';
import { LoginForm } from '@/components/auth/LoginForm';

export const metadata: Metadata = {
  title: 'Вход',
  description: 'Вход в аккаунт Marketplace',
};

export default function LoginPage() {
  return (
    <div className="container flex min-h-[60vh] flex-col items-center justify-center">
      <div className="w-full max-w-sm space-y-6">
        <h1 className="text-2xl font-bold text-center">Вход</h1>
        <LoginForm />
        <p className="text-center text-sm text-muted-foreground">
          <a href="/" className="hover:underline">На главную</a>
        </p>
      </div>
    </div>
  );
}
