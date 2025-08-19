// Storage interface for unified access across web and extension environments
export interface StorageAdapter {
  getItem(key: string): Promise<string | null>;
  setItem(key: string, value: string): Promise<void>;
  removeItem(key: string): Promise<void>;
  clear(): Promise<void>;
}

// Web storage adapter (localStorage/sessionStorage)
export class WebStorageAdapter implements StorageAdapter {
  constructor(private storage: Storage = localStorage) {}

  async getItem(key: string): Promise<string | null> {
    try {
      return this.storage.getItem(key);
    } catch (error) {
      console.error('WebStorage getItem error:', error);
      return null;
    }
  }

  async setItem(key: string, value: string): Promise<void> {
    try {
      this.storage.setItem(key, value);
    } catch (error) {
      console.error('WebStorage setItem error:', error);
      throw error;
    }
  }

  async removeItem(key: string): Promise<void> {
    try {
      this.storage.removeItem(key);
    } catch (error) {
      console.error('WebStorage removeItem error:', error);
      throw error;
    }
  }

  async clear(): Promise<void> {
    try {
      this.storage.clear();
    } catch (error) {
      console.error('WebStorage clear error:', error);
      throw error;
    }
  }
}

// Extension storage adapter (chrome.storage)
export class ExtensionStorageAdapter implements StorageAdapter {
  constructor(private area: 'local' | 'sync' = 'local') {}

  async getItem(key: string): Promise<string | null> {
    try {
      if (typeof chrome !== 'undefined' && chrome.storage) {
        const result = await chrome.storage[this.area].get(key);
        return result[key] || null;
      }
      // Fallback to localStorage if chrome.storage is not available
      return localStorage.getItem(key);
    } catch (error) {
      console.error('ExtensionStorage getItem error:', error);
      return null;
    }
  }

  async setItem(key: string, value: string): Promise<void> {
    try {
      if (typeof chrome !== 'undefined' && chrome.storage) {
        await chrome.storage[this.area].set({ [key]: value });
      } else {
        // Fallback to localStorage if chrome.storage is not available
        localStorage.setItem(key, value);
      }
    } catch (error) {
      console.error('ExtensionStorage setItem error:', error);
      throw error;
    }
  }

  async removeItem(key: string): Promise<void> {
    try {
      if (typeof chrome !== 'undefined' && chrome.storage) {
        await chrome.storage[this.area].remove(key);
      } else {
        // Fallback to localStorage if chrome.storage is not available
        localStorage.removeItem(key);
      }
    } catch (error) {
      console.error('ExtensionStorage removeItem error:', error);
      throw error;
    }
  }

  async clear(): Promise<void> {
    try {
      if (typeof chrome !== 'undefined' && chrome.storage) {
        await chrome.storage[this.area].clear();
      } else {
        // Fallback to localStorage if chrome.storage is not available
        localStorage.clear();
      }
    } catch (error) {
      console.error('ExtensionStorage clear error:', error);
      throw error;
    }
  }
}

// Memory storage adapter (fallback for environments without persistent storage)
export class MemoryStorageAdapter implements StorageAdapter {
  private storage = new Map<string, string>();

  async getItem(key: string): Promise<string | null> {
    return this.storage.get(key) || null;
  }

  async setItem(key: string, value: string): Promise<void> {
    this.storage.set(key, value);
  }

  async removeItem(key: string): Promise<void> {
    this.storage.delete(key);
  }

  async clear(): Promise<void> {
    this.storage.clear();
  }
}

// Environment detection and adapter creation
function createStorageAdapter(): StorageAdapter {
  // Extension environment detection
  if (typeof chrome !== 'undefined' && chrome.storage) {
    return new ExtensionStorageAdapter('local');
  }
  
  // Web environment detection
  if (typeof window !== 'undefined' && window.localStorage) {
    return new WebStorageAdapter(localStorage);
  }
  
  // Fallback to memory storage
  console.warn('No persistent storage available, using memory storage');
  return new MemoryStorageAdapter();
}

// Universal storage instance
export const universalStorage = createStorageAdapter();

// High-level storage utilities
export async function getStorageItem<T>(key: string, defaultValue?: T): Promise<T | null> {
  try {
    const value = await universalStorage.getItem(key);
    if (value === null) {
      return defaultValue || null;
    }
    return JSON.parse(value);
  } catch (error) {
    console.error('getStorageItem error:', error);
    return defaultValue || null;
  }
}

export async function setStorageItem<T>(key: string, value: T): Promise<void> {
  try {
    const stringValue = JSON.stringify(value);
    await universalStorage.setItem(key, stringValue);
  } catch (error) {
    console.error('setStorageItem error:', error);
    throw error;
  }
}

export async function removeStorageItem(key: string): Promise<void> {
  try {
    await universalStorage.removeItem(key);
  } catch (error) {
    console.error('removeStorageItem error:', error);
    throw error;
  }
}

export async function clearStorage(): Promise<void> {
  try {
    await universalStorage.clear();
  } catch (error) {
    console.error('clearStorage error:', error);
    throw error;
  }
}
