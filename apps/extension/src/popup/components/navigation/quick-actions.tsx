import React from 'react';
import { Block, List, ListItem } from 'framework7-react';

export function QuickActions() {
  return (
    <>
      <Block strong>
        <p>Navigate to a YouTube video to use dual-language transcripts.</p>
      </Block>
      
      <List>
        <ListItem 
          title="Quick Translate"
          subtitle="Translate selected text"
          onClick={() => console.log('Quick translate clicked')}
        />
        <ListItem 
          title="Word Practice"
          subtitle="Practice vocabulary"
          onClick={() => console.log('Word practice clicked')}
        />
        <ListItem 
          title="Learning Progress"
          subtitle="View your progress"
          onClick={() => console.log('Progress clicked')}
        />
      </List>
    </>
  );
}