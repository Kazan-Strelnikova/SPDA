import React, { createContext, useState, ReactNode, useEffect } from "react";
import { getUserByToken } from "../http/get-user-by-token";
import Cookies from 'js-cookie';

export interface User {
  name: string;
  lastName: string;
  email: string;
}

interface UserContextType {
  user: User | undefined;
  setUser: (user: User | undefined) => void;
}

export const UserContext = createContext<UserContextType | undefined>(undefined);

interface UserProviderProps {
  children: ReactNode;
}

export const UserProvider: React.FC<UserProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | undefined>();
  

  return (
    <UserContext.Provider value={{ user, setUser }}>
      {children}
    </UserContext.Provider>
  );
};
