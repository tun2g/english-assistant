import { Link } from 'react-router-dom';
import { Button } from '@english/ui';

export function HomeHero() {
  return (
    <section className="text-center space-y-4">
      <h1 className="text-4xl font-bold tracking-tight">
        Master English with Our Learning Platform
      </h1>
      <p className="text-xl text-muted-foreground max-w-2xl mx-auto">
        Improve your English skills with interactive lessons, vocabulary building, and personalized learning paths.
      </p>
      <div className="flex gap-4 justify-center">
        <Link to="/login">
          <Button size="lg">
            Get Started
          </Button>
        </Link>
        <Button variant="outline" size="lg">
          Learn More
        </Button>
      </div>
    </section>
  );
}