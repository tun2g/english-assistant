import { createMutation } from 'react-query-kit';
import { getQueryClient } from '../../lib/query-client';
import { QUERY_KEY } from '../../constants';
import { refreshUserToken } from '../../api/auth-api';
import type { RefreshTokenRequest, AuthResponse } from '../../types';

export const useRefreshTokenMutation = createMutation({
  mutationFn: (tokenData: RefreshTokenRequest): Promise<AuthResponse> => refreshUserToken(tokenData),
  onSuccess: (response: AuthResponse) => {
    const queryClient = getQueryClient();
    queryClient.setQueryData([QUERY_KEY.AUTH.PROFILE], response.user);
  },
});
