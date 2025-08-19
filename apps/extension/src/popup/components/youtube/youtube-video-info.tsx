import React from 'react';
import { Block } from 'framework7-react';
import type { PageInfo } from '../../../shared/types/extension-types';

interface YouTubeVideoInfoProps {
  pageInfo: PageInfo;
}

export function YouTubeVideoInfo({ pageInfo }: YouTubeVideoInfoProps) {
  return (
    <Block strong>
      <h4>🎥 YouTube Video Detected</h4>
      <p className="text-color-gray">{pageInfo.title}</p>
    </Block>
  );
}