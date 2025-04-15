import React, { createContext, useState, ReactNode } from "react";

// Define User Type
export interface User {
  name: string;
  lastName: string;
  email: string;
}

// Define Context Type
interface UserContextType {
  user: User | null;
  setUser: (user: User | null) => void;
}

// Create Context with Default Values
export const UserContext = createContext<UserContextType | undefined>(undefined);

// Create Provider Component
interface UserProviderProps {
  children: ReactNode;
}

export const UserProvider: React.FC<UserProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);

  return (
    <UserContext.Provider value={{ user, setUser }}>
      {children}
    </UserContext.Provider>
  );
};
