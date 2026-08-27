import { Dispatch } from 'react';
import { AppAction } from '../types';
import { useS3Upload } from './useS3Upload';

export const usePresignUpload = (
  dispatch: Dispatch<AppAction>,
) => {
  const { upload } = useS3Upload(dispatch);

  const send = async (
    blob: Blob,
    shotNumber: string,
  ): Promise<void> => {
    const apiBaseURL = (
      import.meta.env.VITE_API_BASE_URL || window.location.origin
    ).replace(/\/$/, '');

    const response = await fetch(
      `${apiBaseURL}/api/storage/upload`,
      {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ shotNumber }),
      },
    );

    if (!response.ok) {
      dispatch({
        type: 'UPLOAD_COMPLETE',
        status: 'error',
        error: new Error(
          `Failed to get presigned URL: ${response.status}`,
        ),
      });

      return;
    }

    const data: {
      uploadURL: string;
      expiresAt: string;
    } = await response.json();

    await upload(
      blob,
      data.uploadURL,
      shotNumber,
    );
  };

  return {
    send,
  };
};
