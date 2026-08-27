import { expect, test } from '@playwright/test';

test.describe('Slack OAuth Phase', () => {
  test.use({
    viewport: {
      width: 375,
      height: 667,
    },
  });

  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test('should display the Slack notification settings', async ({
    page,
  }) => {
    await expect(
      page.getByRole('heading', {
        name: 'Slack通知設定',
      }),
    ).toBeVisible();

    await expect(
      page.getByRole('button', {
        name: 'DM通知を設定',
      }),
    ).toBeVisible();
  });

  test('should be centered and stacked vertically', async ({
    page,
  }) => {
    const container = page.locator('.slack-oauth');

    const heading = page.getByRole('heading', {
      name: 'Slack通知設定',
    });

    const button = page.getByRole('button', {
      name: 'DM通知を設定',
    });

    await expect(container).toBeVisible();
    await expect(heading).toBeVisible();
    await expect(button).toBeVisible();

    const [containerBox, headingBox, buttonBox] =
      await Promise.all([
        container.boundingBox(),
        heading.boundingBox(),
        button.boundingBox(),
      ]);

    if (!containerBox || !headingBox || !buttonBox) {
      throw new Error('Some elements are not visible');
    }

    const getCenter = (box: {
      x: number;
      width: number;
    }) => box.x + box.width / 2;

    const containerCenter = getCenter(containerBox);
    const headingCenter = getCenter(headingBox);
    const buttonCenter = getCenter(buttonBox);

    expect(
      Math.abs(containerCenter - headingCenter),
    ).toBeLessThan(1);

    expect(
      Math.abs(containerCenter - buttonCenter),
    ).toBeLessThan(1);

    expect(
      headingBox.y + headingBox.height,
    ).toBeLessThanOrEqual(buttonBox.y);
  });
});
