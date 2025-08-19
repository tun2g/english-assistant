import { FeatureCard } from './feature-card';

const FEATURES = [
  {
    title: 'Interactive Lessons',
    description: 'Engage with dynamic content that adapts to your learning style',
    content: 'Practice speaking, listening, reading, and writing with our comprehensive lesson library.',
  },
  {
    title: 'Vocabulary Builder',
    description: 'Expand your vocabulary with spaced repetition learning',
    content: 'Learn new words and phrases with our intelligent review system that adapts to your progress.',
  },
  {
    title: 'Progress Tracking',
    description: 'Monitor your improvement with detailed analytics',
    content: 'Track your learning journey with comprehensive statistics and achievement badges.',
  },
] as const;

export function FeaturesSection() {
  return (
    <section className="grid md:grid-cols-3 gap-6">
      {FEATURES.map((feature) => (
        <FeatureCard
          key={feature.title}
          title={feature.title}
          description={feature.description}
          content={feature.content}
        />
      ))}
    </section>
  );
}