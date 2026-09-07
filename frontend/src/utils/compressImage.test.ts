import 'vitest-canvas-mock';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { compressImage } from './compressImage';

describe('compressImage', () => {
  beforeEach(() => {
    vi.restoreAllMocks();

    // URL をモック
    vi.stubGlobal('URL', {
      createObjectURL: vi.fn(() => 'blob:mock-url'),
      revokeObjectURL: vi.fn(),
    });

    // onload のトリガー
    vi.spyOn(HTMLImageElement.prototype, 'src', 'set').mockImplementation(function (this: HTMLImageElement) {
      queueMicrotask(() => {
        this.dispatchEvent(new Event('load'));
      });
    });
  });

  // =============================================================
  // 変換パイプラインをテスト（Blob → URL → Image → Blob）
  // =============================================================
  it('should successfully output a JPEG Blob through the full conversion pipeline', async () => {
    const file = new File(['dummy'], 'photo.png', { type: 'image/png' });
    const result = await compressImage(file);
    
    expect(result).toBeInstanceOf(Blob);
    expect(result.type).toBe('image/jpeg');
  });

  // =============================================================
  // 2000px 制限の圧縮をテスト
  // =============================================================
  it('should resize the long edge to 2000px or less maintaining aspect ratio', async () => {
    // ブラウザ環境（JSDOM）のImageオブジェクトにダミーの画像サイズ（3000x1500）を設定
    vi.spyOn(HTMLImageElement.prototype, 'width', 'get').mockReturnValue(3000);
    vi.spyOn(HTMLImageElement.prototype, 'height', 'get').mockReturnValue(1500);

    // ブラウザ標準の Canvas API (drawImage) を監視
    const drawImageSpy = vi.spyOn(CanvasRenderingContext2D.prototype, 'drawImage');

    const file = new File(['dummy'], 'large.jpg', { type: 'image/jpeg' });
    await compressImage(file);

    expect(drawImageSpy).toHaveBeenCalled();
    const args = drawImageSpy.mock.calls[0];

    // ctx.drawImage(img, dx, dy, dWidth, dHeight) の引数が 2000x1000 に縮小されているか検証
    expect(args[3]).toBe(2000); 
    expect(args[4]).toBe(1000); 
  });

  // =============================================================
  // 品質 85% 指定の書き出しをテスト
  // =============================================================
  it('should call canvas.toBlob with quality parameter set to 0.85', async () => {
    // ブラウザ標準の Canvas API (toBlob) を監視
    const toBlobSpy = vi.spyOn(HTMLCanvasElement.prototype, 'toBlob');

    const file = new File(['dummy'], 'test.jpg', { type: 'image/jpeg' });
    await compressImage(file);

    expect(toBlobSpy).toHaveBeenCalled();
    const mostRecentCall = toBlobSpy.mock.calls[0];

    // toBlob(callback, type, quality) の第2, 第3引数を検証
    expect(mostRecentCall[1]).toBe('image/jpeg'); 
    expect(mostRecentCall[2]).toBe(0.85);         
  });
});
