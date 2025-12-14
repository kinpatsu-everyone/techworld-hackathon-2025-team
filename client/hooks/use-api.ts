import { useSyncExternalStore, useCallback, useRef, useEffect } from 'react';
import { api, EndpointTypes, ApiResponse, ApiError } from '@/lib/client';

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
type Listeners = Map<string, Set<() => void>>;

const cache: CacheStore = new Map();
const listeners: Listeners = new Map();

function emitChange(cacheKey: string) {
  const keyListeners = listeners.get(cacheKey);
  if (keyListeners) {
    keyListeners.forEach((listener) => listener());
  }
}

function subscribe(cacheKey: string, listener: () => void) {
  if (!listeners.has(cacheKey)) {
    listeners.set(cacheKey, new Set());
  }
  listeners.get(cacheKey)!.add(listener);
  return () => {
    listeners.get(cacheKey)?.delete(listener);
  };
}

function getCacheKey<P extends keyof EndpointTypes>(
  endpoint: P,
  request: EndpointTypes[P]['request']
): string {
  return `${endpoint}:${JSON.stringify(request)}`;
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
  request: EndpointTypes[P]['request'],
  options?: {
    headers?: Record<string, string>;
    enabled?: boolean;
  }
): ApiState<EndpointTypes[P]['response']> {
  const cacheKey = getCacheKey(endpoint, request);
  const enabled = options?.enabled ?? true;
  const fetchingRef = useRef(false);
  const requestRef = useRef(request);
  const headersRef = useRef(options?.headers);

  // Update refs on each render
  requestRef.current = request;
  headersRef.current = options?.headers;

  const fetchData = useCallback(async () => {
    if (fetchingRef.current) return;
    fetchingRef.current = true;

    const currentCacheKey = cacheKey;

    // Set loading state
    cache.set(currentCacheKey, {
      data: cache.get(currentCacheKey)?.data ?? null,
      error: null,
      isLoading: true,
      timestamp: Date.now(),
    });
    emitChange(currentCacheKey);

    try {
      const response: ApiResponse<EndpointTypes[P]['response']> = await api(
        endpoint,
        requestRef.current,
        { headers: headersRef.current }
      );

      cache.set(currentCacheKey, {
        data: response.data,
        error: null,
        isLoading: false,
        timestamp: Date.now(),
      });
    } catch (err) {
      console.error('useApi fetch error:', err);
      cache.set(currentCacheKey, {
        data: cache.get(currentCacheKey)?.data ?? null,
        error:
          err instanceof ApiError
            ? err
            : new ApiError(0, 'UNKNOWN', String(err)),
        isLoading: false,
        timestamp: Date.now(),
      });
    } finally {
      fetchingRef.current = false;
      emitChange(currentCacheKey);
    }
  }, [cacheKey, endpoint]);

  // Subscribe to cache changes for this specific key
  const subscribeToKey = useCallback(
    (listener: () => void) => subscribe(cacheKey, listener),
    [cacheKey]
  );

  const getSnapshot = useCallback(
    () => cache.get(cacheKey) as CacheEntry<EndpointTypes[P]['response']> | undefined,
    [cacheKey]
  );

  const entry = useSyncExternalStore(subscribeToKey, getSnapshot, getSnapshot);

  // Initial fetch - only when cacheKey changes and no entry exists
  useEffect(() => {
    if (enabled && !cache.has(cacheKey)) {
      fetchData();
    }
  }, [enabled, cacheKey, fetchData]);

  const mutate = useCallback(async () => {
    cache.delete(cacheKey);
    fetchingRef.current = false;
    await fetchData();
  }, [cacheKey, fetchData]);

  return {
    data: entry?.data ?? null,
    error: entry?.error ?? null,
    isLoading: entry?.isLoading ?? (enabled && !cache.has(cacheKey)),
    mutate,
  };
}

// ============================================================================
// Global Mutate Function
// ============================================================================

export function mutateApi<P extends keyof EndpointTypes>(
  endpoint: P,
  request: EndpointTypes[P]['request']
): void {
  const cacheKey = getCacheKey(endpoint, request);
  cache.delete(cacheKey);
  emitChange(cacheKey);
}

export function clearApiCache(): void {
  const keys = Array.from(cache.keys());
  cache.clear();
  keys.forEach((key) => emitChange(key));
}
