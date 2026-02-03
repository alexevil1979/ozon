/** Next.js app ESLint extend. */
module.exports = {
  extends: ['./index.js', 'next/core-web-vitals', 'next/typescript'],
  rules: {
    '@typescript-eslint/no-explicit-any': 'error',
  },
};
