import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@english/ui';
import { LoginForm } from './components/login-form';
import type { LoginRequest } from '@english/shared';

interface LoginContentProps {
  onSubmit: (data: LoginRequest) => void;
  isLoading: boolean;
  error: string | null;
}

export function LoginContent({ onSubmit, isLoading, error }: LoginContentProps) {
  return (
    <div className="flex justify-center items-center min-h-[60vh]">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle>Login</CardTitle>
          <CardDescription>
            Enter your credentials to access your account
          </CardDescription>
        </CardHeader>
        <CardContent>
          <LoginForm 
            onSubmit={onSubmit}
            isLoading={isLoading}
            error={error}
          />
        </CardContent>
      </Card>
    </div>
  );
}