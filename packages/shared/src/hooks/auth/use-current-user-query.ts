import { createQuery } from 'react-query-kit';
import { QUERY_KEY } from '../../constants';
import { getUserProfile } from '../../api/auth-api';
import type { User } from '../../types';

export const useCurrentUserQuery = createQuery({
  queryKey: [QUERY_KEY.AUTH.PROFILE],
  fetcher: (): Promise<User> => getUserProfile(),
});
