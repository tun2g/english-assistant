import { createMutation } from 'react-query-kit';
import { getQueryClient } from '../../lib/query-client';
import { loginUser } from '../../api/auth-api';
import { QUERY_KEY } from '../../constants';
import type { LoginRequest, AuthResponse } from '../../types';

export const useLoginMutation = createMutation({
  mutationFn: (credentials: LoginRequest): Promise<AuthResponse> => loginUser(credentials),
  onSuccess: data => {
    const queryClient = getQueryClient();

    queryClient.setQueryData([QUERY_KEY.AUTH.PROFILE], data.user);
    // Invalidate other queries that depend on auth
    queryClient.invalidateQueries({ queryKey: [QUERY_KEY.USER.PROFILE] });
  },
});
