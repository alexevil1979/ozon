import { describe, it, expect } from 'vitest';
import { cn } from '@marketplace/utils';

describe('cn (class names)', () => {
  it('merges class names', () => {
    expect(cn('foo', 'bar')).toBe('foo bar');
  });

  it('handles conditional classes', () => {
    expect(cn('base', false && 'hidden', true && 'visible')).toBe('base visible');
  });
});
