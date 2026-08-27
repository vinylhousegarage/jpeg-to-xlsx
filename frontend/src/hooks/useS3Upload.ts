import { Dispatch } from 'react';
import { AppAction } from '../types';

export const useS3Upload = (
  dispatch: Dispatch<AppAction>,
) => {
  const upload = async (
    blob: Blob,
    presignUrl: string,
    shotNumber: string,
  ): Promise<void> => {
    dispatch({
      type: 'START_UPLOAD',
    });

    try {
      const response = await fetch(presignUrl, {
        method: 'PUT',
        body: blob,
        headers: {
          'Content-Type': 'image/jpeg',
          'x-amz-meta-shot-number': shotNumber,
        },
      });

      if (!response.ok) {
        throw new Error(
          `Failed to upload file: ${response.status}`,
        );
      }

      dispatch({
        type: 'UPLOAD_COMPLETE',
        status: 'success',
      });
    } catch (error) {
      dispatch({
        type: 'UPLOAD_COMPLETE',
        status: 'error',
        error:
          error instanceof Error
            ? error
            : new Error('Failed to upload file'),
      });
    }
  };

  return {
    upload,
  };
};
