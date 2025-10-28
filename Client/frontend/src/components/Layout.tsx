import { Navigate, Outlet } from "react-router-dom";
import { useAuth } from "../contexts/AuthContext";

export const PublicLayout = () => {
  return (
    <div>
      <Outlet />
    </div>
  );
};

export const ProtectedLayout = () => {
  const { user } = useAuth();

  if (!user) {
    return <Navigate to="/login" />;
  }

  return (
    <div>
      <Outlet />
    </div>
  );
};
