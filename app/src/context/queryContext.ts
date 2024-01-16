type QueryListener = (updater?: (oldValue: unknown) => unknown) => void;

const queryListeners = new Map<string, Set<QueryListener>>();

export function addQueryListener(query: string, listener: QueryListener) {
  query = query.split("?")[0];

  if (!queryListeners.has(query)) {
    queryListeners.set(query, new Set());
  }

  queryListeners.get(query)?.add(listener);
  return () => removeQueryListener(query, listener);
}

export function removeQueryListener(query: string, listener: QueryListener) {
  queryListeners.get(query)?.delete(listener);
}

export function invalidateQuery(query: string) {
  queryListeners.get(query)?.forEach((listener) => listener());
}

export function updateQuery(query: string, updater: (oldValue: unknown) => unknown) {
  queryListeners.get(query)?.forEach((listener) => listener(updater));
}
