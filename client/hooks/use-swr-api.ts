import useSWR, { SWRConfiguration, mutate as swrMutate } from 'swr';
import useSWRMutation, { SWRMutationConfiguration } from 'swr/mutation';
import { api, EndpointTypes, ApiError } from '@/lib/client';

// ============================================================================
// Types
// ============================================================================

/** useApi の戻り値の型 */
export type UseApiResult<P extends keyof EndpointTypes> = {
  data: EndpointTypes[P]['response'] | null;
  error: ApiError | null;
  isLoading: boolean;
  isValidating: boolean;
  mutate: () => Promise<EndpointTypes[P]['response'] | undefined>;
};

/** useApi のオプション型 */
export type UseApiOptions<P extends keyof EndpointTypes> = {
  headers?: Record<string, string>;
  enabled?: boolean;
} & Omit<SWRConfiguration<EndpointTypes[P]['response'], ApiError>, 'fetcher'>;

/** Trigger関数の型 */
export type TriggerFunction<P extends keyof EndpointTypes> = (
  request: EndpointTypes[P]['request'],
  options?: {
    throwOnError?: boolean;
    revalidate?: boolean;
    populateCache?: boolean;
    optimisticData?: EndpointTypes[P]['response'];
    rollbackOnError?: boolean;
  }
) => Promise<EndpointTypes[P]['response'] | undefined>;

/** useApiMutation の戻り値の型 */
export type UseApiMutationResult<P extends keyof EndpointTypes> = {
  trigger: TriggerFunction<P>;
  data: EndpointTypes[P]['response'] | undefined;
  error: ApiError | undefined;
  isMutating: boolean;
  reset: () => void;
};

/** useApiMutation のオプション型 */
export type UseApiMutationOptions<P extends keyof EndpointTypes> = {
  headers?: Record<string, string>;
} & Omit<
  SWRMutationConfiguration<
    EndpointTypes[P]['response'],
    ApiError,
    string,
    EndpointTypes[P]['request']
  >,
  'fetcher'
>;

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
// useApi - データ取得用 Hook
// ============================================================================

/**
 * 型安全なデータ取得用SWR Hook
 *
 * @param endpoint - APIエンドポイントパス (例: "/monster/v1/GetMonsters")
 * @param request - リクエストボディ (エンドポイントに応じて型が決定)
 * @param options - SWRオプション (任意)
 *
 * @example
 * ```typescript
 * // モンスター一覧を取得
 * const { data, isLoading } = useApi("/monster/v1/GetMonsters", {});
 *
 * // 特定のモンスターを取得
 * const { data } = useApi("/monster/v1/GetMonster", { id: "monster-123" });
 *
 * // オプション付き
 * const { data } = useApi("/monster/v1/GetMonsters", {}, {
 *   refreshInterval: 5000,
 *   revalidateOnFocus: true,
 * });
 *
 * // 条件付きfetch
 * const { data } = useApi("/monster/v1/GetMonster", { id }, { enabled: !!id });
 * ```
 */
export function useApi<P extends keyof EndpointTypes>(
  endpoint: P,
  request: EndpointTypes[P]['request'],
  options?: UseApiOptions<P>
): UseApiResult<P> {
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
// useApiMutation - データ変更用 Hook
// ============================================================================

/**
 * 型安全なデータ変更用SWR Mutation Hook
 *
 * @param endpoint - APIエンドポイントパス (例: "/gemini/v1/AnalyzeImage")
 * @param options - SWR Mutationオプション (任意)
 *
 * @example
 * ```typescript
 * // 画像解析のmutation
 * const { trigger, isMutating } = useApiMutation("/gemini/v1/AnalyzeImage");
 *
 * // triggerはリクエストの型が決まっている
 * const result = await trigger({
 *   image_data: "base64...",
 *   mime_type: "image/jpeg",
 * });
 *
 * // オプション付き
 * const { trigger } = useApiMutation("/gemini/v1/GenerateImage", {
 *   onSuccess: (data) => console.log("Generated:", data),
 *   onError: (error) => console.error("Failed:", error),
 * });
 * ```
 */
export function useApiMutation<P extends keyof EndpointTypes>(
  endpoint: P,
  options?: UseApiMutationOptions<P>
): UseApiMutationResult<P> {
  const { headers, ...swrOptions } = options ?? {};

  const fetcher = async (
    _key: string,
    { arg }: { arg: EndpointTypes[P]['request'] }
  ): Promise<EndpointTypes[P]['response']> => {
    const response = await api(endpoint, arg, { headers });
    return response.data;
  };

  const { trigger, data, error, isMutating, reset } = useSWRMutation<
    EndpointTypes[P]['response'],
    ApiError,
    string,
    EndpointTypes[P]['request']
  >(endpoint, fetcher, swrOptions);

  return {
    trigger: trigger as unknown as TriggerFunction<P>,
    data,
    error,
    isMutating,
    reset,
  };
}

// ============================================================================
// Global Cache Functions
// ============================================================================

/**
 * 特定のAPIキャッシュを再検証
 * @example
 * ```typescript
 * await revalidateApi("/monster/v1/GetMonsters", {});
 * ```
 */
export async function revalidateApi<P extends keyof EndpointTypes>(
  endpoint: P,
  request: EndpointTypes[P]['request']
): Promise<void> {
  const cacheKey = getCacheKey(endpoint, request);
  await swrMutate(cacheKey);
}

/**
 * 全SWRキャッシュをクリア
 */
export function clearAllCache(): void {
  swrMutate(() => true, undefined, { revalidate: false });
}

// ============================================================================
// Re-exports
// ============================================================================

export type { EndpointTypes, ApiError } from '@/lib/client';
