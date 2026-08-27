import { afterEach, describe, expect, it, vi } from "vitest";

import { getImageDimensions } from "./getImageDimensions";

describe("getImageDimensions", () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("returns the image dimensions and revokes the object URL", async () => {
    const blob = new Blob(["dummy"], { type: "image/jpeg" });
    const objectUrl = "blob:test-image";

    vi.spyOn(URL, "createObjectURL").mockReturnValue(objectUrl);
    const revokeObjectURLSpy = vi
      .spyOn(URL, "revokeObjectURL")
      .mockImplementation(() => {});

    class MockImage {
      width = 1920;
      height = 1080;
      onload: (() => void) | null = null;
      onerror: (() => void) | null = null;

      set src(_value: string) {
        queueMicrotask(() => {
          this.onload?.();
        });
      }
    }

    vi.stubGlobal("Image", MockImage);

    await expect(getImageDimensions(blob)).resolves.toEqual({
      width: 1920,
      height: 1080,
    });

    expect(URL.createObjectURL).toHaveBeenCalledWith(blob);
    expect(revokeObjectURLSpy).toHaveBeenCalledWith(objectUrl);
  });

  it("rejects when the image fails to load and revokes the object URL", async () => {
    const blob = new Blob(["dummy"], { type: "image/jpeg" });
    const objectUrl = "blob:test-image";

    vi.spyOn(URL, "createObjectURL").mockReturnValue(objectUrl);
    const revokeObjectURLSpy = vi
      .spyOn(URL, "revokeObjectURL")
      .mockImplementation(() => {});

    class MockImage {
      width = 0;
      height = 0;
      onload: (() => void) | null = null;
      onerror: (() => void) | null = null;

      set src(_value: string) {
        queueMicrotask(() => {
          this.onerror?.();
        });
      }
    }

    vi.stubGlobal("Image", MockImage);

    await expect(getImageDimensions(blob)).rejects.toThrow(
      "Failed to load image for dimension check",
    );

    expect(URL.createObjectURL).toHaveBeenCalledWith(blob);
    expect(revokeObjectURLSpy).toHaveBeenCalledWith(objectUrl);
  });
});
