import { useEffect, useState } from 'react';

export function useDebouncedSearch(value: string, delay = 300): string {
  const [settled, setSettled] = useState(value);
  useEffect(() => {
    const id = setTimeout(() => setSettled(value), delay);
    return () => clearTimeout(id);
  }, [value, delay]);
  return settled;
}
