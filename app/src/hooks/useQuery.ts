import React, { useState } from "react";
import { z } from "zod";
import { addQueryListener } from "../context/queryContext";
import { ApiError, safeFetch } from "./fetch";

type QueryState<T> = {
  data: T | null;
  error: ApiError | null;
  status: "idle" | "loading" | "error" | "success";
};

export function useQuery<T extends z.ZodTypeAny>(url: string, schema: T) {
  const [state, setState] = useState<QueryState<z.infer<T>>>({
    data: null,
    error: null,
    status: "idle",
  });

  React.useEffect(() => {
    let controller = new AbortController();

    const fetchData = async () => {
      try {
        setState((old) => ({ ...old, status: "loading" }));
        const data = await safeFetch(url, schema, {
          signal: controller.signal,
        });
        setState({ data, error: null, status: "success" });
      } catch (error) {
        if (error instanceof ApiError) {
          setState({ data: null, error, status: "error" });
        }
      }
    };

    const onChange = (updater?: (oldValue: unknown) => unknown) => {
      if (updater) {
        setState((old) => ({ ...old, data: updater(old.data) }));
      } else {
        controller.abort();
        controller = new AbortController();
        fetchData();
      }
    };

    const cleanUp = addQueryListener(url, onChange);
    fetchData();

    return () => {
      controller.abort();
      cleanUp();
    };
  }, [url]);

  return {
    data: state.data,
    error: state.error,
    status: state.status,
    isError: state.status === "error",
    isLoading: state.status === "loading",
    isSuccess: state.status === "success",
    isIdle: state.status === "idle",
  };
}
