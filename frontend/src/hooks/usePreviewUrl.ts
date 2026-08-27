import { useState, useEffect } from 'react';

export const usePreviewUrl = (blob: Blob | null) => {
  const [imageUrl, setImageUrl] = useState<string>('');

  useEffect(() => {
    if (!blob) {
      setImageUrl('');
      return;
    }
    const url = URL.createObjectURL(blob);
    setImageUrl(url);
    return () => URL.revokeObjectURL(url);
  }, [blob]);

  return imageUrl;
};
