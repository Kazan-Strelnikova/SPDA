import { Button } from "@mui/material";
import styles from './button.module.scss';
import { ReactNode } from "react";


type ButtonVariant = 'filled' | 'outlined';

interface ButtonProps {
    filled?: boolean;
    outlined?: boolean;
    children: ReactNode;
}
export const ButtonAKAM : React.FC<ButtonProps> = ({filled=false, outlined=false, children}) => {
    return (
        <>
        {filled &&
        <Button variant="contained" disableElevation className={styles.filledButton}>{children}</Button>}
        {
        outlined &&  <Button variant="outlined" disableElevation className={styles.outlinedButton}>{children}</Button>}
        </>
    );
}