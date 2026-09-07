import { useState } from 'react';
import { compressImage } from '../utils';

export const useImageProcessor = (
  onCapture: (blob: Blob) => void,
  onError?: (err: Error) => void,
) => {
  const [isCompressing, setIsCompressing] = useState(false);

  const processImage = async (file: File) => {
    try {
      setIsCompressing(true);

      const compressedBlob = await compressImage(file);

      onCapture(compressedBlob);
    } catch (err) {
      onError?.(
        err instanceof Error
          ? err
          : new Error('圧縮に失敗しました'),
      );
    } finally {
      setIsCompressing(false);
    }
  };

  return {
    isCompressing,
    processImage,
  };
};
