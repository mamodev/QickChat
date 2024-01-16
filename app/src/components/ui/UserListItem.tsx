type Props = {
  profile_picture: string;
  username: string;
  email: string;

  action: React.ReactNode;
};

export function UserListItem(props: Props) {
  const { profile_picture, username, email, action } = props;

  return (
    <div
      className="row f-width v-center spacing-1 p-sm"
      style={{
        justifyContent: "space-between",
      }}
    >
      <div className="row v-center spacing-1">
        <img src={profile_picture} alt={username} width="40" height="40" style={{ borderRadius: "50%" }} />
        <div>
          <div className="text">{username}</div>
          <div className="text-small text-secondary">{email}</div>
        </div>
      </div>

      {action}
    </div>
  );
}
