/** @type {import('eslint').Linter.Config} */
module.exports = {
  root: true,
  extends: ['@marketplace/eslint-config'],
  ignorePatterns: ['node_modules', 'dist', '.next', '*.cjs', 'services/**'],
};
