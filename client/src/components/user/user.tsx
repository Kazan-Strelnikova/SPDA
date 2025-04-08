import { useContext, useEffect } from 'react';
import { UserContext } from '../../contexts/UserContext';
import styles from './user.module.scss';
import { Avatar} from '@mui/material';

import * as React from 'react';
import { styled } from '@mui/material/styles';
import Tooltip, { TooltipProps, tooltipClasses } from '@mui/material/Tooltip';
import Typography from '@mui/material/Typography';
import { ButtonAKAM } from '../button/button';
import Cookies from 'js-cookie';
import { useNavigate } from 'react-router-dom';


const HtmlTooltip = styled(({ className, ...props }: TooltipProps) => (
  <Tooltip arrow placement='bottom-end' {...props} classes={{ popper: className }} />
))(() => ({
  [`& .${tooltipClasses.arrow}`]: {
    color: '#DDE2FF',
  },
  [`& .${tooltipClasses.tooltip}`]: {
    borderRadius: '16px',
    padding: '15px 20px',
    backgroundColor: '#DDE2FF',
    maxWidth: 320,
  },
}));


export const User : React.FC = () => {
  const userContext = useContext(UserContext);
  const navigate = useNavigate();
  if (!userContext) {
    throw new Error("Header must be used within a UserProvider");
  }
  const { user, setUser } = userContext;

  useEffect(() => setUser({id: '1', name: "Константин", lastName: 'Константинов', email: 'justcoolestgiraffe9@gmail.com'}), [setUser]);
  
  function handleLogOut(){
    setUser(null);
    Cookies.remove('token');
  }

  return (
    user ? 
      <HtmlTooltip title={
        <div className={styles.popper}>
          <Typography variant='body1' width='100%'>
            {user.name} {user.lastName}
          </Typography>
          <Typography variant='body2'>
            {user.email}
          </Typography>
          <ButtonAKAM filled onClick={handleLogOut}>
            Выйти
          </ButtonAKAM>
  
        </div>
      }>
        <Avatar className={styles.user}>{user.name.charAt(0).toUpperCase()}</Avatar>
      </HtmlTooltip>
      :  
      <div className={styles.buttons}>
        <ButtonAKAM outlined onClick={() => navigate('/login')}>
          Войти
        </ButtonAKAM>
        <ButtonAKAM filled onClick={() => navigate('/signup')}>
          Зарегестрироваться
        </ButtonAKAM>
      </div>
  );
}