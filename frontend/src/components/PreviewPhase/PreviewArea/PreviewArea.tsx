import { usePreviewUrl } from '../../../hooks/usePreviewUrl';

type Props = {
  blob: Blob;
};

export const PreviewArea = ({
  blob,
}: Props) => {
  const imageUrl = usePreviewUrl(blob);

  if (!imageUrl) {
    return null;
  }

  return (
    <div className="preview-container">
      <img
        src={imageUrl}
        alt="撮影画像"
        style={{
          maxWidth: '100%',
          display: 'block',
          marginBottom: '1rem',
        }}
      />
    </div>
  );
};
