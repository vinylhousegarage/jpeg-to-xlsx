import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { SlackOAuth } from './SlackOAuth';

describe('SlackOAuth', () => {
  it('renders the Slack notification settings heading', () => {
    render(
      <SlackOAuth
        onConnectSlack={vi.fn()}
      />,
    );

    expect(
      screen.getByRole('heading', {
        name: 'Slack通知設定',
      }),
    ).toBeInTheDocument();
  });

  it('calls onConnectSlack when the button is clicked', () => {
    const onConnectSlack = vi.fn();

    render(
      <SlackOAuth
        onConnectSlack={onConnectSlack}
      />,
    );

    const button = screen.getByRole('button', {
      name: 'DM通知を設定',
    });

    fireEvent.click(button);

    expect(onConnectSlack).toHaveBeenCalledTimes(1);
  });
});
