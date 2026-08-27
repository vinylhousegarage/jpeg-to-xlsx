let counter = 1;

export const createShotNumber = (): string => {
  const shotNumber = `SHOT-${String(counter).padStart(3, '0')}`;
  counter += 1;

  return shotNumber;
};
