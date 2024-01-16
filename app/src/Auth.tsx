import { useState } from "react";
import { useMutation } from "./hooks/useMutation";
import { login } from "./context/authContext";
import { z } from "zod";

type Tabs = "login" | "signup";

export default function Auth() {
  const [tab, setTab] = useState<Tabs>("login");

  return (
    <div className="app-container v-center h-center">
      <div className="shadow-300 rounded-sm p-lg" style={{ minWidth: 300 }}>
        {tab === "login" && <Login onChange={setTab} />}
        {tab === "signup" && <Signup onChange={setTab} />}
      </div>
    </div>
  );
}

type TabProps = {
  onChange: (tab: Tabs) => void;
};

const authResponse = z.object({
  user_id: z.string(),
});

const onSuccess = (data: unknown) => {
  try {
    login(authResponse.parse(data).user_id);
  } catch (e) {
    console.log(e);
  }
};

function Login({ onChange }: TabProps) {
  const { mutate, error } = useMutation("POST", "api/login", {
    onSuccess,
  });

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const username = e.currentTarget.username.value;
    const password = e.currentTarget.password.value;
    mutate({ username, password });
  };

  return (
    <form onSubmit={handleSubmit} className="stack spacing-1">
      <p className="title py-sm">Login</p>

      <div>
        <div className="stack spacing-0">
          <label htmlFor="username">Username</label>
          <input className="input rounded-sm" type="text" id="username" placeholder="username" />
        </div>
      </div>

      <div>
        <div className="stack spacing-0">
          <label htmlFor="password">Password</label>
          <input className="input rounded-sm " type="password" id="password" placeholder="*******" />
        </div>
      </div>

      <p className="text-small text-secondary">
        Don't have an account?{" "}
        <a className="link bold" onClick={() => onChange("signup")}>
          Sign up
        </a>
      </p>

      {error && <p className="text c-error">Invalid credentials, plese retry...</p>}
      <button className="btn primary" type="submit">
        Login
      </button>
    </form>
  );
}

function Signup({ onChange }: TabProps) {
  const { mutate, error } = useMutation("POST", "api/register", {
    onSuccess,
  });

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const email = e.currentTarget.email.value;
    const username = e.currentTarget.username.value;
    const password = e.currentTarget.password.value;
    mutate({ username, password, email });
  };

  return (
    <form onSubmit={handleSubmit} className="stack spacing-1">
      <p className="title py-sm">Sign up</p>

      <div>
        <div className="stack spacing-0">
          <label htmlFor="email">Email</label>
          <input className="input rounded-sm" type="email" id="email" placeholder="email" />
        </div>
      </div>
      <div>
        <div className="stack spacing-0">
          <label htmlFor="username">Username</label>
          <input className="input rounded-sm" type="text" id="username" placeholder="username" />
        </div>
      </div>
      <div>
        <div className="stack spacing-0">
          <label htmlFor="password">Password</label>
          <input className="input rounded-sm" type="password" id="password" placeholder="*******" />
        </div>
      </div>

      <p className="text-small text-secondary">
        Already have an account?{" "}
        <a className="link bold" onClick={() => onChange("login")}>
          Login
        </a>
      </p>

      {error && <p className="text c-error">Something went wrong..</p>}
      <button className="btn primary" type="submit">
        Sign up
      </button>
    </form>
  );
}
