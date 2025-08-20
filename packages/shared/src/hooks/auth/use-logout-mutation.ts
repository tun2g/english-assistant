import { createMutation } from 'react-query-kit';
import { getQueryClient } from '../../lib/query-client';
import { logoutUser } from '../../api/auth-api';

export const useLogoutMutation = createMutation({
  mutationFn: (): Promise<null> => logoutUser(),
  onSuccess: () => {
    const queryClient = getQueryClient();
    queryClient.clear();
  },
});
