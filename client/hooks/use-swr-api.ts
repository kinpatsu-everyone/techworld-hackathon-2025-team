import useSWR, { SWRConfiguration, mutate as swrMutate } from 'swr';
import {
  api,
  EndpointTypes,
  ApiError,
  Endpoints,
} from '@/lib/client';

// ============================================================================
// Types
// ============================================================================

export type SWRApiState<T> = {
  data: T | null;
  error: ApiError | null;
  isLoading: boolean;
  isValidating: boolean;
  mutate: () => Promise<T | undefined>;
};

export type UseSWRApiOptions<P extends keyof EndpointTypes> = {
  headers?: Record<string, string>;
  enabled?: boolean;
} & Omit<SWRConfiguration<EndpointTypes[P]['response'], ApiError>, 'fetcher'>;

// ============================================================================
// Helper Functions
// ============================================================================

function getCacheKey<P extends keyof EndpointTypes>(
  endpoint: P,
  request: EndpointTypes[P]['request']
): string {
  return `${endpoint}:${JSON.stringify(request)}`;
}

// ============================================================================
// Generic useApi Hook with SWR
// ============================================================================

/**
 * Type-safe SWR hook for API calls.
 * The request and response types are automatically inferred from the endpoint path.
 *
 * @example
 * ```typescript
 * const { data, error, isLoading } = useApi(Endpoints.GetMonsters, {});
 * // data is typed as GetMonstersResponse | null
 * ```
 */
export function useApi<P extends keyof EndpointTypes>(
  endpoint: P,
  request: EndpointTypes[P]['request'],
  options?: UseSWRApiOptions<P>
): SWRApiState<EndpointTypes[P]['response']> {
  const { headers, enabled = true, ...swrOptions } = options ?? {};
  const cacheKey = enabled ? getCacheKey(endpoint, request) : null;

  const fetcher = async (): Promise<EndpointTypes[P]['response']> => {
    const response = await api(endpoint, request, { headers });
    return response.data;
  };

  const { data, error, isLoading, isValidating, mutate } = useSWR<
    EndpointTypes[P]['response'],
    ApiError
  >(cacheKey, fetcher, swrOptions);

  return {
    data: data ?? null,
    error: error ?? null,
    isLoading,
    isValidating,
    mutate: async () => mutate(),
  };
}

// ============================================================================
// Pre-built Endpoint Hooks
// ============================================================================

/**
 * Fetch health status
 * @example
 * ```typescript
 * const { data, isLoading } = useHealthz();
 * ```
 */
export function useHealthz(options?: Omit<UseSWRApiOptions<'/healthz/v1/Healthz'>, 'headers'>) {
  return useApi(Endpoints.Healthz, {}, options);
}

/**
 * Fetch all monsters
 * @example
 * ```typescript
 * const { data, isLoading, error } = useMonsters();
 * const monsters = data?.monsters ?? [];
 * ```
 */
export function useMonsters(options?: UseSWRApiOptions<'/monster/v1/GetMonsters'>) {
  return useApi(Endpoints.GetMonsters, {}, options);
}

/**
 * Fetch a single monster by ID
 * @example
 * ```typescript
 * const { data, isLoading } = useMonster({ id: 'monster-123' });
 * ```
 */
export function useMonster(
  request: EndpointTypes['/monster/v1/GetMonster']['request'],
  options?: UseSWRApiOptions<'/monster/v1/GetMonster'>
) {
  return useApi(Endpoints.GetMonster, request, options);
}

/**
 * Fetch all trash items
 * @example
 * ```typescript
 * const { data, isLoading } = useTrashs();
 * const trashs = data?.trashs ?? [];
 * ```
 */
export function useTrashs(options?: UseSWRApiOptions<'/trash/v1/GetTrashs'>) {
  return useApi(Endpoints.GetTrashs, {}, options);
}

// ============================================================================
// Global Mutate Functions
// ============================================================================

/**
 * Mutate (revalidate) a specific API cache
 * @example
 * ```typescript
 * await mutateApi(Endpoints.GetMonsters, {});
 * ```
 */
export async function mutateApi<P extends keyof EndpointTypes>(
  endpoint: P,
  request: EndpointTypes[P]['request']
): Promise<void> {
  const cacheKey = getCacheKey(endpoint, request);
  await swrMutate(cacheKey);
}

/**
 * Mutate all monsters cache
 */
export async function mutateMonsters(): Promise<void> {
  await mutateApi(Endpoints.GetMonsters, {});
}

/**
 * Mutate a specific monster cache
 */
export async function mutateMonster(id: string): Promise<void> {
  await mutateApi(Endpoints.GetMonster, { id });
}

/**
 * Mutate all trashs cache
 */
export async function mutateTrashs(): Promise<void> {
  await mutateApi(Endpoints.GetTrashs, {});
}

/**
 * Clear all SWR cache
 * @example
 * ```typescript
 * clearSWRCache();
 * ```
 */
export function clearSWRCache(): void {
  swrMutate(() => true, undefined, { revalidate: false });
}

// ============================================================================
// Re-export useful types
// ============================================================================

export type { EndpointTypes, ApiError } from '@/lib/client';
export { Endpoints } from '@/lib/client';
