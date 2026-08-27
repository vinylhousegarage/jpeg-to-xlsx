import React from 'react';
import { standardButtonStyle } from '../../../styles/button';

type Props = {
  onFileSelected: (file: File) => Promise<void>;
  disabled?: boolean;
};

export const CameraInput: React.FC<Props> = ({
  onFileSelected,
  disabled = false,
}) => {
  const handleFileChange = async (
    event: React.ChangeEvent<HTMLInputElement>,
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
    <label style={standardButtonStyle}>
      カメラを起動

      <input
        type="file"
        accept="image/*"
        capture="environment"
        onChange={handleFileChange}
        disabled={disabled}
        style={{ display: 'none' }}
      />
    </label>
  );
};
