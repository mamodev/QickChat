import { useEffect, useState } from "react";

let userID = window.localStorage.getItem("userID");

type AuthListener = (id: string | null) => void;

const listeners: AuthListener[] = [];

function addAuthListener(listener: AuthListener) {
  listeners.push(listener);
  return () => {
    const index = listeners.indexOf(listener);
    listeners.splice(index, 1);
  };
}

function notifyListeners() {
  listeners.forEach((listener) => listener(userID));
}

export function login(id: string) {
  userID = id;
  window.localStorage.setItem("userID", id);
  notifyListeners();
}

export function logout() {
  userID = null;
  window.localStorage.removeItem("userID");
  notifyListeners();
}

export function useAuth() {
  const [authenticated, setAuthenticated] = useState<boolean>(!!userID);

  useEffect(() => {
    return addAuthListener((id) => setAuthenticated(!!id));
  }, []);

  return authenticated;
}

export function useUserID() {
  if (!userID) {
    throw new Error("No user ID available");
  }

  const [id, setID] = useState<string>(userID);

  useEffect(() => {
    return addAuthListener((id) => setID(id!));
  }, []);

  return id;
}
