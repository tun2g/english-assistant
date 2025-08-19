import { useQuery } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import { useEffect } from 'react';
import { DashboardContent } from '../../components/dashboard/dashboard-content';
import { getUserProfile, hasValidToken, QUERY_KEY } from '@english/shared';

export function DashboardContainer() {
  const navigate = useNavigate();

  // Redirect if not authenticated
  useEffect(() => {
    if (!hasValidToken()) {
      navigate('/login');
    }
  }, [navigate]);

  const { data: userResponse, isLoading } = useQuery({
    queryKey: [QUERY_KEY.AUTH.PROFILE],
    queryFn: getUserProfile,
    enabled: hasValidToken(),
  });

  return (
    <DashboardContent
      user={userResponse?.data || null}
      isLoading={isLoading}
    />
  );
}