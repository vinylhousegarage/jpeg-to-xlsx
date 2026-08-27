import React from 'react';
import { standardButtonStyle } from '../../../styles/button';

type Props = {
  onConnectSlack: () => void;
};

export const SlackOAuth: React.FC<Props> = ({
  onConnectSlack,
}) => {
  return (
    <div
      className="slack-oauth"
      style={{
        maxWidth: '375px',
        margin: '0 auto',
        textAlign: 'center',
      }}
    >
      <h2>Slack通知設定</h2>

      <div
        style={{
          display: 'flex',
          flexDirection: 'column',
          gap: '10px',
          alignItems: 'center',
        }}
      >
        <button
          type="button"
          onClick={onConnectSlack}
          style={standardButtonStyle}
        >
          DM通知を設定
        </button>
      </div>
    </div>
  );
};
