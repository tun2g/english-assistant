import { usePageInfoQuery } from '@/hooks';
import { Button, Card, CardContent } from '@english/ui';

export function QuickActions() {
  const { pageInfo } = usePageInfoQuery();

  return (
    <div className="space-y-4">
      {!pageInfo.isYouTube && (
        <Card>
          <CardContent className="p-4">
            <p className="text-muted-foreground text-sm">
              Navigate to a YouTube video to use dual-language transcripts.
            </p>
          </CardContent>
        </Card>
      )}

      <div className="space-y-2">
        <Button
          variant="outline"
          className="h-auto w-full justify-start p-4"
          onClick={() => console.log('Quick translate clicked')}
        >
          <div className="text-left">
            <div className="font-medium">Quick Translate</div>
            <div className="text-muted-foreground text-sm">Translate selected text</div>
          </div>
        </Button>

        <Button
          variant="outline"
          className="h-auto w-full justify-start p-4"
          onClick={() => console.log('Word practice clicked')}
        >
          <div className="text-left">
            <div className="font-medium">Word Practice</div>
            <div className="text-muted-foreground text-sm">Practice vocabulary</div>
          </div>
        </Button>

        <Button
          variant="outline"
          className="h-auto w-full justify-start p-4"
          onClick={() => console.log('Progress clicked')}
        >
          <div className="text-left">
            <div className="font-medium">Learning Progress</div>
            <div className="text-muted-foreground text-sm">View your progress</div>
          </div>
        </Button>
      </div>
    </div>
  );
}
