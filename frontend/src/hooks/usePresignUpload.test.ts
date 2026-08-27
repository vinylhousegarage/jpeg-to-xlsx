import {
  beforeEach,
  describe,
  expect,
  it,
  vi,
} from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { usePresignUpload } from './usePresignUpload';
import * as UseS3UploadModule from './useS3Upload';

vi.mock('./useS3Upload', () => ({
  useS3Upload: vi.fn(),
}));

describe('usePresignUpload', () => {
  const mockDispatch = vi.fn();
  const mockUpload = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();

    vi.mocked(
      UseS3UploadModule.useS3Upload,
    ).mockReturnValue({
      upload: mockUpload,
    });

    vi.stubGlobal(
      'fetch',
      vi.fn(),
    );
  });

  it('requests a presigned URL and uploads the image', async () => {
    const blob = new Blob(
      ['image'],
      {
        type: 'image/jpeg',
      },
    );

    const shotNumber = 'SHOT-001';
    const uploadURL =
      'https://example.com/presigned-upload';

    vi.mocked(fetch).mockResolvedValue(
      new Response(
        JSON.stringify({
          uploadURL,
          expiresAt: '2026-08-07T15:00:00Z',
        }),
        {
          status: 200,
          headers: {
            'Content-Type': 'application/json',
          },
        },
      ),
    );

    mockUpload.mockResolvedValue(undefined);

    const { result } = renderHook(() =>
      usePresignUpload(mockDispatch),
    );

    await act(async () => {
      await result.current.send(
        blob,
        shotNumber,
      );
    });

    expect(fetch).toHaveBeenCalledTimes(1);

    expect(fetch).toHaveBeenCalledWith(
      `${window.location.origin}/api/storage/upload`,
      {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          shotNumber: 'SHOT-001',
        }),
      },
    );

    expect(mockUpload).toHaveBeenCalledTimes(1);

    expect(mockUpload).toHaveBeenCalledWith(
      blob,
      uploadURL,
      shotNumber,
    );
  });

  it('dispatches an error when presigned URL acquisition fails', async () => {
    const blob = new Blob(
      ['image'],
      {
        type: 'image/jpeg',
      },
    );

    const shotNumber = 'SHOT-001';

    vi.mocked(fetch).mockResolvedValue(
      new Response(null, {
        status: 500,
      }),
    );

    const { result } = renderHook(() =>
      usePresignUpload(mockDispatch),
    );

    await act(async () => {
      await result.current.send(
        blob,
        shotNumber,
      );
    });

    expect(mockDispatch).toHaveBeenCalledWith({
      type: 'UPLOAD_COMPLETE',
      status: 'error',
      error: expect.objectContaining({
        message:
          'Failed to get presigned URL: 500',
      }),
    });

    expect(mockUpload).not.toHaveBeenCalled();
  });

  it('does not upload when the presign request fails', async () => {
    const blob = new Blob(
      ['image'],
      {
        type: 'image/jpeg',
      },
    );

    vi.mocked(fetch).mockResolvedValue(
      new Response(null, {
        status: 400,
      }),
    );

    const { result } = renderHook(() =>
      usePresignUpload(mockDispatch),
    );

    await act(async () => {
      await result.current.send(
        blob,
        'SHOT-001',
      );
    });

    expect(mockUpload).not.toHaveBeenCalled();
  });
});
