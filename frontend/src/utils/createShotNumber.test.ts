import {
  describe,
  expect,
  it,
  vi,
} from 'vitest';

describe('createShotNumber', () => {
  it(
    'creates sequential zero-padded shot numbers',
    async () => {
      vi.resetModules();

      const { createShotNumber } =
        await import(
          './createShotNumber'
        );

      expect(
        createShotNumber(),
      ).toBe('SHOT-001');

      expect(
        createShotNumber(),
      ).toBe('SHOT-002');

      expect(
        createShotNumber(),
      ).toBe('SHOT-003');
    },
  );
});
