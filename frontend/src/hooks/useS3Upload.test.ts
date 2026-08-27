import {
  beforeEach,
  describe,
  expect,
  it,
  vi,
} from 'vitest';
import { useS3Upload } from './useS3Upload';

describe('useS3Upload', () => {
  const mockDispatch = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();

    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue({
        ok: true,
        status: 200,
      }),
    );
  });

  it('dispatches correctly on successful completion and sends metadata', async () => {
    const { upload } = useS3Upload(mockDispatch);

    const mockBlob = new Blob(
      ['test-image'],
      {
        type: 'image/jpeg',
      },
    );

    const shotNumber = '1';

    await upload(
      mockBlob,
      'https://test-s3-presign.url',
      shotNumber,
    );

    expect(fetch).toHaveBeenCalledWith(
      'https://test-s3-presign.url',
      {
        method: 'PUT',
        body: mockBlob,
        headers: {
          'Content-Type': 'image/jpeg',
          'x-amz-meta-shot-number': '1',
        },
      },
    );

    expect(mockDispatch).toHaveBeenNthCalledWith(
      1,
      {
        type: 'START_UPLOAD',
      },
    );

    expect(mockDispatch).toHaveBeenNthCalledWith(
      2,
      {
        type: 'UPLOAD_COMPLETE',
        status: 'success',
      },
    );
  });
});
