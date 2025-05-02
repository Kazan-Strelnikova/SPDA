import { Typography } from '@mui/material';
import styles from './header.module.scss';
import { User } from '../user/user';
import { Link } from 'react-router-dom';

export const Header : React.FC = () => {
    return (
    <header className={styles.root}>
        <Typography variant='h4'>
            <Link to='/' className={styles.title}>
            Eventify
            </Link>
        </Typography>
        
        <User/>
    </header>
    );
}