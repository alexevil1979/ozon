import type { NextConfig } from 'next';
import path from 'path';

/** Next.js 15 App Router. Transpile workspace packages. */
const nextConfig: NextConfig = {
  reactStrictMode: true,
  transpilePackages: ['@marketplace/ui', '@marketplace/types', '@marketplace/utils'],
  experimental: {},
};

export default nextConfig;
