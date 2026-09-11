import { PreviewArea } from './PreviewArea';
import { PreviewActions } from './PreviewActions';

type Props = {
  blob: Blob;
  onRetakeFileSelected: (
    file: File,
  ) => Promise<void>;
  onSend: () => void;
  onCancel: () => void;
  isSending?: boolean;
  isCompressing?: boolean;
};

export const PreviewPhase = ({
  blob,
  onRetakeFileSelected,
  onSend,
  onCancel,
  isSending = false,
  isCompressing = false,
}: Props) => {
  return (
    <div
      style={{
        maxWidth: '375px',
        margin: '0 auto',
        textAlign: 'center',
      }}
    >
      <h2>
        {isCompressing
          ? '画像を処理しています'
          : '画像を確認'}
      </h2>

      <PreviewArea blob={blob} />

      <PreviewActions
        onRetakeFileSelected={
          onRetakeFileSelected
        }
        onSubmit={onSend}
        onCancel={onCancel}
        isSending={isSending}
        isCompressing={isCompressing}
      />
    </div>
  );
};
