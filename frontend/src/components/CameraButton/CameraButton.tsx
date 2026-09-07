import {
  useRef,
  type ChangeEvent,
  type ReactNode,
} from 'react';

import { standardButtonStyle } from '../../styles/button';

type Props = {
  children: ReactNode;
  onFileSelected: (file: File) => Promise<void>;
  disabled?: boolean;
};

export const CameraButton = ({
  children,
  onFileSelected,
  disabled = false,
}: Props) => {
  const fileInputRef =
    useRef<HTMLInputElement>(null);

  const handleClick = () => {
    fileInputRef.current?.click();
  };

  const handleFileChange = async (
    event: ChangeEvent<HTMLInputElement>,
  ) => {
    const file = event.target.files?.[0];

    if (!file) {
      return;
    }

    try {
      await onFileSelected(file);
    } finally {
      event.target.value = '';
    }
  };

  return (
    <>
      <button
        type="button"
        onClick={handleClick}
        disabled={disabled}
        style={standardButtonStyle}
      >
        {children}
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
    </>
  );
};
