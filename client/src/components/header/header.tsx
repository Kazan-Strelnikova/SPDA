import { Typography } from '@mui/material';
import styles from './header.module.scss';
import { User } from '../user/user';

export const Header : React.FC = () => {
    return (
    <header className={styles.root}>
        <Typography variant='h3'className={styles.title}>
            Крутой тайтл
        </Typography>
        
        <User/>
    </header>
    );
}