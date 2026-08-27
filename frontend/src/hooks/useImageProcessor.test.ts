import { act, renderHook } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { useImageProcessor } from './useImageProcessor';
import * as imageUtils from '../components/PreviewPhase/utils/compressImage';

vi.mock('../components/PreviewPhase/utils/compressImage', () => ({
  compressImage: vi.fn(),
}));

describe('useImageProcessor', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('compresses the image and passes the compressed Blob to onCapture', async () => {
    const compressedBlob = new Blob(
      ['compressed-image'],
      { type: 'image/jpeg' },
    );
    const file = new File(
      ['original-image'],
      'test.png',
      { type: 'image/png' },
    );
    const onCapture = vi.fn();

    vi.mocked(
      imageUtils.compressImage,
    ).mockResolvedValue(compressedBlob);

    const { result } = renderHook(() =>
      useImageProcessor(onCapture),
    );

    await act(async () => {
      await result.current.processImage(file);
    });

    expect(imageUtils.compressImage).toHaveBeenCalledTimes(1);
    expect(imageUtils.compressImage).toHaveBeenCalledWith(file);

    expect(onCapture).toHaveBeenCalledTimes(1);
    expect(onCapture).toHaveBeenCalledWith(compressedBlob);

    expect(result.current.isCompressing).toBe(false);
  });
});
