import { renderHook } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { usePreviewUrl } from './usePreviewUrl';

describe('usePreviewUrl', () => {
  it('returns empty string when blob is null', () => {
    const { result } = renderHook(() =>
      usePreviewUrl(null),
    );

    expect(result.current).toBe('');
  });

  it('creates preview URL from blob', () => {
    const createObjectURL = vi
      .spyOn(URL, 'createObjectURL')
      .mockReturnValue('blob:test');

    const blob = new Blob(['image']);

    const { result } = renderHook(() =>
      usePreviewUrl(blob),
    );

    expect(createObjectURL).toHaveBeenCalledWith(blob);
    expect(result.current).toBe('blob:test');

    createObjectURL.mockRestore();
  });

  it('revokes preview URL on unmount', () => {
    vi.spyOn(URL, 'createObjectURL').mockReturnValue('blob:test');

    const revokeObjectURL = vi.spyOn(
      URL,
      'revokeObjectURL',
    );

    const blob = new Blob(['image']);

    const { unmount } = renderHook(() =>
      usePreviewUrl(blob),
    );

    unmount();

    expect(revokeObjectURL).toHaveBeenCalledWith(
      'blob:test',
    );
  });
});
