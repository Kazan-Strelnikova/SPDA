import React, { createContext, useState, ReactNode } from "react";

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
