import useSWR, { SWRConfiguration, mutate as swrMutate } from 'swr';
import useSWRMutation, { SWRMutationConfiguration } from 'swr/mutation';
import { api, EndpointTypes, ApiError, Endpoints } from '@/lib/client';

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
  /** リクエストを実行する関数 (型安全) */
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
 * @param endpoint - APIエンドポイント (Endpoints.XXX)
 * @param request - リクエストボディ (エンドポイントに応じて型が決定)
 * @param options - SWRオプション (任意)
 *
 * @example
 * ```typescript
 * // モンスター一覧を取得
 * const { data, isLoading } = useApi(Endpoints.GetMonsters, {});
 *
 * // 特定のモンスターを取得
 * const { data } = useApi(Endpoints.GetMonster, { id: 'monster-123' });
 *
 * // オプション付き
 * const { data } = useApi(Endpoints.GetMonsters, {}, {
 *   refreshInterval: 5000,
 *   revalidateOnFocus: true,
 * });
 *
 * // 条件付きfetch
 * const { data } = useApi(Endpoints.GetMonster, { id }, { enabled: !!id });
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
 * @param endpoint - APIエンドポイント (Endpoints.XXX)
 * @param options - SWR Mutationオプション (任意)
 *
 * @example
 * ```typescript
 * // 画像解析のmutation
 * const { trigger, isMutating } = useApiMutation(Endpoints.AnalyzeImage);
 *
 * // triggerはリクエストの型が決まっている
 * const result = await trigger({
 *   image_data: 'base64...',
 *   mime_type: 'image/jpeg',
 * });
 *
 * // オプション付き
 * const { trigger } = useApiMutation(Endpoints.GenerateImage, {
 *   onSuccess: (data) => console.log('Generated:', data),
 *   onError: (error) => console.error('Failed:', error),
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
    // triggerをTriggerFunction型にキャスト（SWRの複雑な型を簡略化）
    trigger: trigger as unknown as TriggerFunction<P>,
    data,
    error,
    isMutating,
    reset,
  };
}

// ============================================================================
// Pre-built Hooks (便利なショートカット)
// ============================================================================

/** モンスター一覧を取得 */
export function useMonsters(options?: UseApiOptions<'/monster/v1/GetMonsters'>) {
  return useApi(Endpoints.GetMonsters, {}, options);
}

/** 特定のモンスターを取得 */
export function useMonster(
  request: EndpointTypes['/monster/v1/GetMonster']['request'],
  options?: UseApiOptions<'/monster/v1/GetMonster'>
) {
  return useApi(Endpoints.GetMonster, request, options);
}

/** トラッシュ一覧を取得 */
export function useTrashs(options?: UseApiOptions<'/trash/v1/GetTrashs'>) {
  return useApi(Endpoints.GetTrashs, {}, options);
}

/** ヘルスチェック */
export function useHealthz(options?: UseApiOptions<'/healthz/v1/Healthz'>) {
  return useApi(Endpoints.Healthz, {}, options);
}

// ============================================================================
// Global Cache Functions
// ============================================================================

/**
 * 特定のAPIキャッシュを再検証
 * @example
 * ```typescript
 * await revalidateApi(Endpoints.GetMonsters, {});
 * ```
 */
export async function revalidateApi<P extends keyof EndpointTypes>(
  endpoint: P,
  request: EndpointTypes[P]['request']
): Promise<void> {
  const cacheKey = getCacheKey(endpoint, request);
  await swrMutate(cacheKey);
}

/** モンスター一覧キャッシュを再検証 */
export const revalidateMonsters = () => revalidateApi(Endpoints.GetMonsters, {});

/** 特定モンスターのキャッシュを再検証 */
export const revalidateMonster = (id: string) =>
  revalidateApi(Endpoints.GetMonster, { id });

/** トラッシュ一覧キャッシュを再検証 */
export const revalidateTrashs = () => revalidateApi(Endpoints.GetTrashs, {});

/** 全SWRキャッシュをクリア */
export function clearAllCache(): void {
  swrMutate(() => true, undefined, { revalidate: false });
}

// ============================================================================
// Re-exports
// ============================================================================

export type { EndpointTypes, ApiError } from '@/lib/client';
export { Endpoints } from '@/lib/client';
