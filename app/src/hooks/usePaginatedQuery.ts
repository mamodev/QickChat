import React, { useState } from "react";
import { z } from "zod";
import { addQueryListener } from "../context/queryContext";
import { ApiError, safeFetch } from "./fetch";

type QueryState<T> = {
  pages: T[];
  error: ApiError | null;
  status: "idle" | "loading" | "error" | "success";
  hasMore: boolean;
};

export function usePaginatedQuery<T extends z.ZodTypeAny>(_url: string, schema: T) {
  const [state, setState] = useState<QueryState<z.infer<T>[]>>({
    pages: [],
    hasMore: true,
    error: null,
    status: "idle",
  });

  const abortController = React.useRef<Record<number, AbortController>>({});

  const fetchData = React.useCallback(
    async (pageNumber: number) => {
      try {
        setState((old) => ({ ...old, status: "loading" }));

        const [baseUrl, query] = _url.split("?");
        const searchParams = new URLSearchParams(query);
        searchParams.set("page", pageNumber.toString());

        if (abortController.current[pageNumber]) abortController.current[pageNumber].abort();

        abortController.current[pageNumber] = new AbortController();

        const data = await safeFetch(`${baseUrl}?${searchParams.toString()}`, schema.array(), {
          signal: abortController.current[pageNumber].signal,
        });

        setState((old) => {
          if (pageNumber < old.pages.length) {
            return {
              pages: old.pages.map((page, index) => (index === pageNumber ? data : page)),
              hasMore: data.length > 0,
              error: null,
              status: "success",
            };
          }

          return {
            pages: [...old.pages, data],
            hasMore: data.length > 0,
            error: null,
            status: "success",
          };
        });
      } catch (error) {
        if (error instanceof ApiError) {
          setState((old) => ({ ...old, error: error as ApiError, status: "error" }));
        }
      }
    },
    [_url]
  );

  React.useEffect(() => {
    fetchData(0);
  }, [_url, fetchData]);

  React.useEffect(() => {
    const onInvalidate = () => {
      for (let i = 0; i < state.pages.length; i++) {
        fetchData(i);
      }
    };

    const cleanUp = addQueryListener(_url, onInvalidate);

    return () => {
      cleanUp();
    };
  }, [state.pages, _url, fetchData]);

  const fetchMore = () => {
    if (state.hasMore) {
      fetchData(state.pages.length);
    }
  };

  return {
    data: state.pages.reduce((acc, page) => [...acc, ...page], []),
    hasMore: state.hasMore,
    fetchMore,
    error: state.error,
    status: state.status,
    isError: state.status === "error",
    isLoading: state.status === "loading",
    isSuccess: state.status === "success",
    isIdle: state.status === "idle",
  };
}
