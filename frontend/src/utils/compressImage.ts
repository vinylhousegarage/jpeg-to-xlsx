import { getImageDimensions } from './getImageDimensions';

export const compressImage = async (file: File | Blob): Promise<Blob> => {
  // 画像の元サイズを取得
  const dimensions = await getImageDimensions(file);

  // 縮小が必要か判定し、新しい縦横サイズ（アスペクト比維持）を計算
  const MAX_LONG_EDGE = 2000;
  let targetWidth = dimensions.width;
  let targetHeight = dimensions.height;

  if (targetWidth > MAX_LONG_EDGE || targetHeight > MAX_LONG_EDGE) {
    if (targetWidth > targetHeight) {
      // 横長画像の場合：横幅を 2000px に固定し、縦幅を縮小
      targetHeight = Math.round((targetHeight * MAX_LONG_EDGE) / targetWidth);
      targetWidth = MAX_LONG_EDGE;
    } else {
      // 縦長画像（または正方形）の場合：縦幅を 2000px に固定し、横幅を縮小
      targetWidth = Math.round((targetWidth * MAX_LONG_EDGE) / targetHeight);
      targetHeight = MAX_LONG_EDGE;
    }
  }

  // Canvas を使って画像をリサイズ描画
  return new Promise((resolve, reject) => {
    const img = new Image();
    const objectUrl = URL.createObjectURL(file);

    img.src = objectUrl;

    img.onload = () => {
      // メモリ解放
      URL.revokeObjectURL(objectUrl);

      const canvas = document.createElement('canvas');
      canvas.width = targetWidth;
      canvas.height = targetHeight;

      const ctx = canvas.getContext('2d');
      if (!ctx) {
        reject(new Error('Failed to get 2D context from canvas'));
        return;
      }

      // 計算したサイズで Canvas に描画
      ctx.drawImage(img, 0, 0, targetWidth, targetHeight);

      // 品質 85%（0.85）の image/jpeg 形式で Blob として書き出し
      canvas.toBlob(
        (blob) => {
          if (blob) {
            resolve(blob);
          } else {
            reject(new Error('Canvas toBlob serialization failed'));
          }
        },
        'image/jpeg',
        0.85
      );
    };

    img.onerror = () => {
      URL.revokeObjectURL(objectUrl);
      reject(new Error('Failed to load image for compression'));
    };
  });
};
