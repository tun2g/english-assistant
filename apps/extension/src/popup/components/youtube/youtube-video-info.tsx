import React from 'react';
import type { PageInfo } from '../../../shared/types/extension-types';

interface YouTubeVideoInfoProps {
  pageInfo: PageInfo;
}

export function YouTubeVideoInfo({ pageInfo }: YouTubeVideoInfoProps) {
  return (
    <div className="border-b p-4">
      <h4 className="mb-2 font-medium">🎥 YouTube Video Detected</h4>
      <p className="text-muted-foreground text-sm">{pageInfo.title}</p>
    </div>
  );
}
