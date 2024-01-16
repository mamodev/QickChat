// fetch middleware to handle authentication state

import { logout } from "../context/authContext";
import { z } from "zod";

const baseUrl = "";

export class ApiError extends Error {
  constructor(message: string, public code: number) {
    super(message);
  }
}

export async function ffetch(url: string, options: RequestInit = {}) {
  try {
    const response = await fetch(baseUrl + url, options);

    if (response.status === 401) {
      logout();
    }

    const txt = await response.text();
    let json;

    try {
      json = JSON.parse(txt);
    } catch (error) {
      throw new ApiError(txt, response.status);
    }

    if (!response.ok) {
      throw new ApiError(json.message, response.status);
    }

    return json;
  } catch (error) {
    if (error instanceof ApiError) {
      throw error;
    }

    if (error instanceof Error) {
      if (error.name === "AbortError") {
        throw error;
      }

      throw new ApiError(error.message, 500);
    }

    throw new ApiError("Unknown error", 500);
  }
}

export async function safeFetch<T extends z.ZodTypeAny>(
  url: string,
  schema: T,
  options: RequestInit = {}
): Promise<z.infer<T>> {
  const response = await ffetch(url, options);

  try {
    const json = schema.parse(response);
    return json;
  } catch (error) {
    throw new ApiError("Invalid response from server", 500);
  }
}
