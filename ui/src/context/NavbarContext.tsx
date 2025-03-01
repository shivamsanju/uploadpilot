import { createContext, ReactNode, useContext } from 'react';
import { useLocalStorageState } from '../hooks/localstorage';

// Define the context type
interface NavbarContextType {
  opened: boolean;
  toggle: () => void;
}

// Create the context with a default value
const NavbarContext = createContext<NavbarContextType | undefined>(undefined);

// Create a provider component
export const NavbarProvider = ({ children }: { children: ReactNode }) => {
  const [opened, setOpened] = useLocalStorageState('navbar-toggle', false);

  const toggle = () => setOpened(prev => !prev);

  return (
    <NavbarContext.Provider value={{ opened, toggle }}>
      {children}
    </NavbarContext.Provider>
  );
};

// Custom hook for consuming the context
export const useNavbar = () => {
  const context = useContext(NavbarContext);
  if (!context) {
    throw new Error('useNavbar must be used within a NavbarProvider');
  }
  return context;
};
