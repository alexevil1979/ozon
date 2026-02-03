/** Next.js app ESLint. Extends shared config + Next. */
module.exports = {
  root: true,
  extends: ['@marketplace/eslint-config', 'next/core-web-vitals', 'next/typescript'],
  ignorePatterns: ['.next', 'node_modules'],
};
