import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useMutation } from '@tanstack/react-query';
import { LoginContent } from '../../components/login/login-content';
import { loginUser, setAuthTokens, type LoginRequest } from '@english/shared';

export function LoginContainer() {
  const navigate = useNavigate();
  const [error, setError] = useState<string | null>(null);

  const loginMutation = useMutation({
    mutationFn: async (data: LoginRequest) => {
      const response = await loginUser(data);
      return response.data;
    },
    onSuccess: (authResponse) => {
      setAuthTokens(authResponse);
      navigate('/dashboard');
    },
    onError: (error) => {
      console.error('Login failed:', error);
      setError('Login failed. Please check your credentials.');
    },
  });

  function handleSubmit(data: LoginRequest) {
    setError(null);
    loginMutation.mutate(data);
  }

  return (
    <LoginContent
      onSubmit={handleSubmit}
      isLoading={loginMutation.isPending}
      error={error}
    />
  );
}