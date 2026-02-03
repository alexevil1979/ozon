import type { Metadata } from 'next';
import { Header } from '@/components/buyer/Header';
import { Footer } from '@/components/buyer/Footer';

export const metadata: Metadata = {
  title: { default: 'Marketplace', template: '%s | Marketplace' },
  description: 'Покупайте и продавайте. Каталог товаров.',
};

export default function BuyerLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="flex min-h-screen flex-col">
      <Header />
      <main className="flex-1">{children}</main>
      <Footer />
    </div>
  );
}
