import { useRef, type ChangeEvent } from 'react';
import { standardButtonStyle } from '../../../styles/button';

type Props = {
  onRetakeFileSelected: (
    file: File,
  ) => Promise<void>;
  onSubmit: () => void;
  isSending?: boolean;
  isCompressing?: boolean;
};

export const PreviewActions = ({
  onRetakeFileSelected,
  onSubmit,
  isSending = false,
  isCompressing = false,
}: Props) => {
  const fileInputRef =
    useRef<HTMLInputElement>(null);

  const disabled =
    isSending || isCompressing;

  const handleRetake = () => {
    fileInputRef.current?.click();
  };

  const handleFileChange = async (
    event: ChangeEvent<HTMLInputElement>,
  ) => {
    const file =
      event.target.files?.[0];

    if (!file) {
      return;
    }

    try {
      await onRetakeFileSelected(file);
    } finally {
      event.target.value = '';
    }
  };

  return (
    <div
      className="button-group"
      style={{
        display: 'flex',
        justifyContent: 'center',
        gap: '10px',
        width: '100%',
      }}
    >
      <button
        type="button"
        onClick={handleRetake}
        disabled={disabled}
        style={standardButtonStyle}
      >
        撮り直し
      </button>

      <input
        ref={fileInputRef}
        type="file"
        accept="image/*"
        capture="environment"
        onChange={handleFileChange}
        disabled={disabled}
        style={{ display: 'none' }}
      />

      <button
        type="button"
        onClick={onSubmit}
        disabled={disabled}
        style={standardButtonStyle}
      >
        画像を確定し送信
      </button>
    </div>
  );
};
