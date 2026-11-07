import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { createRefetchScheduler } from "./refetch";

describe("createRefetchScheduler", () => {
  beforeEach(() => vi.useFakeTimers());
  afterEach(() => vi.useRealTimers());

  it("batches changes that arrive close together", () => {
    const refetch = vi.fn();
    const scheduler = createRefetchScheduler<string>(refetch);
    scheduler.change("a");
    vi.advanceTimersByTime(100);
    scheduler.change("b");
    vi.advanceTimersByTime(199);
    expect(refetch).not.toHaveBeenCalled();
    vi.advanceTimersByTime(1);
    expect(refetch).toHaveBeenCalledExactlyOnceWith(["a", "b"]);
  });

  it("spreads a ready refetch over the jitter window", () => {
    const refetch = vi.fn();
    const scheduler = createRefetchScheduler<string>(refetch, {
      random: () => 0.5,
    });
    scheduler.ready();
    vi.advanceTimersByTime(2499);
    expect(refetch).not.toHaveBeenCalled();
    vi.advanceTimersByTime(1);
    expect(refetch).toHaveBeenCalledExactlyOnceWith(null);
  });

  it("refetches everything early when a change arrives before a ready refetch", () => {
    const refetch = vi.fn();
    const scheduler = createRefetchScheduler<string>(refetch, {
      random: () => 0.9,
    });
    scheduler.ready();
    scheduler.change("a");
    vi.advanceTimersByTime(300);
    expect(refetch).toHaveBeenCalledExactlyOnceWith(null);
    vi.advanceTimersByTime(10_000);
    expect(refetch).toHaveBeenCalledOnce();
  });

  it("drops a scheduled refetch on cancel", () => {
    const refetch = vi.fn();
    const scheduler = createRefetchScheduler<string>(refetch);
    scheduler.change("a");
    scheduler.cancel();
    vi.advanceTimersByTime(10_000);
    expect(refetch).not.toHaveBeenCalled();
  });
});
