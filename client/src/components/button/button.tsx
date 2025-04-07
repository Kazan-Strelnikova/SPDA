import { Button } from "@mui/material";
import styles from './button.module.scss';
import { MouseEventHandler, ReactNode } from "react";


type ButtonVariant = 'filled' | 'outlined';

interface ButtonProps {
    filled?: boolean;
    outlined?: boolean;
    onClick?: MouseEventHandler<HTMLButtonElement>;
    children: ReactNode;
}
export const ButtonAKAM : React.FC<ButtonProps> = ({filled=false, outlined=false, onClick, children}) => {
    return (
        <>
        {filled &&
        <Button onClick={onClick} variant="contained" disableElevation className={styles.filledButton}>{children}</Button>}
        {
        outlined &&  <Button onClick={onClick} variant="outlined" disableElevation className={styles.outlinedButton}>{children}</Button>}
        </>
    );
}