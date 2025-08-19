import { HomeHero } from './components/home-hero';
import { FeaturesSection } from './components/features-section';

export function HomeContent() {
  return (
    <div className="space-y-8">
      <HomeHero />
      <FeaturesSection />
    </div>
  );
}