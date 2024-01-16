import { useEffect, useRef } from "react";

export function usePolling(func: () => Promise<void>, interval: number) {
  const timeoutRef = useRef<number | undefined>(undefined);

  useEffect(() => {
    const execFunc = async () => {
      try {
        await func();
      } catch (e) {
        console.error(e);
      }

      if (!timeoutRef.current) {
        return;
      }

      timeoutRef.current = setTimeout(execFunc, interval) as unknown as number;
    };

    timeoutRef.current = setTimeout(execFunc, interval) as unknown as number;

    return () => {
      const timeout = timeoutRef.current;
      timeoutRef.current = undefined;
      clearTimeout(timeout);
    };
  }, [interval, func]);
}
