import { z } from "zod";
import { logout, useUserID } from "../context/authContext";
import { useQuery } from "../hooks/useQuery";
import { useState } from "react";
import { UserSearch } from "./UserSearch";
import { Request } from "./Requests";
import { Dialog } from "./ui/Dialog";
import { useMutation } from "../hooks/useMutation";

const userSchema = z.object({
  id: z.string(),
  username: z.string(),
  email: z.string(),
});

export function Profile() {
  const [sarchUser, setSearchUser] = useState(false);
  const [showRequests, setShowRequests] = useState(false);

  const id = useUserID();

  const { data: user } = useQuery(`/api/user/${id}`, userSchema);

  const { mutate: _logout } = useMutation("POST", "/api/logout");

  function handleLogout() {
    _logout(null).then(() => {
      logout();
    });
  }

  return (
    <>
      {sarchUser && (
        <Dialog title="Contacts" open={sarchUser} onClose={() => setSearchUser(false)}>
          <UserSearch />
        </Dialog>
      )}

      {showRequests && (
        <Dialog title="Requests" open={showRequests} onClose={() => setShowRequests(false)}>
          <Request />
        </Dialog>
      )}
      <div className="stack v-center f-width">
        <div className="f-width">
          <hr />
          <div className="row v-center spacing-1 px-sm">
            <img className="rounded-sm" src={`/api/user/${id}/avatar`} width={40} height={40} />
            <div className="stack f-width">
              <div className="row v-center" style={{ justifyContent: "space-between" }}>
                <p className="title">{user?.username}</p>{" "}
                <a className="text-small text-secondary" onClick={handleLogout}>
                  logout ✖
                </a>
              </div>
              <p className="text-small text-secondary">{user?.email}</p>
            </div>
          </div>
          <hr />
        </div>

        <div className="row spacing-1" style={{ paddingBottom: 8 }}>
          <button
            className="btn primary"
            onClick={() => {
              setSearchUser(true);
            }}
          >
            Contatti
          </button>

          <button
            className="btn secondary"
            onClick={() => {
              setShowRequests(true);
            }}
          >
            Richieste
          </button>
        </div>
      </div>
    </>
  );
}
