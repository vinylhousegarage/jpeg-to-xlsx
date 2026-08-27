import { test, expect } from '@playwright/test';
import { fileURLToPath } from 'url';
import { dirname, resolve } from 'path';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

test('has title', async ({ page }) => {
  await page.goto('file://' + resolve(__dirname, '../index.html'));
  await expect(page).toHaveTitle(/React App/);
});
