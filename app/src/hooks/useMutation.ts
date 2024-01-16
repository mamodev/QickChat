import React, { useState } from "react";
import { ffetch } from "./fetch";

type MutationOptions = {
  onSuccess?: (data: unknown) => void;
  onError?: (error: Error) => void;
  onSettled?: (data: unknown, error: Error | null) => void;
};

type Method = "POST" | "PUT" | "PATCH" | "DELETE";

type MutationFetchState = "idle" | "loading" | "success" | "error";

type MutationState = {
  error: Error | null;
  status: MutationFetchState;
};

export function useMutation(method: Method, url: string, options?: MutationOptions) {
  const [state, setState] = useState<MutationState>({
    status: "idle",
    error: null,
  });

  const mutate = React.useCallback(
    async (data: unknown, params: Record<string, string> = {}) => {
      setState({
        status: "loading",
        error: null,
      });

      let urlWithParams = url;

      Object.keys(params).forEach((key) => {
        urlWithParams = urlWithParams.replace(`:${key}`, params[key]);
      });

      try {
        const json = await ffetch(urlWithParams, {
          method,
          headers: {
            "Content-Type": "application/json",
          },
          body: data ? JSON.stringify(data) : undefined,
        });

        options?.onSuccess?.(json);
        options?.onSettled?.(json, null);

        setState({
          status: "success",
          error: null,
        });

        return json;
      } catch (unknownError) {
        if (unknownError instanceof Error) {
          const err = unknownError;
          options?.onError?.(err);
          options?.onSettled?.(null, err);
          setState({
            status: "error",
            error: err,
          });

          throw err;
        }
      }
    },
    [method, url, options]
  );

  return {
    mutate,
    error: state.error,
    status: state.status,
    isLoading: state.status === "loading",
    isSuccess: state.status === "success",
    isError: state.status === "error",
    isIdle: state.status === "idle",
  };
}
