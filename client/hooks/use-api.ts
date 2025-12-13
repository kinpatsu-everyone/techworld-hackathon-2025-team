import { useSyncExternalStore, useCallback, useRef } from "react";
import { api, EndpointTypes, ApiResponse, ApiError } from "@/lib/client";

// ============================================================================
// Cache Store
// ============================================================================

type CacheEntry<T> = {
  data: T | null;
  error: ApiError | null;
  isLoading: boolean;
  timestamp: number;
};

type CacheStore = Map<string, CacheEntry<unknown>>;
type Listeners = Set<() => void>;

const cache: CacheStore = new Map();
const listeners: Listeners = new Set();

function emitChange() {
  listeners.forEach((listener) => listener());
}

function subscribe(listener: () => void) {
  listeners.add(listener);
  return () => listeners.delete(listener);
}

function getCacheKey<P extends keyof EndpointTypes>(
  endpoint: P,
  request: EndpointTypes[P]["request"]
): string {
  return `${endpoint}:${JSON.stringify(request)}`;
}

function getSnapshot() {
  return cache;
}

function getServerSnapshot() {
  return cache;
}

// ============================================================================
// API State Type
// ============================================================================

export type ApiState<T> = {
  data: T | null;
  error: ApiError | null;
  isLoading: boolean;
  mutate: () => Promise<void>;
};

// ============================================================================
// useApi Hook
// ============================================================================

export function useApi<P extends keyof EndpointTypes>(
  endpoint: P,
  request: EndpointTypes[P]["request"],
  options?: {
    headers?: Record<string, string>;
    enabled?: boolean;
  }
): ApiState<EndpointTypes[P]["response"]> {
  const cacheKey = getCacheKey(endpoint, request);
  const enabled = options?.enabled ?? true;
  const fetchingRef = useRef(false);

  const fetchData = useCallback(async () => {
    if (fetchingRef.current) return;
    fetchingRef.current = true;

    // Set loading state
    cache.set(cacheKey, {
      data: cache.get(cacheKey)?.data ?? null,
      error: null,
      isLoading: true,
      timestamp: Date.now(),
    });
    emitChange();

    try {
      const response: ApiResponse<EndpointTypes[P]["response"]> = await api(
        endpoint,
        request,
        { headers: options?.headers }
      );

      cache.set(cacheKey, {
        data: response.data,
        error: null,
        isLoading: false,
        timestamp: Date.now(),
      });
    } catch (err) {
      cache.set(cacheKey, {
        data: cache.get(cacheKey)?.data ?? null,
        error: err instanceof ApiError ? err : new ApiError(0, "UNKNOWN", String(err)),
        isLoading: false,
        timestamp: Date.now(),
      });
    } finally {
      fetchingRef.current = false;
      emitChange();
    }
  }, [cacheKey, endpoint, request, options?.headers]);

  // Subscribe to cache changes
  const currentCache = useSyncExternalStore(subscribe, getSnapshot, getServerSnapshot);

  // Initial fetch (only if enabled and no cache entry exists)
  const entry = currentCache.get(cacheKey) as CacheEntry<EndpointTypes[P]["response"]> | undefined;

  if (enabled && !entry && !fetchingRef.current) {
    // Trigger fetch on first render
    fetchData();
  }

  const mutate = useCallback(async () => {
    await fetchData();
  }, [fetchData]);

  return {
    data: entry?.data ?? null,
    error: entry?.error ?? null,
    isLoading: entry?.isLoading ?? (enabled && !entry),
    mutate,
  };
}

// ============================================================================
// Global Mutate Function
// ============================================================================

export function mutateApi<P extends keyof EndpointTypes>(
  endpoint: P,
  request: EndpointTypes[P]["request"]
): void {
  const cacheKey = getCacheKey(endpoint, request);
  cache.delete(cacheKey);
  emitChange();
}

export function clearApiCache(): void {
  cache.clear();
  emitChange();
}
