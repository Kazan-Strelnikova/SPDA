import { useContext, useEffect } from 'react';
import { UserContext } from '../../contexts/UserContext';
import styles from './user.module.scss';
import { Avatar, Zoom } from '@mui/material';
import { AccountCircleOutlined} from '@mui/icons-material';

import * as React from 'react';
import { styled } from '@mui/material/styles';
import Tooltip, { TooltipProps, tooltipClasses } from '@mui/material/Tooltip';
import Typography from '@mui/material/Typography';
import { ButtonAKAM } from '../button/button';

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

  if (!userContext) {
    throw new Error("Header must be used within a UserProvider");
  }

  const { user, setUser } = userContext;
  useEffect(() => setUser({id: '1', name: "Константин Константинов", email: 'kostyan.kostyanuch@mailmymail.commomom'}), [setUser]);
  

  return (
    user && user?.name ? 
      <HtmlTooltip title={
        <div className={styles.popper}>
          <Typography variant='body1' width='100%'>
            {user?.name}
          </Typography>
          <Typography variant='body2'>
            {user?.email}
          </Typography>
          <ButtonAKAM filled>
            Выйти
          </ButtonAKAM>
  
        </div>
      }>
        <Avatar className={styles.user}>{user.name.charAt(0).toUpperCase()}</Avatar>
      </HtmlTooltip>
      :  
      <HtmlTooltip title={
        <div className={styles.popper}>
          <Typography variant='body1' width='100%'>
            Ещё не вошли в аккаунт?
          </Typography>
          <div className={styles.buttons}>
            <ButtonAKAM outlined>
              Войти
            </ButtonAKAM>
            <ButtonAKAM filled>
              Зарегестрироваться
            </ButtonAKAM>
          </div>
        </div>
      }>
        <Avatar className={styles.user}><AccountCircleOutlined className={styles.userIcon} fontSize='large'/></Avatar>
      </HtmlTooltip>  
  );
}