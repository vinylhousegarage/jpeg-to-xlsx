import { render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { PreviewArea } from '../PreviewArea/PreviewArea';
import * as UsePreviewUrlModule from '../../../hooks/usePreviewUrl';

vi.mock('../../../hooks/usePreviewUrl', () => ({
  usePreviewUrl: vi.fn(),
}));

describe('PreviewArea', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders the preview image when an image URL is available', () => {
    const blob = new Blob(['test-image'], {
      type: 'image/jpeg',
    });

    vi.mocked(
      UsePreviewUrlModule.usePreviewUrl,
    ).mockReturnValue('blob:http://localhost/preview');

    render(<PreviewArea blob={blob} />);

    const image = screen.getByRole('img', {
      name: '撮影画像',
    });

    expect(image).toBeInTheDocument();
    expect(image).toHaveAttribute(
      'src',
      'blob:http://localhost/preview',
    );

    expect(
      UsePreviewUrlModule.usePreviewUrl,
    ).toHaveBeenCalledWith(blob);
  });

  it('renders nothing when an image URL is unavailable', () => {
    const blob = new Blob(['test-image'], {
      type: 'image/jpeg',
    });

    vi.mocked(
      UsePreviewUrlModule.usePreviewUrl,
    ).mockReturnValue('');

    const { container } = render(
      <PreviewArea blob={blob} />,
    );

    expect(container).toBeEmptyDOMElement();

    expect(
      screen.queryByRole('img', {
        name: '撮影画像',
      }),
    ).not.toBeInTheDocument();
  });
});
