export const getImageDimensions = (blob: Blob): Promise<{ width: number; height: number }> => {
  return new Promise((resolve, reject) => {
    const img = new Image();
    const objectUrl = URL.createObjectURL(blob);

    img.src = objectUrl;

    // 成功時
    img.onload = () => {
      resolve({ width: img.width, height: img.height });
      URL.revokeObjectURL(objectUrl); // メモリ解放
    };

    // 失敗時
    img.onerror = () => {
      reject(new Error('Failed to load image for dimension check'));
      URL.revokeObjectURL(objectUrl); // メモリ解放
    };
  });
};
