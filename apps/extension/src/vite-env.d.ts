/// <reference types="vite/client" />

// CSS module declarations for Vite
declare module '*.css' {
  const content: string;
  export default content;
}

declare module '*.scss' {
  const content: string;
  export default content;
}

declare module '*.sass' {
  const content: string;
  export default content;
}

// CSS module declarations with ?inline suffix for Shadow DOM injection
declare module '*.css?inline' {
  const content: string;
  export default content;
}

declare module '*.scss?inline' {
  const content: string;
  export default content;
}

declare module '*.sass?inline' {
  const content: string;
  export default content;
}

// CSS module declarations with ?url suffix
declare module '*.css?url' {
  const content: string;
  export default content;
}

declare module '*.scss?url' {
  const content: string;
  export default content;
}

// Other asset types for completeness
declare module '*.svg' {
  const content: string;
  export default content;
}

declare module '*.png' {
  const content: string;
  export default content;
}

declare module '*.jpg' {
  const content: string;
  export default content;
}

declare module '*.jpeg' {
  const content: string;
  export default content;
}

declare module '*.gif' {
  const content: string;
  export default content;
}

declare module '*.webp' {
  const content: string;
  export default content;
}

declare module '*.woff' {
  const content: string;
  export default content;
}

declare module '*.woff2' {
  const content: string;
  export default content;
}

declare module '*.ttf' {
  const content: string;
  export default content;
}

declare module '*.otf' {
  const content: string;
  export default content;
}

// Environment variables interface
interface ImportMetaEnv {
  readonly VITE_API_URL: string;
  readonly VITE_APP_NAME: string;
  readonly VITE_PORT: string;
  readonly BUILD_TARGET: string;
  readonly NODE_ENV: string;
  readonly INLINE_RUNTIME_CHUNK: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}

// Chrome extension specific global types
declare global {
  interface Window {
    __englishExtensionDebug?: {
      getComponentDebugInfo?: () => any[];
      reactRoots?: number;
    };
    __englishExtensionReactRoots?: number;
  }
}
