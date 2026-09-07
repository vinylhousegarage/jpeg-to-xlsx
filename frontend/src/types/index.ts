export type InputPhase = {
  type: 'input';
};

export type PreviewPhase = {
  type: 'preview';
  shotNumber: string;
  file: Blob;
};

export type UploadPhase = {
  type: 'upload';
  shotNumber: string;
};

export type ResultPhase = {
  type: 'result';
  shotNumber: string;
  status: 'success' | 'error';
  error?: Error;
};

export type AppPhase =
  | InputPhase
  | PreviewPhase
  | UploadPhase
  | ResultPhase;

export type AppState = {
  isSlackLinked: boolean;
  phase: AppPhase;
};

export type AppAction =
  | {
      type: 'SET_SLACK_LINKED';
      isSlackLinked: boolean;
    }
  | {
      type: 'SET_PREVIEW';
      file: Blob;
      shotNumber: string;
    }
  | { type: 'SEND' }
  | { type: 'START_UPLOAD' }
  | {
      type: 'UPLOAD_COMPLETE';
      status: 'success' | 'error';
      error?: Error;
    }
